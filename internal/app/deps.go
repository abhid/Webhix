package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/GaIsBAX/Webhix/internal/config"
	"github.com/GaIsBAX/Webhix/internal/repos"
	"github.com/GaIsBAX/Webhix/internal/store"
)

type dependencies struct {
	db           *store.Database
	repositories *repositories
}

func newDependencies(ctx context.Context, cfg *config.Config) (*dependencies, error) {
	db, err := store.New(ctx, cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Migrate(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, errors.Join(
				fmt.Errorf("migrate database: %w", err),
				fmt.Errorf("close database after migration failure: %w", closeErr),
			)
		}

		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return &dependencies{
		db:           db,
		repositories: newRepositories(db),
	}, nil
}

type repositories struct {
	hook  *repos.HookRepository
	serve *repos.Serve
}

func newRepositories(db *store.Database) *repositories {
	return &repositories{
		hook:  repos.NewHookRepository(db.DB),
		serve: repos.NewServe(db.DB),
	}
}

func (d *dependencies) close() error {
	if d.db != nil {
		return d.db.Close()
	}

	return nil
}
