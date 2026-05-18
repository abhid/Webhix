package cli

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/cli/serve"
	"github.com/GaIsBAX/Webhix/internal/cli/version"
	"github.com/GaIsBAX/Webhix/internal/config"
)

func Run(
	ctx context.Context,
	cfg *config.Config,
	args []string,
	versionService version.Service,
	serveFactory serve.ServiceFactory,
) error {
	root := NewRootCommand(ctx, cfg, versionService, serveFactory)
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
