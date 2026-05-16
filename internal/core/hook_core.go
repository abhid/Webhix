package core

import (
	"context"

	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/GaIsBAX/Webhix/pkg"
)

type HookRepository interface {
	CreateHook(ctx context.Context, token string) (domain.Hook, error)
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
