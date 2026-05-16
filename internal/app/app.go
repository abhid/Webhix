package app

import (
	"log"
	"net/http"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/core"
	"github.com/GaIsBAX/Webhix/internal/server"
	"github.com/GaIsBAX/Webhix/internal/store"
)

type App struct {
	mux *http.ServeMux
	cfg *config.Config
}

func New(cfg *config.Config) (*App, error) {
	mux := http.NewServeMux()

	deps, err := NewDeps(cfg)
	if err != nil {
		return nil, err
	}

	hookRepository := store.NewHookRepository(deps.DB.DB)
	hookService := core.NewHookService(hookRepository)
	hookHandler := server.NewHookHandler(mux, hookService)

	hookHandler.RegisterRoutes()

	return &App{
		mux: mux,
		cfg: cfg,
	}, nil
}

func (a *App) Start() error {
	log.Printf("starting webhix server on %s", a.cfg.WebHixAddr)
	return http.ListenAndServe(a.cfg.WebHixAddr, a.mux)
}

// TODO: Gracefull shutdown
