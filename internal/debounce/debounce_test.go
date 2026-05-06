package debounce

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Cooldown <= 0 {
		t.Fatalf("expected positive cooldown, got %v", cfg.Cooldown)
	}
}

func TestNew_ZeroCooldown_UsesDefault(t *testing.T) {
	d := New(Config{Cooldown: 0})
	if d.cfg.Cooldown != DefaultConfig().Cooldown {
		t.Fatalf("expected default cooldown, got %v", d.cfg.Cooldown)
	}
}

func TestAllow_FirstEvent_IsAllowed(t *testing.T) {
	d := New(Config{Cooldown: time.Minute})
	if !d.Allow("svc-a") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_WithinCooldown_IsDenied(t *testing.T) {
	now := time.Now()
	d := New(Config{Cooldown: time.Minute})
	d.now = func() time.Time { return now }

	d.Allow("svc-a") // first — allowed
	if d.Allow("svc-a") {
		t.Fatal("expected second event within cooldown to be denied")
	}
}

func TestAllow_AfterCooldown_IsAllowed(t *testing.T) {
	base := time.Now()
	current := base
	d := New(Config{Cooldown: time.Minute})
	d.now = func() time.Time { return current }

	d.Allow("svc-a")

	current = base.Add(2 * time.Minute)
	if !d.Allow("svc-a") {
		t.Fatal("expected event after cooldown to be allowed")
	}
}

func TestAllow_DifferentKeys_Independent(t *testing.T) {
	now := time.Now()
	d := New(Config{Cooldown: time.Minute})
	d.now = func() time.Time { return now }

	d.Allow("svc-a")
	if !d.Allow("svc-b") {
		t.Fatal("expected independent key to be allowed")
	}
}

func TestReset_AllowsImmediateRetrigger(t *testing.T) {
	now := time.Now()
	d := New(Config{Cooldown: time.Minute})
	d.now = func() time.Time { return now }

	d.Allow("svc-a")
	d.Reset("svc-a")
	if !d.Allow("svc-a") {
		t.Fatal("expected allow after reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	d := New(Config{Cooldown: time.Minute})
	if d.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", d.Len())
	}
	d.Allow("svc-a")
	d.Allow("svc-b")
	if d.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", d.Len())
	}
	d.Reset("svc-a")
	if d.Len() != 1 {
		t.Fatalf("expected 1 key after reset, got %d", d.Len())
	}
}
