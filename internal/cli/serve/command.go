package serve

import (
	"context"
	"log/slog"
	"time"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/spf13/cobra"
)

const (
	ServeGroup = "Serve"
	ServeTitle = ""
)

type Service interface {
	RunServe(ctx context.Context, retention time.Duration) error
}

type ServiceFactory func() (Service, error)

func NewCommand(ctx context.Context, cfg *config.Config, factory ServiceFactory) *cobra.Command {
	opts := DefaultOptions()

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start webhix server",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(cfg); err != nil {
				return err
			}

			service, err := factory()
			if err != nil {
				slog.Error("init app", "err", err)
				return err
			}

			return service.RunServe(ctx, opts.Retention)
		},
	}

	RegisterFlags(cmd, cfg, &opts)

	return cmd
}
