package serve

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/app"
)

func run(ctx context.Context, application *app.App, opts Options) error {
	return application.RunServe(ctx, app.ServeOptions{
		Retention: opts.Retention,
	})
}
