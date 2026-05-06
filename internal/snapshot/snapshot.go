// Package snapshot captures and compares point-in-time container state
// so drift can be evaluated against a known-good baseline.
package snapshot

import (
	"sync"
	"time"

	"github.com/user/driftwatch/internal/docker"
)

// Snapshot holds the container state captured at a specific moment.
type Snapshot struct {
	CapturedAt time.Time
	Containers []docker.ContainerInfo
}

// Store holds the most recent snapshot and provides thread-safe access.
type Store struct {
	mu       sync.RWMutex
	current  *Snapshot
}

// New returns an empty snapshot Store.
func New() *Store {
	return &Store{}
}

// Save replaces the current snapshot with a new one built from containers.
func (s *Store) Save(containers []docker.ContainerInfo) {
	snap := &Snapshot{
		CapturedAt: time.Now(),
		Containers: make([]docker.ContainerInfo, len(containers)),
	}
	copy(snap.Containers, containers)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = snap
}

// Get returns the most recent snapshot and whether one exists.
func (s *Store) Get() (*Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.current == nil {
		return nil, false
	}
	return s.current, true
}

// Clear removes the stored snapshot.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = nil
}

// Diff returns containers present in prev but absent (by service name) in next.
func Diff(prev, next *Snapshot) []docker.ContainerInfo {
	if prev == nil || next == nil {
		return nil
	}

	nextIndex := make(map[string]struct{}, len(next.Containers))
	for _, c := range next.Containers {
		nextIndex[c.Service] = struct{}{}
	}

	var removed []docker.ContainerInfo
	for _, c := range prev.Containers {
		if _, ok := nextIndex[c.Service]; !ok {
			removed = append(removed, c)
		}
	}
	return removed
}
