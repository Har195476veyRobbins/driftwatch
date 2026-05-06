// Package history records drift detection results over time,
// allowing driftwatch to surface trend information and avoid
// duplicate alerts for persistent drift.
package history

import (
	"sync"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Entry is a single recorded drift-check result.
type Entry struct {
	Time    time.Time
	Results []drift.Result
	Drifted bool
}

// Store holds an in-memory ring buffer of recent drift entries.
type Store struct {
	mu      sync.RWMutex
	entries []Entry
	cap     int
}

// New creates a Store that retains at most capacity entries.
// If capacity is <= 0 it defaults to 100.
func New(capacity int) *Store {
	if capacity <= 0 {
		capacity = 100
	}
	return &Store{cap: capacity}
}

// Record appends a new entry, evicting the oldest when full.
func (s *Store) Record(results []drift.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	drifted := false
	for _, r := range results {
		if r.Drifted {
			drifted = true
			break
		}
	}

	e := Entry{
		Time:    time.Now().UTC(),
		Results: results,
		Drifted: drifted,
	}

	if len(s.entries) >= s.cap {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, e)
}

// Last returns the most recent entry and true, or false if empty.
func (s *Store) Last() (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.entries) == 0 {
		return Entry{}, false
	}
	return s.entries[len(s.entries)-1], true
}

// All returns a shallow copy of all stored entries, oldest first.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Len returns the number of entries currently stored.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
