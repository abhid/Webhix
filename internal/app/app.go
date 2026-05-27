package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/core"
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
	if a.config.TLSDomain != "" {
		return a.startTLS(ctx)
	}
	return a.startPlain(ctx)
}

func (a *App) startPlain(ctx context.Context) error {
	slog.Info("webhix started", "addr", a.config.Addr, "base_url", a.config.BaseURL)
	return a.run(ctx, a.server.ListenAndServe)
}

func (a *App) startTLS(ctx context.Context) error {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(a.config.TLSDomain),
		Cache:      autocert.DirCache(a.config.TLSCacheDir),
	}

	a.server.Addr = ":443"
	a.server.TLSConfig = m.TLSConfig()

	redirect := &http.Server{
		Addr: ":80",
		Handler: m.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := &url.URL{
				Scheme:   "https",
				Host:     a.config.TLSDomain,
				Path:     r.URL.Path,
				RawQuery: r.URL.RawQuery,
			}
			w.Header().Set("Location", target.String())
			w.WriteHeader(http.StatusMovedPermanently)
		})),
		ReadHeaderTimeout: readHeaderTimeout,
	}
	defer shutdownServer(redirect, "redirect server")

	serverErr := make(chan error, 2)

	go func() {
		err := redirect.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("redirect server (:80): %w", err)
		}
	}()

	go func() {
		err := a.server.ListenAndServeTLS("", "")
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serverErr <- err
	}()

	slog.Info("webhix started (TLS)", "domain", a.config.TLSDomain, "base_url", a.config.BaseURL)

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return a.Shutdown(ctx)
	}
}

func (a *App) run(ctx context.Context, listen func() error) error {
	serverErr := make(chan error, 1)
	go func() {
		err := listen()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serverErr <- err
	}()

	select {
	case err := <-serverErr:
		return err

	case <-ctx.Done():
		return a.Shutdown(ctx)
	}
}

func shutdownServer(s *http.Server, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("shutdown", "server", name, "err", err)
	}
}

func (a *App) RunServe(ctx context.Context, retention time.Duration) error {
	a.deps.services.serve.StartRetentionCleaner(
		ctx,
		core.ServeRunOptions{Retention: retention, ReadOnly: a.config.ReadOnly},
		func(err error) { slog.Error("retention cleaner", "err", err) },
	)
	return a.Start(ctx)
}

func (a *App) Shutdown(ctx context.Context) error {
	slog.Info("shutting down")
	a.deps.infra.hub.Close()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
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
