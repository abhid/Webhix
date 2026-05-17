package serve

import (
	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/spf13/cobra"
)

const (
	flagAddr    = "addr"
	flagDBPath  = "db-path"
	flagBaseURL = "base-url"

	flagPassword  = "password"
	flagSecretKey = "secret-key"

	flagMaxBodySize    = "max-body-size"
	flagTrustedProxies = "trusted-proxies"
)

func RegisterFlags(cmd *cobra.Command, cfg *config.Config, opt *Options) {
	flags := cmd.Flags()

	flags.StringVarP(&cfg.Addr, flagAddr, "a", cfg.Addr, "address to listen on")
	flags.StringVarP(&cfg.DBPath, flagDBPath, "d", cfg.DBPath, "path to SQLite database directory")

	flags.StringVar(&cfg.BaseURL, flagBaseURL, cfg.BaseURL, "public base URL used for endpoint links")
	flags.StringVar(&cfg.Password, flagPassword, cfg.Password, "basic auth password (env: WEBHIX_PASSWORD)")
	flags.StringVar(&cfg.SecretKey, flagSecretKey, cfg.SecretKey, "API secret key via X-Webhix-Key or Bearer (env: WEBHIX_SECRET_KEY)")
	flags.Int64Var(&cfg.MaxBodySize, flagMaxBodySize, cfg.MaxBodySize, "maximum webhook request body size in bytes")
	flags.StringSliceVar(&cfg.TrustedProxies, flagTrustedProxies, cfg.TrustedProxies, "trusted proxy CIDRs")
}
