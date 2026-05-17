package forward

import "github.com/spf13/cobra"

const (
	flagTo          = "to"
	flagServer      = "server"
	flagAuthToken   = "auth-token"
	flagRewriteHost = "rewrite-host"
)

func RegisterFlags(cmd *cobra.Command, opt *Options) {
	flags := cmd.Flags()

	flags.StringVar(&opt.To, flagTo, opt.To, "target URL to forward requests to")
	flags.StringVar(&opt.Server, flagServer, opt.Server, "Webhix server URL")
	flags.StringVar(&opt.AuthToken, flagAuthToken, opt.AuthToken, "auth token for Webhix server (env: WEBHIX_SECRET_KEY)")
	flags.BoolVar(&opt.RewriteHost, flagRewriteHost, opt.RewriteHost, "rewrite Host header to match target")

	if err := cmd.MarkFlagRequired(flagTo); err != nil {
		panic(err)
	}
}
