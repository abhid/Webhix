package version

import "fmt"

const (
	outputText = "text"
	outputJSON = "json"
	outputYAML = "yaml"
)

type Options struct {
	Output string
}

func NewOptions() *Options {
	return &Options{
		Output: outputText,
	}
}

func (o Options) Validate() error {
	switch o.Output {
	case outputText, outputJSON, outputYAML:
		return nil
	default:
		return fmt.Errorf("unsupported output format %q (want text, json, or yaml)", o.Output)
	}
}
