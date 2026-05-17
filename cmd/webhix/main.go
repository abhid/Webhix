package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/GaIsBAX/Webhix/internal/cli"
	"github.com/GaIsBAX/Webhix/internal/config"
	_ "github.com/GaIsBAX/Webhix/pkg"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	if err := cli.Run(ctx, cfg, os.Args[1:]); err != nil {
		slog.Error("app run", "err", err)
		os.Exit(1)
	}
}
