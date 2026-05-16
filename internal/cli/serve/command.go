package serve

import "github.com/spf13/cobra"

const (
	ServeGroup = "Serve"
	ServeTitle = ""
)

func NewCommand() *cobra.Command {
	opts := DefaultOptions()

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start webhix server",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			return run(opts)
		},
	}

	return cmd
}

func run(opts *Options) error {
	panic("impl me")
}
