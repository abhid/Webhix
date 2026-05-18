package server

import (
	"encoding/json"
	"net/http"
)

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

func DecodeRequest[T any](req *http.Request) (*T, error) {
	var body T

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
