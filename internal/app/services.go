package app

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/cli"
	"github.com/GaIsBAX/Webhix/internal/cli/serve"
	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/core"
	"github.com/GaIsBAX/Webhix/internal/domain"
)

type services struct {
	hook  *core.HookService
	serve *core.Serve
}

func newServices(repositories *repositories) *services {
	return &services{
		hook:  core.NewHookService(repositories.hook),
		serve: core.NewServe(repositories.serve),
	}
}

func NewVersionService() *core.Version {
	return core.NewVersion()
}

func Start(ctx context.Context, cfg *config.Config, args []string) error {
	versionService := NewVersionService()
	serveFactory := serve.ServiceFactoryFunc(func(ctx context.Context, cfg *config.Config) (serve.Service, domain.ServeStartFunc, error) {
		application, err := New(ctx, cfg)
		if err != nil {
			return nil, nil, err
		}

		return application.services.serve, application.Start, nil
	})

	return cli.Run(ctx, cfg, args, versionService, serveFactory)
}
