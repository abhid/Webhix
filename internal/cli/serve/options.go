package serve

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/GaIsBAX/Webhix/internal/config"
)

type Options struct {
	Addr    string
	DBPath  string
	BaseURL string

	MaxBodySize    int
	TrustedProxies []string
}

func DefaultOptions() Options {
	return Options{}
}

func NewOptions(cfg *config.Config) Options {
	opts := DefaultOptions()

	if cfg != nil {
		opts.Addr = cfg.Addr
		opts.DBPath = cfg.DBPath
		opts.BaseURL = cfg.BaseURL
		opts.TrustedProxies = cfg.TrustedProxies
	}

	return opts
}

func (o *Options) Validate() error {
	if strings.TrimSpace(o.Addr) == "" {
		return fmt.Errorf("addr cannot be empty")
	}

	u, err := url.Parse(o.BaseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid base URL:\n  got:  %s\n  want: https://hooks.example.com", o.BaseURL)
	}

	return nil
}

func (o *Options) Apply(cfg *config.Config) {
	cfg.Addr = o.Addr
	cfg.DBPath = o.DBPath
	cfg.BaseURL = o.BaseURL
	cfg.TrustedProxies = o.TrustedProxies
}
