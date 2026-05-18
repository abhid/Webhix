package serve

import (
	"context"
	"log/slog"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/spf13/cobra"
)

const (
	ServeGroup = "Serve"
	ServeTitle = ""
)

type Service interface {
	Run(ctx context.Context, opts domain.ServeRunOptions, start domain.ServeStartFunc, onRetentionError func(error)) error
}

type ServiceFactory interface {
	New(ctx context.Context, cfg *config.Config) (Service, domain.ServeStartFunc, error)
}

type ServiceFactoryFunc func(ctx context.Context, cfg *config.Config) (Service, domain.ServeStartFunc, error)

func (f ServiceFactoryFunc) New(ctx context.Context, cfg *config.Config) (Service, domain.ServeStartFunc, error) {
	return f(ctx, cfg)
}

func NewCommand(ctx context.Context, cfg *config.Config, factory ServiceFactory) *cobra.Command {
	opts := DefaultOptions()

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start webhix server",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(cfg); err != nil {
				return err
			}

			service, start, err := factory.New(ctx, cfg)
			if err != nil {
				slog.Error("init app", "err", err)
				return err
			}

			return run(ctx, service, start, cfg, opts)
		},
	}

	RegisterFlags(cmd, cfg, &opts)

	return cmd
}
