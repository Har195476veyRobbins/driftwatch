package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/ratelimit"
)

func TestDefaultConfig(t *testing.T) {
	cfg := ratelimit.DefaultConfig()
	if cfg.Window <= 0 {
		t.Errorf("expected positive Window, got %v", cfg.Window)
	}
	if cfg.MaxEvents <= 0 {
		t.Errorf("expected positive MaxEvents, got %d", cfg.MaxEvents)
	}
}

func TestNew_AppliesDefaults(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{})
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestAllow_UnderLimit(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Window: time.Minute, MaxEvents: 3})
	for i := 0; i < 3; i++ {
		if !l.Allow("svc") {
			t.Fatalf("event %d should be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Window: time.Minute, MaxEvents: 2})
	l.Allow("svc")
	l.Allow("svc")
	if l.Allow("svc") {
		t.Error("third event should be denied")
	}
}

func TestAllow_SeparateKeys_Independent(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Window: time.Minute, MaxEvents: 1})
	if !l.Allow("a") {
		t.Error("first event for 'a' should be allowed")
	}
	if !l.Allow("b") {
		t.Error("first event for 'b' should be allowed")
	}
	if l.Allow("a") {
		t.Error("second event for 'a' should be denied")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	// Use a very short window so events expire quickly.
	l := ratelimit.New(ratelimit.Config{Window: 50 * time.Millisecond, MaxEvents: 1})
	if !l.Allow("svc") {
		t.Fatal("first event should be allowed")
	}
	if l.Allow("svc") {
		t.Fatal("second event should be denied within window")
	}
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("svc") {
		t.Error("event after window expiry should be allowed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Window: time.Minute, MaxEvents: 1})
	l.Allow("svc")
	l.Reset("svc")
	if !l.Allow("svc") {
		t.Error("event after Reset should be allowed")
	}
}
