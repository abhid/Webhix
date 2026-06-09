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
)

const DefaultMaxBodySize int64 = 5 << 20 // 5MB

type HookService interface {
	ListHooks(ctx context.Context) ([]domain.Hook, error)
	CreateHook(ctx context.Context, token string) (domain.Hook, error)
	ReceiveWebhook(ctx context.Context, token string, params domain.CreateWebhookRequestParams) (domain.WebhookRequest, domain.HookResponse, error)
	ListWebhookRequests(ctx context.Context, token string) ([]domain.WebhookRequest, error)
	GetHookResponse(ctx context.Context, token string) (domain.HookResponse, error)
	SetHookResponse(ctx context.Context, token string, params domain.UpsertHookResponseParams) (domain.HookResponse, error)
}

type EventBroker interface {
	Done() <-chan struct{}
	Subscribe(token string) (<-chan []byte, func())
	Publish(token string, data []byte)
}

type HookOptions struct {
	BaseURL     string
	MaxBodySize int64
	ReadOnly    bool
}

type HookDeps struct {
	Mux     *http.ServeMux
	Service HookService
	Hub     EventBroker
	Opts    HookOptions
}

type Hook struct {
	deps *HookDeps
}

func NewHook(deps *HookDeps) *Hook {
	if deps.Opts.MaxBodySize <= 0 {
		deps.Opts.MaxBodySize = DefaultMaxBodySize
	}

	return &Hook{deps: deps}
}

func (h *Hook) RegisterRoutes() {
	h.deps.Mux.HandleFunc("GET /api/endpoints", h.ListEndpoints)
	h.deps.Mux.HandleFunc("POST /api/endpoints", h.CreateEndpoint)
	h.deps.Mux.HandleFunc("GET /api/endpoints/{token}/requests", h.ListRequests)
	h.deps.Mux.HandleFunc("GET /api/endpoints/{token}/events", h.StreamEvents)
	h.deps.Mux.HandleFunc("GET /api/endpoints/{token}/response", h.GetResponse)
	h.deps.Mux.HandleFunc("PUT /api/endpoints/{token}/response", h.SetResponse)
	h.deps.Mux.HandleFunc("/r/{token}", h.ReceiveWebhook)
}

func (h *Hook) ListEndpoints(w http.ResponseWriter, r *http.Request) {
	hooks, err := h.deps.Service.ListHooks(r.Context())
	if err != nil {
		slog.Error("list endpoints", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	contracts := make([]EndpointListItemContract, len(hooks))
	for i, hook := range hooks {
		contracts[i] = EndpointListItemContract{
			ID:           hook.ID,
			Token:        hook.Token,
			Name:         hook.Name,
			URL:          h.deps.Opts.BaseURL + "/r/" + hook.Token,
			CreatedAt:    hook.CreatedAt,
			RequestCount: hook.RequestCount,
		}
	}

	data, err := json.Marshal(contracts)
	if err != nil {
		slog.Error("marshal endpoints", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	SendSuccess(w, http.StatusOK, data)
}

func (h *Hook) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	if h.readOnly(w) {
		return
	}

	contract, err := DecodeRequest[CreateEndpointRequestContract](r)
	if err != nil {
		slog.Error("create endpoint", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	hook, err := h.deps.Service.CreateHook(r.Context(), contract.Name)
	if err != nil {
		slog.Error("create endpoint", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	data, err := json.Marshal(CreateEndpointResponseContract{
		ID:    hook.ID,
		Token: hook.Token,
		Name:  hook.Name,
		URL:   h.deps.Opts.BaseURL + "/r/" + hook.Token,
	})
	if err != nil {
		slog.Error("marshal endpoint", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	SendSuccess(w, http.StatusCreated, data)
}

func (h *Hook) ReceiveWebhook(w http.ResponseWriter, r *http.Request) {
	if h.readOnly(w) {
		return
	}

	token := r.PathValue("token")

	headersJSON, err := json.Marshal(r.Header)
	if err != nil {
		slog.Error("marshal request headers", "err", err)
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.deps.Opts.MaxBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			SendError(w, http.StatusRequestEntityTooLarge, WithDetails(ErrPayloadTooLarge, ErrorDetailContract{
				Field:   "body",
				Message: fmt.Sprintf("body exceeds %d bytes limit", h.deps.Opts.MaxBodySize),
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

	req, customResp, err := h.deps.Service.ReceiveWebhook(r.Context(), token, domain.CreateWebhookRequestParams{
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

	h.deps.Hub.Publish(token, data)

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

func (h *Hook) ListRequests(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	reqs, err := h.deps.Service.ListWebhookRequests(r.Context(), token)
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

func (h *Hook) StreamEvents(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	ch, unsubscribe := h.deps.Hub.Subscribe(token)
	defer unsubscribe()

	flusher, canFlush := w.(http.Flusher)

	for {
		select {
		case <-r.Context().Done():
			return

		case <-h.deps.Hub.Done():
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

func (h *Hook) GetResponse(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	resp, err := h.deps.Service.GetHookResponse(r.Context(), token)
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

func (h *Hook) SetResponse(w http.ResponseWriter, r *http.Request) {
	if h.readOnly(w) {
		return
	}

	token := r.PathValue("token")

	contract, err := DecodeRequest[SetHookResponseRequestContract](r)
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

	resp, err := h.deps.Service.SetHookResponse(r.Context(), token, params)
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

func (h *Hook) readOnly(w http.ResponseWriter) bool {
	if !h.deps.Opts.ReadOnly {
		return false
	}

	SendError(w, http.StatusForbidden, ErrReadOnly)
	return true
}
