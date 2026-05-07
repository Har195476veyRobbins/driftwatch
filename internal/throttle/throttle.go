// Package throttle provides a token-bucket style throttle that limits
// how frequently drift alerts can be emitted per service.
package throttle

import (
	"sync"
	"time"
)

// Config holds throttle configuration.
type Config struct {
	// MinInterval is the minimum time between allowed events for a given key.
	MinInterval time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MinInterval: 5 * time.Minute,
	}
}

// Throttle tracks the last allowed event time per key and suppresses
// subsequent events that arrive within MinInterval.
type Throttle struct {
	mu     sync.Mutex
	cfg    Config
	lastAt map[string]time.Time
}

// New creates a new Throttle. Zero-value Config fields are replaced with defaults.
func New(cfg Config) *Throttle {
	if cfg.MinInterval <= 0 {
		cfg.MinInterval = DefaultConfig().MinInterval
	}
	return &Throttle{
		cfg:    cfg,
		lastAt: make(map[string]time.Time),
	}
}

// Allow returns true if the event for key should be allowed through,
// i.e. no event for that key has been allowed within MinInterval.
// When true is returned the internal timestamp for key is updated.
func (t *Throttle) Allow(key string) bool {
	return t.allowAt(key, time.Now())
}

// allowAt is the testable core of Allow.
func (t *Throttle) allowAt(key string, now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	last, seen := t.lastAt[key]
	if seen && now.Sub(last) < t.cfg.MinInterval {
		return false
	}
	t.lastAt[key] = now
	return true
}

// Reset clears the throttle state for key, allowing the next event through
// regardless of when it arrives.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastAt, key)
}

// Len returns the number of keys currently tracked.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.lastAt)
}
