package schedule

import (
	"context"
	"time"
)

// Config holds scheduler configuration.
type Config struct {
	// Interval between drift checks.
	Interval time.Duration
	// Jitter adds up to this duration of random delay to avoid thundering herd.
	Jitter time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval: 30 * time.Second,
		Jitter:   0,
	}
}

// Scheduler triggers a callback on a fixed interval until the context is done.
type Scheduler struct {
	cfg    Config
	ticker *time.Ticker
}

// New creates a new Scheduler with the given config.
func New(cfg Config) *Scheduler {
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultConfig().Interval
	}
	return &Scheduler{cfg: cfg}
}

// Run starts the scheduling loop, calling fn on each tick.
// It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context, fn func(ctx context.Context)) error {
	s.ticker = time.NewTicker(s.cfg.Interval)
	defer s.ticker.Stop()

	// Fire immediately on start.
	fn(ctx)

	for {
		select {
		case <-s.ticker.C:
			fn(ctx)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Reset changes the interval on a running scheduler.
func (s *Scheduler) Reset(interval time.Duration) {
	if s.ticker != nil && interval > 0 {
		s.ticker.Reset(interval)
		s.cfg.Interval = interval
	}
}
