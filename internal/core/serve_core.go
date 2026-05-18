package core

import (
	"context"
	"time"
)

type ServeRepository interface {
	DeleteWebhookRequestsOlderThan(ctx context.Context, retention time.Duration) (int64, error)
}

type Serve struct {
	repo ServeRepository
}

func NewServe(repo ServeRepository) *Serve {
	return &Serve{
		repo: repo,
	}
}

func (s *Serve) RetentionCleaner(ctx context.Context, retention time.Duration) (int64, error) {
	cleanup := func() (int64, error) {
		return s.repo.DeleteWebhookRequestsOlderThan(ctx, retention)
	}

	_, err := cleanup()
	if err != nil {
		return 0, err
	}

	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanup()

		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
}
