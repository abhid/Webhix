package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/GaIsBAX/Webhix/internal/core"
	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type versionInfoContract struct {
	Version string `json:"version" yaml:"version"`
	Commit  string `json:"commit"  yaml:"commit"`
	Built   string `json:"built"   yaml:"built"`
	Go      string `json:"go"      yaml:"go"`
}

func NewCommand(ctx context.Context) *cobra.Command {
	service := core.NewVersion()
	opts := NewOptions()

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			return print(cmd.OutOrStdout(), service.Info(), opts.Output)
		},
	}

	RegisterFlags(cmd, opts)

	return cmd
}

func print(w io.Writer, info domain.VersionInfo, output string) error {
	switch output {
	case outputJSON:
		encoder := json.NewEncoder(w)
		return encoder.Encode(toContract(info))

	case outputYAML:
		encoder := yaml.NewEncoder(w)
		if err := encoder.Encode(toContract(info)); err != nil {
			return err
		}
		return encoder.Close()

	default:
		_, err := fmt.Fprintf(
			w,
			"webhix %s\ncommit: %s\nbuilt: %s\ngo: %s\n",
			info.Version,
			info.Commit,
			info.Built,
			info.Go,
		)
		return err
	}
}

func toContract(info domain.VersionInfo) versionInfoContract {
	return versionInfoContract{
		Version: info.Version,
		Commit:  info.Commit,
		Built:   info.Built,
		Go:      info.Go,
	}
}
