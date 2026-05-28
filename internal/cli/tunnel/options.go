package tunnel

import "fmt"

type Options struct {
	LocalPort   int
	RelayServer string
	AuthToken   string
	Subdomain   string
}

func DefaultOptions() Options {
	return Options{
		RelayServer: "wss://relay.webhix.online/tunnel",
	}
}

func (o *Options) Validate() error {
	if o.LocalPort <= 0 || o.LocalPort > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", o.LocalPort)
	}
	return nil
}
