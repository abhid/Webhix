package serve

import "github.com/spf13/cobra"

const (
	flagAddr    = "addr"
	flagDBPath  = "db-path"
	flagBaseURL = "base-url"

	flagTrustedProxies = "trusted-proxies"
)

func RegisterFlags(cmd *cobra.Command, opt *Options) {
	flags := cmd.Flags()

	flags.StringVarP(&opt.Addr, flagAddr, "a", opt.Addr, "TODO")
	flags.StringVarP(&opt.DBPath, flagDBPath, "d", opt.DBPath, "TODO")

	flags.StringVar(&opt.BaseURL, flagBaseURL, opt.BaseURL, "TODO")
	flags.StringSliceVar(&opt.TrustedProxies, flagTrustedProxies, opt.TrustedProxies, "trusted proxy CIDRs")
}
