package serve

import "github.com/spf13/cobra"

const (
	flagAddr    = "addr"
	flagDBPath  = "db-path"
	flagBaseURL = "base-url"

	flagPassword  = "password"
	flagSecretKey = "secret-key"

	flagTrustedProxies = "trusted-proxies"
)

func RegisterFlags(cmd *cobra.Command, opt *Options) {
	flags := cmd.Flags()

	flags.StringVarP(&opt.Addr, flagAddr, "a", opt.Addr, "address to listen on")
	flags.StringVarP(&opt.DBPath, flagDBPath, "d", opt.DBPath, "path to SQLite database directory")

	flags.StringVar(&opt.BaseURL, flagBaseURL, opt.BaseURL, "public base URL used for endpoint links")
	flags.StringVar(&opt.Password, flagPassword, opt.Password, "basic auth password (env: WEBHIX_PASSWORD)")
	flags.StringVar(&opt.SecretKey, flagSecretKey, opt.SecretKey, "API secret key via X-Webhix-Key or Bearer (env: WEBHIX_SECRET_KEY)")
	flags.StringSliceVar(&opt.TrustedProxies, flagTrustedProxies, opt.TrustedProxies, "trusted proxy CIDRs")
}
