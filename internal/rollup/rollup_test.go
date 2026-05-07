package rollup

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeResults(services ...string) []drift.Result {
	out := make([]drift.Result, 0, len(services))
	for _, s := range services {
		out = append(out, drift.Result{Service: s, Drifted: true})
	}
	return out
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Window != 30*time.Second {
		t.Fatalf("expected 30s window, got %v", cfg.Window)
	}
}

func TestNew_ZeroWindow_UsesDefault(t *testing.T) {
	r := New(Config{})
	if r.cfg.Window != DefaultConfig().Window {
		t.Fatalf("expected default window, got %v", r.cfg.Window)
	}
}

func TestLen_Empty(t *testing.T) {
	r := New(DefaultConfig())
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
}

func TestAdd_And_Len(t *testing.T) {
	r := New(DefaultConfig())
	r.Add(makeResults("web"))
	r.Add(makeResults("db"))
	if r.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", r.Len())
	}
}

func TestFlatten_DeduplicatesServices(t *testing.T) {
	r := New(DefaultConfig())
	r.Add(makeResults("web", "db"))
	r.Add(makeResults("web")) // duplicate service
	flat := r.Flatten()
	if len(flat) != 2 {
		t.Fatalf("expected 2 unique services, got %d", len(flat))
	}
}

func TestFlatten_Empty_ReturnsEmpty(t *testing.T) {
	r := New(DefaultConfig())
	if got := r.Flatten(); len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestEvict_RemovesStaleEntries(t *testing.T) {
	now := time.Now()
	r := New(Config{Window: 10 * time.Second})

	// inject a fake clock
	r.now = func() time.Time { return now }
	r.Add(makeResults("old-service"))

	// advance clock beyond the window
	r.now = func() time.Time { return now.Add(15 * time.Second) }
	r.Add(makeResults("new-service"))

	if r.Len() != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", r.Len())
	}
	flat := r.Flatten()
	if len(flat) != 1 || flat[0].Service != "new-service" {
		t.Fatalf("expected only new-service, got %v", flat)
	}
}

func TestFlatten_KeepsLatestResult(t *testing.T) {
	r := New(DefaultConfig())
	r.Add([]drift.Result{{Service: "web", Drifted: false}})
	r.Add([]drift.Result{{Service: "web", Drifted: true}})
	flat := r.Flatten()
	if len(flat) != 1 {
		t.Fatalf("expected 1 result, got %d", len(flat))
	}
	if !flat[0].Drifted {
		t.Fatal("expected latest (drifted=true) result to win")
	}
}
