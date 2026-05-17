package forward

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/spf13/cobra"
)

func NewCommand(ctx context.Context, cfg *config.Config) *cobra.Command {
	opts := DefaultOptions()
	opts.AuthToken = cfg.SecretKey

	cmd := &cobra.Command{
		Use:   "forward <token>",
		Short: "Forward incoming webhook requests to a local server",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Token = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			return run(ctx, opts)
		},
	}

	RegisterFlags(cmd, &opts)

	return cmd
}
