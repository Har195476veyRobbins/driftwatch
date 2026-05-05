// Package metrics provides lightweight in-process counters for drift
// detection runs, making it easy to expose operational statistics.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of all counters.
type Snapshot struct {
	TotalRuns      int64
	DriftDetected  int64
	LastRunAt      time.Time
	LastDriftAt    time.Time
	ServicesScanned int64
}

// Collector accumulates drift-detection statistics.
type Collector struct {
	mu             sync.RWMutex
	totalRuns      int64
	driftDetected  int64
	lastRunAt      time.Time
	lastDriftAt    time.Time
	servicesScanned int64
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{}
}

// RecordRun records a completed detection run.
// servicesScanned is the number of compose services evaluated.
// hasDrift should be true when at least one drift result was reported.
func (c *Collector) RecordRun(servicesScanned int, hasDrift bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UTC()
	c.totalRuns++
	c.lastRunAt = now
	c.servicesScanned += int64(servicesScanned)

	if hasDrift {
		c.driftDetected++
		c.lastDriftAt = now
	}
}

// Snapshot returns a consistent copy of the current counters.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Snapshot{
		TotalRuns:       c.totalRuns,
		DriftDetected:   c.driftDetected,
		LastRunAt:       c.lastRunAt,
		LastDriftAt:     c.lastDriftAt,
		ServicesScanned: c.servicesScanned,
	}
}
