package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/GaIsBAX/Webhix/internal/domain"
)

type HookService interface {
	CreateHook(ctx context.Context, token string) (domain.Hook, error)
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
	h.mux.HandleFunc("POST /api/endpoints/{name}", h.CreateEndpoint)
	h.mux.HandleFunc("/r/{token}", h.ReceiveWebhook)
}

func (h *HookHandler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		SendError(w, http.StatusBadRequest, WithDetails(ErrBadRequest, ErrorDetailContract{
			Field:   "name",
			Message: "name is missing",
		}))
		return
	}

	hook, err := h.service.CreateHook(r.Context(), name)
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
	// TODO: захват входящего вебхука (headers, body, query, metadata)
	w.WriteHeader(http.StatusOK)
}

func Send(w http.ResponseWriter, status int, msg *ResponseContract) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func SendSuccess(w http.ResponseWriter, status int, data []byte) {
	Send(w, status, NewSuccessResponseContract(data))
}

func SendError(w http.ResponseWriter, status int, msg ErrorContract) {
	Send(w, status, NewErrorResponseContract(msg))
}
