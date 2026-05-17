package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/core"
	"github.com/GaIsBAX/Webhix/internal/hub"
	"github.com/GaIsBAX/Webhix/internal/server"
	"github.com/GaIsBAX/Webhix/internal/server/middleware"
	"github.com/GaIsBAX/Webhix/internal/store"
	"github.com/GaIsBAX/Webhix/internal/web"
)

type App struct {
	srv  *http.Server
	cfg  *config.Config
	deps *Deps
	hub  *hub.Hub
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	mux := http.NewServeMux()

	deps, err := NewDeps(ctx, cfg)
	if err != nil {
		return nil, err
	}

	eventHub := hub.New()

	hookRepository := store.NewHookRepository(deps.DB.DB)
	hookService := core.NewHookService(hookRepository)
	hookHandler := server.NewHookHandler(mux, hookService, cfg.BaseURL, eventHub)

	hookHandler.RegisterRoutes()

	staticSub, err := fs.Sub(web.Static, "static")
	if err != nil {
		return nil, err
	}
	staticFS := http.FileServer(http.FS(staticSub))
	mux.Handle("/ui/", http.StripPrefix("/ui/", staticFS))
	mux.Handle("/", staticFS)

	password, secretKey, err := resolveAuth(cfg)
	if err != nil {
		return nil, fmt.Errorf("auth setup: %w", err)
	}

	handler := http.Handler(mux)
	handler = middleware.NewAuth(password, secretKey).Protect(handler)

	if len(cfg.TrustedProxies) > 0 {
		trustedProxies := middleware.NewTrustedProxies(cfg.TrustedProxies)
		if trustedProxies == nil {
			return nil, fmt.Errorf("invalid trusted proxies")
		}

		handler = trustedProxies.BehindProxy(handler)
	}

	return &App{
		srv:  &http.Server{Addr: cfg.Addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second},
		cfg:  cfg,
		deps: deps,
		hub:  eventHub,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("webhix started", "addr", a.cfg.Addr, "base_url", a.cfg.BaseURL)
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return a.Shutdown()
	}
}

func resolveAuth(cfg *config.Config) (password, secretKey string, err error) {
	password = cfg.Password
	secretKey = cfg.SecretKey

	if password == "" && secretKey == "" {
		b := make([]byte, 16)
		if _, err = rand.Read(b); err != nil {
			return "", "", fmt.Errorf("generate password: %w", err)
		}
		password = base64.URLEncoding.EncodeToString(b)[:22]
		slog.Warn("no auth configured — generated password for this session", "password", password)
	}

	return password, secretKey, nil
}

func (a *App) Shutdown() error {
	slog.Info("shutting down")
	a.hub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed, forcing close", "err", err)
		if closeErr := a.srv.Close(); closeErr != nil {
			slog.Error("server close failed", "err", closeErr)
		}
	}

	if err := a.deps.teardownInfrastructure(); err != nil {
		slog.Error("teardown error", "err", err)
	}

	return nil
}
