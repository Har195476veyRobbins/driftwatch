package metrics_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()

	if s.TotalRuns != 0 {
		t.Errorf("expected TotalRuns=0, got %d", s.TotalRuns)
	}
	if s.DriftDetected != 0 {
		t.Errorf("expected DriftDetected=0, got %d", s.DriftDetected)
	}
	if !s.LastRunAt.IsZero() {
		t.Error("expected LastRunAt to be zero")
	}
}

func TestRecordRun_NoDrift(t *testing.T) {
	c := metrics.New()
	before := time.Now().UTC()
	c.RecordRun(3, false)
	after := time.Now().UTC()

	s := c.Snapshot()
	if s.TotalRuns != 1 {
		t.Errorf("expected TotalRuns=1, got %d", s.TotalRuns)
	}
	if s.DriftDetected != 0 {
		t.Errorf("expected DriftDetected=0, got %d", s.DriftDetected)
	}
	if s.ServicesScanned != 3 {
		t.Errorf("expected ServicesScanned=3, got %d", s.ServicesScanned)
	}
	if s.LastRunAt.Before(before) || s.LastRunAt.After(after) {
		t.Error("LastRunAt out of expected range")
	}
	if !s.LastDriftAt.IsZero() {
		t.Error("expected LastDriftAt to remain zero when no drift")
	}
}

func TestRecordRun_WithDrift(t *testing.T) {
	c := metrics.New()
	c.RecordRun(5, true)

	s := c.Snapshot()
	if s.DriftDetected != 1 {
		t.Errorf("expected DriftDetected=1, got %d", s.DriftDetected)
	}
	if s.LastDriftAt.IsZero() {
		t.Error("expected LastDriftAt to be set")
	}
}

func TestRecordRun_Accumulates(t *testing.T) {
	c := metrics.New()
	c.RecordRun(2, false)
	c.RecordRun(4, true)
	c.RecordRun(1, true)

	s := c.Snapshot()
	if s.TotalRuns != 3 {
		t.Errorf("expected TotalRuns=3, got %d", s.TotalRuns)
	}
	if s.DriftDetected != 2 {
		t.Errorf("expected DriftDetected=2, got %d", s.DriftDetected)
	}
	if s.ServicesScanned != 7 {
		t.Errorf("expected ServicesScanned=7, got %d", s.ServicesScanned)
	}
}
