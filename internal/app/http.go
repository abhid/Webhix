package app

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/hub"
	"github.com/GaIsBAX/Webhix/internal/server"
	"github.com/GaIsBAX/Webhix/internal/server/middleware"
	"github.com/GaIsBAX/Webhix/internal/web"
)

const readHeaderTimeout = 5 * time.Second

func newMux(cfg *config.Config, services *services, events *hub.Hub) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	registerWebhookRoutes(mux, cfg, services.hook, events)

	if err := registerStaticRoutes(mux); err != nil {
		return nil, err
	}

	return mux, nil
}

func registerWebhookRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
	hookService server.HookService,
	events *hub.Hub,
) {
	handler := server.NewHookHandler(
		mux,
		hookService,
		events,
		server.HookHandlerOptions{
			BaseURL:     cfg.BaseURL,
			MaxBodySize: cfg.MaxBodySize,
			ReadOnly:    cfg.ReadOnly,
		},
	)

	handler.RegisterRoutes()
}

func registerStaticRoutes(mux *http.ServeMux) error {
	staticSub, err := fs.Sub(web.Static, "static")
	if err != nil {
		return err
	}

	staticHandler := http.FileServer(http.FS(staticSub))

	mux.Handle("/ui/", http.StripPrefix("/ui/", staticHandler))
	mux.Handle("/", staticHandler)

	return nil
}

func newHTTPHandler(cfg *config.Config, mux *http.ServeMux) (http.Handler, error) {
	handler := http.Handler(mux)

	auth, err := newAuthMiddleware(cfg)
	if err != nil {
		return nil, err
	}

	handler = auth.Protect(handler)

	if len(cfg.TrustedProxies) > 0 {
		trustedProxies := middleware.NewTrustedProxies(cfg.TrustedProxies)
		if trustedProxies == nil {
			return nil, fmt.Errorf("invalid trusted proxies")
		}

		handler = trustedProxies.BehindProxy(handler)
	}

	return handler, nil
}

func newAuthMiddleware(cfg *config.Config) (*middleware.Auth, error) {
	password, secretKey, err := authCredentials(cfg)
	if err != nil {
		return nil, fmt.Errorf("auth setup: %w", err)
	}

	return middleware.NewAuth(password, secretKey), nil
}

func newHTTPServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

func authCredentials(cfg *config.Config) (password, secretKey string, err error) {
	if cfg.Password == "" && cfg.SecretKey == "" {
		return "", "", fmt.Errorf("auth is required: set WEBHIX_PASSWORD or WEBHIX_SECRET_KEY")
	}

	return cfg.Password, cfg.SecretKey, nil
}
