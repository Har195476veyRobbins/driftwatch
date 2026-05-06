// Package webhook provides HTTP webhook delivery for drift alerts.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Config holds webhook delivery configuration.
type Config struct {
	URL     string
	Timeout time.Duration
	Secret  string // optional HMAC secret for signing
}

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Timestamp  time.Time         `json:"timestamp"`
	DriftCount int               `json:"drift_count"`
	Services   []string          `json:"services"`
	Level      string            `json:"level"`
	Meta       map[string]string `json:"meta,omitempty"`
}

// Client delivers webhook payloads over HTTP.
type Client struct {
	cfg    Config
	http   *http.Client
	signer Signer
}

// New returns a new webhook Client.
func New(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	var s Signer
	if cfg.Secret != "" {
		s = NewHMACSigner(cfg.Secret)
	}
	return &Client{
		cfg:    cfg,
		http:   &http.Client{Timeout: timeout},
		signer: s,
	}
}

// Send encodes the payload and POSTs it to the configured URL.
func (c *Client) Send(ctx context.Context, p Payload) error {
	if c.cfg.URL == "" {
		return nil
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.signer != nil {
		sig := c.signer.Sign(body)
		req.Header.Set("X-DriftWatch-Signature", sig)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
