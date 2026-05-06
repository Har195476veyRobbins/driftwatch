// Package audit provides an append-only event log for drift detection runs,
// recording when drift was detected, which services were affected, and any
// alert level that was triggered.
package audit

import (
	"sync"
	"time"

	"github.com/user/driftwatch/internal/alert"
)

// EventKind classifies an audit event.
type EventKind string

const (
	EventDriftDetected EventKind = "drift_detected"
	EventNoDrift       EventKind = "no_drift"
	EventError         EventKind = "error"
)

// Event is a single audit log entry.
type Event struct {
	Timestamp   time.Time
	Kind        EventKind
	Level       alert.Level
	Services    []string
	DriftCount  int
	Message     string
}

// Log is a thread-safe, bounded audit event log.
type Log struct {
	mu       sync.RWMutex
	events   []Event
	capacity int
}

// New returns a Log that retains at most capacity events.
// If capacity is <= 0 the default of 256 is used.
func New(capacity int) *Log {
	if capacity <= 0 {
		capacity = 256
	}
	return &Log{capacity: capacity}
}

// Record appends an event to the log, evicting the oldest entry when full.
func (l *Log) Record(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.events) >= l.capacity {
		l.events = l.events[1:]
	}
	l.events = append(l.events, e)
}

// All returns a snapshot of all stored events, oldest first.
func (l *Log) All() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Len returns the number of stored events.
func (l *Log) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.events)
}

// Clear removes all events from the log.
func (l *Log) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = nil
}
