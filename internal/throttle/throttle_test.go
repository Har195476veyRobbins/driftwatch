package throttle

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MinInterval != 5*time.Minute {
		t.Fatalf("expected 5m, got %v", cfg.MinInterval)
	}
}

func TestNew_ZeroInterval_UsesDefault(t *testing.T) {
	th := New(Config{})
	if th.cfg.MinInterval != DefaultConfig().MinInterval {
		t.Fatalf("expected default interval, got %v", th.cfg.MinInterval)
	}
}

func TestAllow_FirstEvent_IsAllowed(t *testing.T) {
	th := New(Config{MinInterval: time.Second})
	if !th.Allow("svc-a") {
		t.Fatal("first event should be allowed")
	}
}

func TestAllow_WithinInterval_IsDenied(t *testing.T) {
	th := New(Config{MinInterval: time.Minute})
	now := time.Now()
	th.allowAt("svc-a", now)
	if th.allowAt("svc-a", now.Add(30*time.Second)) {
		t.Fatal("event within interval should be denied")
	}
}

func TestAllow_AfterInterval_IsAllowed(t *testing.T) {
	th := New(Config{MinInterval: time.Minute})
	now := time.Now()
	th.allowAt("svc-a", now)
	if !th.allowAt("svc-a", now.Add(61*time.Second)) {
		t.Fatal("event after interval should be allowed")
	}
}

func TestAllow_SeparateKeys_AreIndependent(t *testing.T) {
	th := New(Config{MinInterval: time.Minute})
	now := time.Now()
	th.allowAt("svc-a", now)
	// svc-b has never been seen, so it should be allowed
	if !th.allowAt("svc-b", now.Add(5*time.Second)) {
		t.Fatal("independent key should be allowed")
	}
}

func TestReset_AllowsNextEvent(t *testing.T) {
	th := New(Config{MinInterval: time.Minute})
	now := time.Now()
	th.allowAt("svc-a", now)
	th.Reset("svc-a")
	if !th.allowAt("svc-a", now.Add(time.Second)) {
		t.Fatal("event after reset should be allowed")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	th := New(Config{MinInterval: time.Minute})
	if th.Len() != 0 {
		t.Fatal("expected 0 keys initially")
	}
	th.Allow("svc-a")
	th.Allow("svc-b")
	if th.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", th.Len())
	}
	th.Reset("svc-a")
	if th.Len() != 1 {
		t.Fatalf("expected 1 key after reset, got %d", th.Len())
	}
}
