package schedule

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
}

func TestNew_UsesDefaultWhenZeroInterval(t *testing.T) {
	s := New(Config{Interval: 0})
	if s.cfg.Interval != 30*time.Second {
		t.Errorf("expected default interval, got %v", s.cfg.Interval)
	}
}

func TestScheduler_FiresImmediately(t *testing.T) {
	var count int32

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	s := New(Config{Interval: 10 * time.Second})
	go func() {
		_ = s.Run(ctx, func(_ context.Context) {
			atomic.AddInt32(&count, 1)
		})
	}()

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected fn to be called immediately on start")
	}
}

func TestScheduler_CancelsOnContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	s := New(Config{Interval: 50 * time.Millisecond})
	done := make(chan error, 1)
	go func() {
		done <- s.Run(ctx, func(_ context.Context) {})
	}()

	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("scheduler did not stop after context cancellation")
	}
}

func TestScheduler_TicksMultipleTimes(t *testing.T) {
	var count int32

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	s := New(Config{Interval: 60 * time.Millisecond})
	_ = s.Run(ctx, func(_ context.Context) {
		atomic.AddInt32(&count, 1)
	})

	// initial call + ~3 ticks within 250ms
	if got := atomic.LoadInt32(&count); got < 3 {
		t.Errorf("expected at least 3 calls, got %d", got)
	}
}
