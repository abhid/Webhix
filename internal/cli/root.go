package cli

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/cli/serve"
	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCommand(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "webhix",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	addGroup(cmd, serve.ServeGroup, serve.ServeTitle)

	cmd.AddCommand(serve.NewCommand(ctx, cfg))

	return cmd
}

func addGroup(cmd *cobra.Command, id, title string) {
	cmd.AddGroup(&cobra.Group{
		ID:    id,
		Title: title,
	})
}
