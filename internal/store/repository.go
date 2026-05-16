package store

import (
	"context"
	"database/sql"

	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/GaIsBAX/Webhix/internal/store/sqlc"
)

type HookRepository struct {
	q *sqlc.Queries
}

func NewHookRepository(db sqlc.DBTX) *HookRepository {
	return &HookRepository{
		q: sqlc.New(db),
	}
}

func (r *HookRepository) CreateHook(ctx context.Context, token string) (domain.Hook, error) {
	hook, err := r.q.CreateHook(ctx, sqlc.CreateHookParams{
		Token: token,
		Name:  sql.NullString{},
	})
	if err != nil {
		return domain.Hook{}, err
	}

	return toDomainHook(hook), nil
}

func (r *HookRepository) GetHookByToken(ctx context.Context, token string) (domain.Hook, error) {
	hook, err := r.q.GetHookByToken(ctx, token)
	if err != nil {
		return domain.Hook{}, err
	}

	return toDomainHook(hook), nil
}

func toDomainHook(hook sqlc.Hook) domain.Hook {
	return domain.Hook{
		ID:        hook.ID,
		Token:     hook.Token,
		Name:      hook.Name.String,
		CreatedAt: hook.CreatedAt,
		UpdatedAt: hook.UpdatedAt,
	}
}
