package store

import (
	"context"
	"database/sql"
	"errors"

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
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Hook{}, domain.ErrNotFound
		}
		return domain.Hook{}, err
	}

	return toDomainHook(hook), nil
}

func (r *HookRepository) CreateWebhookRequest(ctx context.Context, params domain.CreateWebhookRequestParams) (domain.WebhookRequest, error) {

	req, err := r.q.CreateWebhookRequest(ctx, sqlc.CreateWebhookRequestParams{
		HookID:      params.HookID,
		Method:      params.Method,
		Path:        params.Path,
		Query:       params.Query,
		Headers:     params.Headers,
		Body:        params.Body,
		RemoteAddr:  sql.NullString{String: params.RemoteAddr, Valid: params.RemoteAddr != ""},
		ContentType: sql.NullString{String: params.ContentType, Valid: params.ContentType != ""},
		BodySize:    params.BodySize,
	})
	if err != nil {
		return domain.WebhookRequest{}, err
	}

	return toDomainWebhookRequest(req), nil
}

func (r *HookRepository) ListWebhookRequests(ctx context.Context, hookID int64) ([]domain.WebhookRequest, error) {
	rows, err := r.q.ListWebhookRequestsByHookID(ctx, hookID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.WebhookRequest, len(rows))
	for i, row := range rows {
		result[i] = toDomainWebhookRequest(row)
	}

	return result, nil
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

func toDomainWebhookRequest(req sqlc.WebhookRequest) domain.WebhookRequest {
	return domain.WebhookRequest{
		ID:          req.ID,
		HookID:      req.HookID,
		Method:      req.Method,
		Path:        req.Path,
		Query:       req.Query,
		Headers:     req.Headers,
		Body:        req.Body,
		RemoteAddr:  req.RemoteAddr.String,
		ContentType: req.ContentType.String,
		BodySize:    req.BodySize,
		ReceivedAt:  req.ReceivedAt,
	}
}
