// Package rollup aggregates multiple drift results within a time window
// into a single summary, reducing notification noise during flapping.
package rollup

import (
	"sync"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Config controls rollup window behaviour.
type Config struct {
	// Window is the duration over which results are accumulated.
	Window time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 30 * time.Second,
	}
}

// Entry holds a single accumulated result set with its timestamp.
type Entry struct {
	RecordedAt time.Time
	Results    []drift.Result
}

// Rollup accumulates drift results within a sliding window and exposes
// a flattened view of all unique drifting services seen in that period.
type Rollup struct {
	mu      sync.Mutex
	cfg     Config
	entries []Entry
	now     func() time.Time
}

// New creates a Rollup with the given config.
// If cfg.Window is zero the default is applied.
func New(cfg Config) *Rollup {
	if cfg.Window <= 0 {
		cfg.Window = DefaultConfig().Window
	}
	return &Rollup{cfg: cfg, now: time.Now}
}

// Add appends a result set to the rollup window, evicting stale entries first.
func (r *Rollup) Add(results []drift.Result) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.evict()
	r.entries = append(r.entries, Entry{
		RecordedAt: r.now(),
		Results:    results,
	})
}

// Flatten returns deduplicated drift results across all entries in the window.
// The most recent result for each service name is returned.
func (r *Rollup) Flatten() []drift.Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.evict()
	seen := make(map[string]drift.Result)
	for _, e := range r.entries {
		for _, res := range e.Results {
			seen[res.Service] = res
		}
	}
	out := make([]drift.Result, 0, len(seen))
	for _, v := range seen {
		out = append(out, v)
	}
	return out
}

// Len returns the number of entries currently within the window.
func (r *Rollup) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evict()
	return len(r.entries)
}

// evict removes entries older than the configured window. Caller must hold mu.
func (r *Rollup) evict() {
	cutoff := r.now().Add(-r.cfg.Window)
	i := 0
	for i < len(r.entries) && r.entries[i].RecordedAt.Before(cutoff) {
		i++
	}
	r.entries = r.entries[i:]
}
