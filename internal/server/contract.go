package server

import (
	"encoding/json"
	"time"

	"github.com/GaIsBAX/Webhix/internal/domain"
	"github.com/GaIsBAX/Webhix/pkg"
)

type CreateEndpointResponseContract struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name,omitempty"`
	URL   string `json:"url"`
}

type WebhookRequestContract struct {
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

type ResponseContract struct {
	Success   bool            `json:"success"`
	Data      json.RawMessage `json:"body,omitempty"`
	Error     *ErrorContract  `json:"error,omitempty"`
	RequestID string          `json:"requestID"`
	Timestamp time.Time       `json:"timestamp"`
}

type ErrorContract struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Details []ErrorDetailContract `json:"details,omitempty"`
}

type ErrorDetailContract struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func NewSuccessResponseContract(data []byte) *ResponseContract {
	return &ResponseContract{
		Success:   true,
		Data:      json.RawMessage(data),
		RequestID: pkg.GeneratePrefixedString("re"),
		Timestamp: time.Now().UTC(),
	}
}

func NewErrorResponseContract(apiErr ErrorContract) *ResponseContract {
	return &ResponseContract{
		Success:   false,
		Error:     &apiErr,
		RequestID: pkg.GeneratePrefixedString("re"),
		Timestamp: time.Now().UTC(),
	}
}

func WithDetails(err ErrorContract, details ...ErrorDetailContract) ErrorContract {
	err.Details = details
	return err
}

func toWebhookRequestContract(req domain.WebhookRequest) WebhookRequestContract {
	return WebhookRequestContract{
		ID:          req.ID,
		HookID:      req.HookID,
		Method:      req.Method,
		Path:        req.Path,
		Query:       req.Query,
		Headers:     req.Headers,
		Body:        req.Body,
		RemoteAddr:  req.RemoteAddr,
		ContentType: req.ContentType,
		BodySize:    req.BodySize,
		ReceivedAt:  req.ReceivedAt,
	}
}
