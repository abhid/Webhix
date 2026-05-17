package app

import "github.com/GaIsBAX/Webhix/internal/core"

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
