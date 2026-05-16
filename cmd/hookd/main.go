package main

import (
	"log/slog"
	"os"

	"github.com/GaIsBAX/Webhix/internal/app"
	"github.com/GaIsBAX/Webhix/internal/config"
	_ "github.com/GaIsBAX/Webhix/pkg"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("init app", "err", err)
		os.Exit(1)
	}

	if err := application.Start(); err != nil {
		slog.Error("server", "err", err)
		os.Exit(1)
	}
}
