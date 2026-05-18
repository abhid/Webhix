package serve

import (
	"context"
	"log/slog"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/core"
)

func run(ctx context.Context, service Service, start core.ServeStartFunc, cfg *config.Config, opts Options) error {
	return service.Run(ctx, core.ServeRunOptions{
		Retention: opts.Retention,
		ReadOnly:  cfg.ReadOnly,
	}, start, func(err error) {
		slog.Error("retention cleaner", "err", err)
	})
}
