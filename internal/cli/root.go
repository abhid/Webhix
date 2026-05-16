package cli

import (
	"github.com/GaIsBAX/Webhix/internal/cli/serve"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "webhix",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	addGroup(cmd, serve.ServeGroup, serve.ServeTitle)

	cmd.AddCommand(serve.NewCommand())

	return cmd
}

func addGroup(cmd *cobra.Command, id, title string) {
	cmd.AddGroup(&cobra.Group{
		ID:    id,
		Title: title,
	})
}
