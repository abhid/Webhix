package serve

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/GaIsBAX/Webhix/internal/config"
)

type Options struct{}

func DefaultOptions() Options {
	return Options{}
}

func (o *Options) Validate(cfg *config.Config) error {
	if strings.TrimSpace(cfg.Addr) == "" {
		return fmt.Errorf("addr cannot be empty")
	}

	u, err := url.Parse(cfg.BaseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid base URL:\n  got:  %s\n  want: https://hooks.example.com", cfg.BaseURL)
	}

	if cfg.MaxBodySize <= 0 {
		return fmt.Errorf("max body size must be greater than 0")
	}

	return nil
}
