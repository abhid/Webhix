package serve

import (
	"fmt"
	"strings"
)

type Options struct {
	Addr string
}

func DefaultOptions() *Options {
	return &Options{
		Addr: ":8000",
	}
}

func (o *Options) Validate() error {
	if strings.TrimSpace(o.Addr) == "" {
		return fmt.Errorf("addr cannot be empty")
	}

	return nil
}
