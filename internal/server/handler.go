package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GaIsBAX/Webhix/internal/domain"
)

type HookService interface {
	CreateHook(ctx context.Context, token string) (domain.Hook, error)
}

type HookHandler struct {
	mux     *http.ServeMux
	service HookService
}

func NewHookHandler(mux *http.ServeMux, srv HookService) *HookHandler {
	return &HookHandler{
		mux:     mux,
		service: srv,
	}
}

func (h *HookHandler) RegisterRoutes() {
	h.mux.HandleFunc("POST /webhook/r/{name}", h.WebHookHandler)
}

func (h *HookHandler) WebHookHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		SendError(w, http.StatusBadRequest, WithDetails(ErrBadRequest, ErrorDetailContract{
			Field:   "id",
			Message: "id is missing",
		}))
		return
	}

	hook, err := h.service.CreateHook(r.Context(), name)
	if err != nil {
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	data, err := json.Marshal(hook)
	if err != nil {
		SendError(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	SendSuccess(w, http.StatusCreated, data)
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
