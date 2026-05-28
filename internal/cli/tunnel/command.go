package tunnel

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/spf13/cobra"
)

func NewCommand(ctx context.Context, cfg *config.Config) *cobra.Command {
	opts := DefaultOptions()

	if cfg.SecretKey != "" {
		opts.AuthToken = cfg.SecretKey
	}

	cmd := &cobra.Command{
		Use:   "tunnel <port>",
		Short: "Expose a local port via a public webhix.online URL",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid port %q: must be a number", args[0])
			}
			opts.LocalPort = port

			if err := opts.Validate(); err != nil {
				return err
			}

			// TODO(v0.3): implement tunnel runner
			// See docs/tunnel-protocol.md for the WebSocket protocol spec.
			cmd.Println("webhix tunnel is not yet implemented.")

			_ = ctx
			return nil
		},
	}

	RegisterFlags(cmd, &opts)

	return cmd
}
