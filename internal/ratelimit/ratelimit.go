// Package ratelimit provides a simple token-bucket rate limiter for
// suppressing repeated alert notifications within a configurable window.
package ratelimit

import (
	"sync"
	"time"
)

// Config holds configuration for the rate limiter.
type Config struct {
	// Window is the duration during which at most MaxEvents are allowed.
	Window time.Duration
	// MaxEvents is the maximum number of events permitted within Window.
	MaxEvents int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window:    5 * time.Minute,
		MaxEvents: 3,
	}
}

// Limiter tracks event occurrences per named key and reports whether a new
// event should be allowed through.
type Limiter struct {
	cfg    Config
	mu     sync.Mutex
	bucket map[string][]time.Time
}

// New creates a Limiter using the provided Config. Zero-value fields are
// replaced with defaults.
func New(cfg Config) *Limiter {
	def := DefaultConfig()
	if cfg.Window <= 0 {
		cfg.Window = def.Window
	}
	if cfg.MaxEvents <= 0 {
		cfg.MaxEvents = def.MaxEvents
	}
	return &Limiter{
		cfg:    cfg,
		bucket: make(map[string][]time.Time),
	}
}

// Allow returns true if the event identified by key is within the allowed
// rate. It records the event timestamp and prunes stale entries.
func (l *Limiter) Allow(key string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := now.Add(-l.cfg.Window)
	times := l.bucket[key]

	// Prune events outside the window.
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.cfg.MaxEvents {
		l.bucket[key] = valid
		return false
	}

	l.bucket[key] = append(valid, now)
	return true
}

// Reset clears all recorded events for the given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.bucket, key)
}
