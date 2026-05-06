// Package alert provides threshold-based alerting for drift detection runs.
// An alert is triggered when the number of drifted services exceeds a
// configurable threshold within a sliding window of recent runs.
package alert

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelCrit  Level = "crit"
	LevelNone  Level = "none"
)

// Config holds threshold configuration for alerting.
type Config struct {
	// WarnThreshold triggers a warn-level alert when drifted services >= value.
	WarnThreshold int
	// CritThreshold triggers a crit-level alert when drifted services >= value.
	CritThreshold int
}

// DefaultConfig returns sensible alert thresholds.
func DefaultConfig() Config {
	return Config{
		WarnThreshold: 1,
		CritThreshold: 3,
	}
}

// Alert describes a triggered alert.
type Alert struct {
	Level     Level
	Message   string
	DriftCount int
	TriggeredAt time.Time
}

// Evaluator checks drift results against configured thresholds.
type Evaluator struct {
	cfg Config
	out io.Writer
}

// New creates an Evaluator with the given config and output writer.
func New(cfg Config, out io.Writer) *Evaluator {
	return &Evaluator{cfg: cfg, out: out}
}

// Evaluate inspects results and returns an Alert (Level may be LevelNone).
func (e *Evaluator) Evaluate(results []drift.Result) Alert {
	count := 0
	for _, r := range results {
		if r.Drifted {
			count++
		}
	}

	lvl := LevelNone
	switch {
	case e.cfg.CritThreshold > 0 && count >= e.cfg.CritThreshold:
		lvl = LevelCrit
	case e.cfg.WarnThreshold > 0 && count >= e.cfg.WarnThreshold:
		lvl = LevelWarn
	}

	a := Alert{
		Level:       lvl,
		DriftCount:  count,
		TriggeredAt: time.Now().UTC(),
	}
	if lvl != LevelNone {
		a.Message = fmt.Sprintf("[%s] drift detected in %d service(s)", lvl, count)
		fmt.Fprintln(e.out, a.Message)
	}
	return a
}
