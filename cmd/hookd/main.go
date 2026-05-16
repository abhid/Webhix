package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := application.Start(ctx); err != nil {
		slog.Error("server", "err", err)
		os.Exit(1)
	}
}
