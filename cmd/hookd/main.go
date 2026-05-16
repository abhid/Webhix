package main

import (
	"github.com/GaIsBAX/Webhix/internal/app"
	"github.com/GaIsBAX/Webhix/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		// TODO: пока будет паника
		panic(err)
	}

	application, err := app.New(cfg)
	if err != nil {
		// TODO: пока будет паника
		panic(err)
	}

	if err := application.Start(); err != nil {
		// TODO: пока будет паника
		panic(err)
	}
}
