package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/core"
	"github.com/GaIsBAX/Webhix/internal/server"
	"github.com/GaIsBAX/Webhix/internal/store"
)

type App struct {
	srv  *http.Server
	cfg  *config.Config
	deps *Deps
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	mux := http.NewServeMux()

	deps, err := NewDeps(ctx, cfg)
	if err != nil {
		return nil, err
	}

	hookRepository := store.NewHookRepository(deps.DB.DB)
	hookService := core.NewHookService(hookRepository)
	hookHandler := server.NewHookHandler(mux, hookService, cfg.BaseURL)

	hookHandler.RegisterRoutes()

	return &App{
		srv:  &http.Server{Addr: cfg.Addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second},
		cfg:  cfg,
		deps: deps,
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

func (a *App) Shutdown() error {
	slog.Info("shutting down")

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
