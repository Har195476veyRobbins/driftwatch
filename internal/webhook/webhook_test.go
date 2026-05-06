package webhook_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"driftwatch/internal/webhook"
)

func TestSend_Success(t *testing.T) {
	var received webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := webhook.New(webhook.Config{URL: srv.URL})
	p := webhook.Payload{
		Timestamp:  time.Now(),
		DriftCount: 2,
		Services:   []string{"web", "db"},
		Level:      "warn",
	}
	if err := client.Send(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.DriftCount != 2 {
		t.Errorf("expected drift_count 2, got %d", received.DriftCount)
	}
}

func TestSend_EmptyURL_NoOp(t *testing.T) {
	client := webhook.New(webhook.Config{})
	if err := client.Send(context.Background(), webhook.Payload{}); err != nil {
		t.Fatalf("expected nil error for empty URL, got %v", err)
	}
}

func TestSend_NonSuccessStatus_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := webhook.New(webhook.Config{URL: srv.URL})
	if err := client.Send(context.Background(), webhook.Payload{}); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSend_WithSecret_SetsSignatureHeader(t *testing.T) {
	var sig string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sig = r.Header.Get("X-DriftWatch-Signature")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := webhook.New(webhook.Config{URL: srv.URL, Secret: "mysecret"})
	if err := client.Send(context.Background(), webhook.Payload{Level: "crit"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sig == "" {
		t.Error("expected X-DriftWatch-Signature header to be set")
	}
	if len(sig) < 8 || sig[:7] != "sha256=" {
		t.Errorf("signature format unexpected: %s", sig)
	}
}
