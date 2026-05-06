// Package suppress provides a suppression list for drift alerts,
// allowing specific services to be silenced for a configured duration.
package suppress

import (
	"sync"
	"time"
)

// Entry represents a single suppression rule.
type Entry struct {
	Service   string
	ExpiresAt time.Time
}

// List holds active suppressions keyed by service name.
type List struct {
	mu      sync.Mutex
	entries map[string]time.Time
	now     func() time.Time
}

// New creates an empty suppression list.
func New() *List {
	return &List{
		entries: make(map[string]time.Time),
		now:     time.Now,
	}
}

// Suppress silences alerts for the given service until the duration elapses.
// Calling Suppress again before expiry extends the window.
func (l *List) Suppress(service string, d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries[service] = l.now().Add(d)
}

// IsSuppressed reports whether the service is currently suppressed.
func (l *List) IsSuppressed(service string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	exp, ok := l.entries[service]
	if !ok {
		return false
	}
	if l.now().After(exp) {
		delete(l.entries, service)
		return false
	}
	return true
}

// Lift removes a suppression for the given service immediately.
func (l *List) Lift(service string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, service)
}

// Active returns a snapshot of all currently active suppressions.
func (l *List) Active() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	out := make([]Entry, 0, len(l.entries))
	for svc, exp := range l.entries {
		if now.Before(exp) {
			out = append(out, Entry{Service: svc, ExpiresAt: exp})
		}
	}
	return out
}
