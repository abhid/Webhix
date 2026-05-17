package cli

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/config"
)

func Run(ctx context.Context, cfg *config.Config, args []string) error {
	root := NewRootCommand(ctx, cfg)
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
