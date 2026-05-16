package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func New(cfg *config.Config) (*App, error) {
	mux := http.NewServeMux()

	deps, err := NewDeps(cfg)
	if err != nil {
		return nil, err
	}

	hookRepository := store.NewHookRepository(deps.DB.DB)
	hookService := core.NewHookService(hookRepository)
	hookHandler := server.NewHookHandler(mux, hookService, cfg.BaseURL)

	hookHandler.RegisterRoutes()

	return &App{
		srv:  &http.Server{Addr: cfg.Addr, Handler: mux},
		cfg:  cfg,
		deps: deps,
	}, nil
}

func (a *App) Start() error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("webhix started", "addr", a.cfg.Addr, "base_url", a.cfg.BaseURL)
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errCh:
		return err
	case <-quit:
		return a.Shutdown()
	}
}

func (a *App) Shutdown() error {
	slog.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	}

	if err := a.deps.teardownInfrastructure(); err != nil {
		slog.Error("teardown error", "err", err)
	}

	return nil
}
