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
	GetHookResponse(ctx context.Context, hookID int64) (domain.HookResponse, error)
	UpsertHookResponse(ctx context.Context, hookID int64, params domain.UpsertHookResponseParams) (domain.HookResponse, error)
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

func (s *HookService) ReceiveWebhook(ctx context.Context, token string, params domain.CreateWebhookRequestParams) (domain.WebhookRequest, domain.HookResponse, error) {
	hook, err := s.repo.GetHookByToken(ctx, token)
	if err != nil {
		return domain.WebhookRequest{}, domain.HookResponse{}, err
	}

	params.HookID = hook.ID

	req, err := s.repo.CreateWebhookRequest(ctx, params)
	if err != nil {
		return domain.WebhookRequest{}, domain.HookResponse{}, err
	}

	resp, _ := s.repo.GetHookResponse(ctx, hook.ID)

	return req, resp, nil
}

func (s *HookService) ListWebhookRequests(ctx context.Context, token string) ([]domain.WebhookRequest, error) {
	hook, err := s.repo.GetHookByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.repo.ListWebhookRequests(ctx, hook.ID)
}

func (s *HookService) GetHookResponse(ctx context.Context, token string) (domain.HookResponse, error) {
	hook, err := s.repo.GetHookByToken(ctx, token)
	if err != nil {
		return domain.HookResponse{}, err
	}

	return s.repo.GetHookResponse(ctx, hook.ID)
}

func (s *HookService) SetHookResponse(ctx context.Context, token string, params domain.UpsertHookResponseParams) (domain.HookResponse, error) {
	hook, err := s.repo.GetHookByToken(ctx, token)
	if err != nil {
		return domain.HookResponse{}, err
	}

	return s.repo.UpsertHookResponse(ctx, hook.ID, params)
}
