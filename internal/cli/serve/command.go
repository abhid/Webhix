package serve

import (
	"context"
	"log/slog"

	"github.com/GaIsBAX/Webhix/internal/app"
	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/spf13/cobra"
)

const (
	ServeGroup = "Serve"
	ServeTitle = ""
)

func NewCommand(ctx context.Context, cfg *config.Config) *cobra.Command {
	opts := DefaultOptions()

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start webhix server",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(cfg); err != nil {
				return err
			}

			app, err := app.New(ctx, cfg)
			if err != nil {
				slog.Error("init app", "err", err)
				return err
			}

			return run(ctx, app)
		},
	}

	RegisterFlags(cmd, cfg, &opts)

	return cmd
}

func run(ctx context.Context, app *app.App) error {
	if err := app.Start(ctx); err != nil {
		slog.Error("server", "err", err)
		return err
	}

	return nil
}
