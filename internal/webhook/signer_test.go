package webhook_test

import (
	"strings"
	"testing"

	"driftwatch/internal/webhook"
)

func TestHMACSigner_Sign_Prefix(t *testing.T) {
	s := webhook.NewHMACSigner("secret")
	sig := s.Sign([]byte("hello"))
	if !strings.HasPrefix(sig, "sha256=") {
		t.Errorf("expected sha256= prefix, got %s", sig)
	}
}

func TestHMACSigner_Sign_Deterministic(t *testing.T) {
	s := webhook.NewHMACSigner("secret")
	a := s.Sign([]byte("payload"))
	b := s.Sign([]byte("payload"))
	if a != b {
		t.Errorf("expected identical signatures, got %s vs %s", a, b)
	}
}

func TestHMACSigner_Sign_DifferentSecrets(t *testing.T) {
	s1 := webhook.NewHMACSigner("secret1")
	s2 := webhook.NewHMACSigner("secret2")
	if s1.Sign([]byte("data")) == s2.Sign([]byte("data")) {
		t.Error("expected different signatures for different secrets")
	}
}
