package cli

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/cli/serve"
	"github.com/GaIsBAX/Webhix/internal/config"
)

func Run(
	ctx context.Context,
	cfg *config.Config,
	args []string,
	serveFactory serve.ServiceFactory,
) error {
	root := NewRootCommand(ctx, cfg, serveFactory)
	root.SetArgs(args)

	return root.Execute()
}
