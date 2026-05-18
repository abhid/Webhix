package core

import (
	"context"
	"errors"
	"time"

	"github.com/GaIsBAX/Webhix/internal/domain"
)

type ServeRepository interface {
	DeleteWebhookRequestsOlderThan(ctx context.Context, retention time.Duration) (int64, error)
	GetCountRequests(ctx context.Context) (int64, error)
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
			if _, err := cleanup(); err != nil {
				return 0, err
			}

		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
}

func (s *Serve) Run(ctx context.Context, opts domain.ServeRunOptions, start domain.ServeStartFunc, onRetentionError func(error)) error {
	s.StartRetentionCleaner(ctx, opts, onRetentionError)
	return start(ctx)
}

func (s *Serve) StartRetentionCleaner(ctx context.Context, opts domain.ServeRunOptions, onError func(error)) {
	if opts.Retention <= 0 || opts.ReadOnly {
		return
	}

	go func() {
		if _, err := s.RetentionCleaner(ctx, opts.Retention); err != nil && onError != nil {
			onError(err)
		}
	}()
}

func (s *Serve) RequestLimitGuard(ctx context.Context, limit int64) error {
	count, err := s.repo.GetCountRequests(ctx)
	if err != nil {
		return err
	}

	if count > limit {
		return errors.New("request: rate limit exceeded")
	}

	return nil
}
