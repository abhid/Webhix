package tunnel

import "github.com/spf13/cobra"

const (
	flagRelay     = "relay"
	flagAuthToken = "auth-token"
	flagSubdomain = "subdomain"
)

func RegisterFlags(cmd *cobra.Command, opt *Options) {
	flags := cmd.Flags()

	flags.StringVar(&opt.RelayServer, flagRelay, opt.RelayServer,
		"relay server WebSocket URL (env: WEBHIX_RELAY_SERVER)")
	flags.StringVar(&opt.AuthToken, flagAuthToken, opt.AuthToken,
		"Pro auth token from webhix.online (env: WEBHIX_TUNNEL_TOKEN)")
	flags.StringVar(&opt.Subdomain, flagSubdomain, opt.Subdomain,
		"reserved subdomain to request (Pro only)")
}
