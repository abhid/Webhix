package domain

import "time"

type Hook struct {
	ID        int64     `json:"id"`
	Token     string    `json:"token"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WebhookRequest struct {
	ID          int64     `json:"id"`
	HookID      int64     `json:"hookId"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Query       string    `json:"query,omitempty"`
	Headers     string    `json:"headers"`
	Body        []byte    `json:"body,omitempty"`
	RemoteAddr  string    `json:"remoteAddr,omitempty"`
	ContentType string    `json:"contentType,omitempty"`
	BodySize    int64     `json:"bodySize"`
	ReceivedAt  time.Time `json:"receivedAt"`
}

type HookResponse struct {
	ID         int64     `json:"id"`
	HookID     int64     `json:"hookId"`
	StatusCode int64     `json:"statusCode"`
	Headers    string    `json:"headers"`
	Body       []byte    `json:"body,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
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
