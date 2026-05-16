package core

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/GaIsBAX/Webhix/pkg"
)

type HookRepository interface {
	CreateHook(ctx context.Context, token string) (domain.Hook, error)
	GetHookByToken(ctx context.Context, token string) (domain.Hook, error)
	CreateWebhookRequest(ctx context.Context, params domain.CreateWebhookRequestParams) (domain.WebhookRequest, error)
	ListWebhookRequests(ctx context.Context, hookID int64) ([]domain.WebhookRequest, error)
}

type HookService struct {
	repo HookRepository
}

func NewHookService(repo HookRepository) *HookService {
	return &HookService{
		repo: repo,
	}
}

func (s *HookService) CreateHook(ctx context.Context, token string) (domain.Hook, error) {
	if token == "" {
		token = pkg.GeneratePrefixedString("ho")
	}

	return s.repo.CreateHook(ctx, token)
}

func (s *HookService) ReceiveWebhook(ctx context.Context, token string, params domain.CreateWebhookRequestParams) (domain.WebhookRequest, error) {
	hook, err := s.repo.GetHookByToken(ctx, token)
	if err != nil {
		return domain.WebhookRequest{}, err
	}

	params.HookID = hook.ID

	return s.repo.CreateWebhookRequest(ctx, params)
}

func (s *HookService) ListWebhookRequests(ctx context.Context, token string) ([]domain.WebhookRequest, error) {
	hook, err := s.repo.GetHookByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.repo.ListWebhookRequests(ctx, hook.ID)
}
