package forward

import (
	"fmt"
	"strings"
)

type Options struct {
	Token       string
	To          string
	Server      string
	AuthToken   string
	RewriteHost bool
}

func DefaultOptions() Options {
	return Options{
		Server: "http://localhost:8080",
	}
}

func (o *Options) Validate() error {
	if strings.TrimSpace(o.Token) == "" {
		return fmt.Errorf("token is required")
	}
	if strings.TrimSpace(o.To) == "" {
		return fmt.Errorf("--to is required")
	}
	return nil
}
