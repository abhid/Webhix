package server

import "net/http"

// HealthHandler responds with 200 OK for liveness/readiness probes.
// It performs no authentication or backend checks so platform health
// checks (Coolify, Docker, Kubernetes) can reach it without credentials.
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
