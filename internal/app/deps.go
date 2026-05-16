package app

import (
	"errors"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/store"
)

type Deps struct {
	DB  *store.Database
	cfg *config.Config
}

func NewDeps(cfg *config.Config) (*Deps, error) {
	deps := &Deps{
		cfg: cfg,
	}

	if err := deps.setupInfrastructure(); err != nil {
		return nil, err
	}

	return deps, nil
}

func (d *Deps) setupInfrastructure() error {
	var errs []error

	database, err := store.New(d.cfg.DBPath)
	if err != nil {
		errs = append(errs, err)
	}

	d.DB = database

	if d.DB != nil {
		if err := d.DB.Migrate(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (d *Deps) teardownInfrastructure() error {
	var errs []error

	if d.DB != nil {
		if err := d.DB.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
