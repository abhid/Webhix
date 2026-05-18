package version

import "github.com/spf13/cobra"

const flagOutput = "output"

func RegisterFlags(cmd *cobra.Command, opts *Options) {
	cmd.Flags().StringVar(&opts.Output, flagOutput, opts.Output, "output format: text, json, or yaml")
}
