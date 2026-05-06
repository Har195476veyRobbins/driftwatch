// Package webhook implements HTTP webhook delivery for driftwatch alerts.
//
// A Client is constructed with a Config specifying the target URL, request
// timeout, and an optional HMAC-SHA256 secret.  When a secret is set every
// outbound request carries an X-DriftWatch-Signature header so receivers can
// verify authenticity.
//
// Usage:
//
//	client := webhook.New(webhook.Config{URL: "https://example.com/hook"})
//	err := client.Send(ctx, webhook.Payload{...})
package webhook
