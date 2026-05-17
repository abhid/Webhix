package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/hub"
)

const shutdownTimeout = 10 * time.Second

type App struct {
	server   *http.Server
	config   *config.Config
	deps     *dependencies
	events   *hub.Hub
	services *services
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	deps, err := newDependencies(ctx, cfg)
	if err != nil {
		return nil, err
	}

	services := newServices(deps.repositories)
	events := hub.New()

	mux, err := newMux(cfg, services, events)
	if err != nil {
		return nil, err
	}

	handler, err := newHTTPHandler(cfg, mux)
	if err != nil {
		return nil, err
	}

	return &App{
		server:   newHTTPServer(cfg, handler),
		config:   cfg,
		deps:     deps,
		events:   events,
		services: services,
	}, nil
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
		return a.Shutdown()
	}
}

func (a *App) Shutdown() error {
	slog.Info("shutting down")
	a.events.Close()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed, forcing close", "err", err)
		if closeErr := a.server.Close(); closeErr != nil {
			slog.Error("server close failed", "err", closeErr)
		}
	}

	if err := a.deps.close(); err != nil {
		slog.Error("teardown error", "err", err)
	}

	return nil
}

type ServeOptions struct {
	Retention time.Duration
}

func (a *App) RunServe(ctx context.Context, opts ServeOptions) error {
	go func() {
		if _, err := a.services.serve.RetentionCleaner(ctx, opts.Retention); err != nil {
			slog.Error("retention cleaner", "err", err)
		}
	}()

	return a.Start(ctx)
}
