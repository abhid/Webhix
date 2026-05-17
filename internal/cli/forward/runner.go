package forward

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type webhookEvent struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Headers string `json:"headers"`
	Body    []byte `json:"body"`
}

func run(ctx context.Context, opts Options) error {
	targetURL := opts.To
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	sseURL := fmt.Sprintf("%s/api/endpoints/%s/events", opts.Server, opts.Token)

	slog.Info("forwarding", "token", opts.Token, "to", targetURL)

	for {
		if ctx.Err() != nil {
			return nil
		}

		if err := stream(ctx, opts, sseURL, targetURL); err != nil {
			slog.Warn("stream interrupted, reconnecting", "err", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(3 * time.Second):
		}
	}
}

func stream(ctx context.Context, opts Options, sseURL, targetURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sseURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/event-stream")
	if opts.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+opts.AuthToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("close stream body", "err", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		var event webhookEvent
		if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &event); err != nil {
			slog.Warn("parse event", "err", err)
			continue
		}

		go forward(ctx, opts, event, targetURL)
	}

	return scanner.Err()
}

func forward(ctx context.Context, opts Options, event webhookEvent, targetURL string) {
	var rawHeaders map[string][]string
	if err := json.Unmarshal([]byte(event.Headers), &rawHeaders); err != nil {
		slog.Warn("parse headers", "err", err)
	}

	var body io.Reader
	if len(event.Body) > 0 {
		body = bytes.NewReader(event.Body)
	}

	req, err := http.NewRequestWithContext(ctx, event.Method, targetURL, body)
	if err != nil {
		slog.Error("create forward request", "err", err)
		return
	}

	for k, vals := range rawHeaders {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	if opts.RewriteHost {
		if u, err := url.Parse(targetURL); err == nil {
			req.Host = u.Host
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("forward failed", "method", event.Method, "path", event.Path, "err", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("close forward body", "err", err)
		}
	}()

	slog.Info("forwarded", "method", event.Method, "path", event.Path, "status", resp.StatusCode)
}
