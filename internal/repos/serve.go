package repos

import (
	"context"
	"time"

	"github.com/GaIsBAX/Webhix/internal/store/sqlc"
)

type Serve struct {
	q *sqlc.Queries
}

func NewServe(db sqlc.DBTX) *Serve {
	return &Serve{
		q: sqlc.New(db),
	}
}

func (r *Serve) DeleteWebhookRequestsOlderThan(ctx context.Context, retention time.Duration) (int64, error) {
	res, err := r.q.DeleteWebhookRequestsOlderThan(ctx, retention)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}
