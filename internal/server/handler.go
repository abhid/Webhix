package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/GaIsBAX/Webhix/internal/domain"
)

const maxBodySize = 5 << 20 // 5MB

type HookService interface {
	CreateHook(ctx context.Context, token string) (domain.Hook, error)
	ReceiveWebhook(ctx context.Context, token string, params domain.CreateWebhookRequestParams) (domain.WebhookRequest, error)
	ListWebhookRequests(ctx context.Context, token string) ([]domain.WebhookRequest, error)
}

type HookHandler struct {
	mux     *http.ServeMux
	service HookService
	baseURL string
}

func NewHookHandler(mux *http.ServeMux, srv HookService, baseURL string) *HookHandler {
	return &HookHandler{
		mux:     mux,
		service: srv,
		baseURL: baseURL,
	}
}

func (h *HookHandler) RegisterRoutes() {
	h.mux.HandleFunc("POST /api/endpoints", h.CreateEndpoint)
	h.mux.HandleFunc("GET /api/endpoints/{token}/requests", h.ListRequests)
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

	headersJSON, _ := json.Marshal(r.Header)

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
		slog.Error("read webhook body", "token", token, "err", err)
		SendError(w, http.StatusBadRequest, WithDetails(ErrBadRequest, ErrorDetailContract{
			Field:   "body",
			Message: "failed to read body",
		}))
		return
	}

	req, err := h.service.ReceiveWebhook(r.Context(), token, domain.CreateWebhookRequestParams{
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
