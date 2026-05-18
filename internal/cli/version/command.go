package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Service interface {
	Info() domain.VersionInfo
}

func NewCommand(ctx context.Context, service Service) *cobra.Command {
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
		return encoder.Encode(info)
	
	case outputYAML:
		encoder := yaml.NewEncoder(w)
		defer encoder.Close()
		return encoder.Encode(info)
	
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
