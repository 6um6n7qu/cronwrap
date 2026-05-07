package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload is the JSON body sent to a webhook endpoint on job events.
type Payload struct {
	JobName   string        `json:"job_name"`
	Event     string        `json:"event"`
	ExitCode  int           `json:"exit_code,omitempty"`
	Duration  time.Duration `json:"duration_ms,omitempty"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// Client sends webhook notifications to a configured URL.
type Client struct {
	URL        string
	Secret     string
	httpClient *http.Client
}

// New returns a new Client. If httpClient is nil, a default client with a
// 10-second timeout is used.
func New(url, secret string, httpClient *http.Client) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("webhook: url must not be empty")
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{URL: url, Secret: secret, httpClient: httpClient}, nil
}

// Send serialises payload and POSTs it to the configured URL. An optional
// X-Webhook-Secret header is added when a secret is configured.
func (c *Client) Send(ctx context.Context, p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Secret != "" {
		req.Header.Set("X-Webhook-Secret", c.Secret)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, c.URL)
	}
	return nil
}
