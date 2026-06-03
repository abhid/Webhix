package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Options struct {
	Server    string
	AuthToken string
}

type apiResponse struct {
	Success bool            `json:"success"`
	Body    json.RawMessage `json:"body"`
	Error   *apiError       `json:"error"`
}

type apiError struct {
	Message string `json:"message"`
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

func DefaultOptions() Options {
	return Options{
		Server: "http://localhost:8080",
	}
}

func (o *Options) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, o.Server+path, body)
	if err != nil {
		return nil, err
	}
	if o.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+o.AuthToken)
	}
	return req, nil
}

func apiGet(ctx context.Context, opts *Options, path string, out any) error {
	req, err := opts.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("close response body", "err", err)
		}
	}()

	var ar apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return err
	}
	if !ar.Success {
		if ar.Error != nil {
			return errors.New(ar.Error.Message)
		}
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}
	return json.Unmarshal(ar.Body, out)
}

func maskToken(token string) string {
	if len(token) <= 14 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
