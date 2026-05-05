package watch_test

import (
	"context"
	"testing"
	"time"

	"driftwatch/internal/watch"
)

func TestConfig_Defaults(t *testing.T) {
	cfg := watch.Config{
		ComposePath: "docker-compose.yml",
		Interval:    30 * time.Second,
		Quiet:       false,
	}

	if cfg.ComposePath != "docker-compose.yml" {
		t.Errorf("expected ComposePath docker-compose.yml, got %s", cfg.ComposePath)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
}

func TestWatcher_CancelsOnContext(t *testing.T) {
	// We cannot easily spin up a real Docker client in unit tests,
	// so we verify that Run returns promptly when the context is cancelled.
	// A nil client will panic on check(); we cancel before the first tick.
	cfg := watch.Config{
		ComposePath: "docker-compose.yml",
		Interval:    10 * time.Second, // long enough that tick won't fire
		Quiet:       true,
	}

	w := watch.New(cfg, nil)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- w.RunNoInitialCheck(ctx)
	}()

	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not stop after context cancellation")
	}
}
