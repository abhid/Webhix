package domain

import (
	"fmt"
	"time"
)

type Hook struct {
	ID        int64
	Token     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WebhookRequest struct {
	ID          int64
	HookID      int64
	Method      string
	Path        string
	Query       string
	Headers     string
	Body        []byte
	RemoteAddr  string
	ContentType string
	BodySize    int64
	ReceivedAt  time.Time
}

type HookResponse struct {
	ID         int64
	HookID     int64
	StatusCode int64
	Headers    map[string]string
	Body       []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Page struct {
	Limit  int64
	Offset int64
}

const (
	DefaultPageLimit int64 = 100
	MaxPageLimit     int64 = 1000
)

// Normalize clamps limit/offset to safe bounds.
func (p Page) Normalize() Page {
	if p.Limit <= 0 {
		p.Limit = DefaultPageLimit
	}
	if p.Limit > MaxPageLimit {
		p.Limit = MaxPageLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}

type WebhookRequestPage struct {
	Requests []WebhookRequest
	Total    int64
	Limit    int64
	Offset   int64
}

type CreateWebhookRequestParams struct {
	HookID      int64
	Method      string
	Path        string
	Query       string
	Headers     string
	Body        []byte
	RemoteAddr  string
	ContentType string
	BodySize    int64
}

type UpsertHookResponseParams struct {
	StatusCode int64
	Headers    map[string]string
	Body       []byte
}

func (p UpsertHookResponseParams) Validate() error {
	if p.StatusCode < 100 || p.StatusCode > 599 {
		return fmt.Errorf("statusCode must be between 100 and 599")
	}
	return nil
}
