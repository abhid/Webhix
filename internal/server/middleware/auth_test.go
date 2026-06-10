package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newOKHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestProtectAllowsHealthzWithoutCredentials(t *testing.T) {
	auth := NewAuth("secret-password", "")
	handler := auth.Protect(newOKHandler())

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestProtectRejectsProtectedPathWithoutCredentials(t *testing.T) {
	auth := NewAuth("secret-password", "")
	handler := auth.Protect(newOKHandler())

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/endpoints", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
