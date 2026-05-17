package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/GaIsBAX/Webhix/internal/hub"
)

const maxBodySize = 5 << 20 // 5MB

type HookService interface {
	CreateHook(ctx context.Context, token string) (domain.Hook, error)
	ReceiveWebhook(ctx context.Context, token string, params domain.CreateWebhookRequestParams) (domain.WebhookRequest, domain.HookResponse, error)
	ListWebhookRequests(ctx context.Context, token string) ([]domain.WebhookRequest, error)
	GetHookResponse(ctx context.Context, token string) (domain.HookResponse, error)
	SetHookResponse(ctx context.Context, token string, params domain.UpsertHookResponseParams) (domain.HookResponse, error)
}

type HookHandler struct {
	mux     *http.ServeMux
	service HookService
	baseURL string
	hub     *hub.Hub
}

func NewHookHandler(mux *http.ServeMux, srv HookService, baseURL string, hub *hub.Hub) *HookHandler {
	return &HookHandler{
		mux:     mux,
		service: srv,
		baseURL: baseURL,
		hub:     hub,
	}
}

func (h *HookHandler) RegisterRoutes() {
	h.mux.HandleFunc("POST /api/endpoints", h.CreateEndpoint)
	h.mux.HandleFunc("GET /api/endpoints/{token}/requests", h.ListRequests)
	h.mux.HandleFunc("GET /api/endpoints/{token}/events", h.StreamEvents)
	h.mux.HandleFunc("GET /api/endpoints/{token}/response", h.GetResponse)
	h.mux.HandleFunc("PUT /api/endpoints/{token}/response", h.SetResponse)
	h.mux.HandleFunc("/r/{token}", h.ReceiveWebhook)
}

func (h *HookHandler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	contract, err := DecodeContract[CreateEndpointRequestContract](r)
	if err != nil {
		slog.Error("create endpoint", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	hook, err := h.service.CreateHook(r.Context(), contract.Name)
	if err != nil {
		slog.Error("create endpoint", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	data, err := json.Marshal(CreateEndpointResponseContract{
		ID:    hook.ID,
		Token: hook.Token,
		Name:  hook.Name,
		URL:   h.baseURL + "/r/" + hook.Token,
	})
	if err != nil {
		slog.Error("marshal endpoint", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	SendSuccess(w, http.StatusCreated, data)
}

func (h *HookHandler) ReceiveWebhook(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	headersJSON, err := json.Marshal(r.Header)
	if err != nil {
		slog.Error("marshal request headers", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			SendError(w, http.StatusRequestEntityTooLarge, WithDetails(ErrPayloadTooLarge, ErrorDetailContract{
				Field:   "body",
				Message: "body exceeds 5MB limit",
			}))
			return
		}
		slog.Error("read webhook body", "err", err)
		SendError(w, http.StatusBadRequest, WithDetails(ErrBadRequest, ErrorDetailContract{
			Field:   "body",
			Message: "failed to read body",
		}))
		return
	}

	req, customResp, err := h.service.ReceiveWebhook(r.Context(), token, domain.CreateWebhookRequestParams{
		Method:      r.Method,
		Path:        r.URL.Path,
		Query:       r.URL.RawQuery,
		Headers:     string(headersJSON),
		Body:        body,
		RemoteAddr:  r.RemoteAddr,
		ContentType: r.Header.Get("Content-Type"),
		BodySize:    int64(len(body)),
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			SendError(w, http.StatusNotFound, ErrNotFound)
			return
		}
		slog.Error("receive webhook", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	data, err := json.Marshal(toWebhookRequestContract(req))
	if err != nil {
		slog.Error("marshal webhook request", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	h.hub.Publish(token, data)

	if customResp.StatusCode > 0 {
		for k, v := range customResp.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(int(customResp.StatusCode))
		if len(customResp.Body) > 0 {
			if _, err := w.Write(customResp.Body); err != nil {
				slog.Error("write custom response body", "err", err)
			}
		}
		return
	}

	SendSuccess(w, http.StatusOK, data)
}

func (h *HookHandler) ListRequests(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	reqs, err := h.service.ListWebhookRequests(r.Context(), token)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			SendError(w, http.StatusNotFound, ErrNotFound)
			return
		}
		slog.Error("list requests", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	contracts := make([]WebhookRequestContract, len(reqs))
	for i, req := range reqs {
		contracts[i] = toWebhookRequestContract(req)
	}

	data, err := json.Marshal(contracts)
	if err != nil {
		slog.Error("marshal requests", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	SendSuccess(w, http.StatusOK, data)
}

func (h *HookHandler) StreamEvents(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	ch, unsubscribe := h.hub.Subscribe(token)
	defer unsubscribe()

	flusher, canFlush := w.(http.Flusher)

	for {
		select {
		case <-r.Context().Done():
			return
		case <-h.hub.Done():
			return
		case data := <-ch:
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return
			}
			if canFlush {
				flusher.Flush()
			}
		}
	}
}

func (h *HookHandler) GetResponse(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	resp, err := h.service.GetHookResponse(r.Context(), token)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			SendError(w, http.StatusNotFound, ErrNotFound)
			return
		}
		slog.Error("get hook response", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	data, err := json.Marshal(toHookResponseContract(resp))
	if err != nil {
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	SendSuccess(w, http.StatusOK, data)
}

func (h *HookHandler) SetResponse(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	contract, err := DecodeContract[SetHookResponseRequestContract](r)
	if err != nil {
		SendError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	params := domain.UpsertHookResponseParams{
		StatusCode: int64(contract.StatusCode),
		Headers:    contract.Headers,
		Body:       []byte(contract.Body),
	}
	if err := params.Validate(); err != nil {
		SendError(w, http.StatusBadRequest, WithDetails(ErrBadRequest, ErrorDetailContract{
			Field:   "statusCode",
			Message: err.Error(),
		}))
		return
	}

	resp, err := h.service.SetHookResponse(r.Context(), token, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			SendError(w, http.StatusNotFound, ErrNotFound)
			return
		}
		slog.Error("set hook response", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	data, err := json.Marshal(toHookResponseContract(resp))
	if err != nil {
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	SendSuccess(w, http.StatusOK, data)
}

