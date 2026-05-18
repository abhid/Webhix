package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/server"
	"github.com/GaIsBAX/Webhix/internal/server/middleware"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type App struct {
	server *http.Server

	config *config.Config
	deps   *dependencies
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	deps, err := newDependencies(ctx, cfg)
	if err != nil {
		return nil, err
	}

	handler, err := newHTTPHandler(deps.mux, cfg)
	if err != nil {
		if closeErr := deps.close(); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &App{
		server: server,
		config: cfg,
		deps:   deps,
	}, nil
}

func newHTTPHandler(mux *http.ServeMux, cfg *config.Config) (http.Handler, error) {
	var middlewares []func(http.Handler) http.Handler

	if len(cfg.TrustedProxies) > 0 {
		trustedProxies := middleware.NewTrustedProxies(cfg.TrustedProxies)
		if trustedProxies == nil {
			return nil, ErrInvalidTrustedProxies
		}
		middlewares = append(middlewares, trustedProxies.BehindProxy)
	}

	if cfg.Password != "" || cfg.SecretKey != "" {
		auth := middleware.NewAuth(cfg.Password, cfg.SecretKey)
		middlewares = append(middlewares, auth.Protect)
	}

	return server.Chain(mux, middlewares...), nil
}

func (a *App) Start(ctx context.Context) error {
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("webhix started", "addr", a.config.Addr, "base_url", a.config.BaseURL)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err

	case <-ctx.Done():
		return a.Shutdown(ctx)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	slog.Info("shutting down")
	a.deps.infra.hub.Close()

	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed, forcing close", "err", err)
		if closeErr := a.server.Close(); closeErr != nil {
			slog.Error("server close failed", "err", closeErr)
		}
	}

	if err := a.deps.close(); err != nil {
		slog.Error("teardown error", "err", err)
		return err
	}

	return nil
}
