// Package debounce provides a mechanism to suppress repeated drift
// notifications within a configurable cooldown window. When the same
// service triggers a drift event multiple times in quick succession,
// only the first event is forwarded until the cooldown expires.
package debounce

import (
	"sync"
	"time"
)

// Config holds debounce configuration.
type Config struct {
	// Cooldown is the minimum duration between notifications for the same key.
	Cooldown time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Cooldown: 5 * time.Minute,
	}
}

// Debouncer tracks the last notification time per key and decides
// whether a new event should be allowed through.
type Debouncer struct {
	mu       sync.Mutex
	cfg      Config
	lastSeen map[string]time.Time
	now      func() time.Time
}

// New creates a new Debouncer. If cfg.Cooldown is zero, DefaultConfig is used.
func New(cfg Config) *Debouncer {
	if cfg.Cooldown <= 0 {
		cfg = DefaultConfig()
	}
	return &Debouncer{
		cfg:      cfg,
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event for key should be forwarded.
// It returns false when the key was seen within the cooldown window.
func (d *Debouncer) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if last, ok := d.lastSeen[key]; ok {
		if now.Sub(last) < d.cfg.Cooldown {
			return false
		}
	}
	d.lastSeen[key] = now
	return true
}

// Reset clears the debounce state for key, allowing the next event
// through immediately regardless of the cooldown window.
func (d *Debouncer) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.lastSeen, key)
}

// Len returns the number of keys currently tracked.
func (d *Debouncer) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.lastSeen)
}
