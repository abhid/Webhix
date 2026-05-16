package server

import (
	"encoding/json"
	"time"

	"github.com/GaIsBAX/Webhix/pkg"
)

type CreateEndpointResponseContract struct {
	ID        int64  `json:"id"`
	Token     string `json:"token"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url"`
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
