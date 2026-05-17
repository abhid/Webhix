package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBehindProxyUsesNearestUntrustedForwardedForAddress(t *testing.T) {
	tp := NewTrustedProxies([]string{"10.0.0.0/8"})
	require.NotNil(t, tp)

	handler := tp.BehindProxy(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "1.2.3.4:0", r.RemoteAddr)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://internal.local/hooks", nil)
	req.RemoteAddr = "10.1.2.3:4567"
	req.Header.Set("X-Forwarded-For", "8.8.8.8, 1.2.3.4")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}
