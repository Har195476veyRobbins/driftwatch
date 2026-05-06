package audit_test

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/alert"
	"github.com/user/driftwatch/internal/audit"
)

func makeEvent(kind audit.EventKind, services ...string) audit.Event {
	return audit.Event{
		Kind:       kind,
		Level:      alert.LevelWarn,
		Services:   services,
		DriftCount: len(services),
		Message:    "test",
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	l := audit.New(0)
	if l == nil {
		t.Fatal("expected non-nil log")
	}
}

func TestRecord_And_Len(t *testing.T) {
	l := audit.New(10)
	l.Record(makeEvent(audit.EventNoDrift))
	l.Record(makeEvent(audit.EventDriftDetected, "svc-a"))
	if l.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", l.Len())
	}
}

func TestRecord_SetsTimestamp(t *testing.T) {
	l := audit.New(10)
	before := time.Now().UTC()
	l.Record(makeEvent(audit.EventNoDrift))
	events := l.All()
	if events[0].Timestamp.Before(before) {
		t.Error("timestamp should be set to approximately now")
	}
}

func TestRecord_PreservesExplicitTimestamp(t *testing.T) {
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e := makeEvent(audit.EventNoDrift)
	e.Timestamp = ts
	l := audit.New(10)
	l.Record(e)
	if !l.All()[0].Timestamp.Equal(ts) {
		t.Error("explicit timestamp should be preserved")
	}
}

func TestLog_Evicts_WhenFull(t *testing.T) {
	l := audit.New(3)
	for i := 0; i < 5; i++ {
		l.Record(makeEvent(audit.EventNoDrift))
	}
	if l.Len() != 3 {
		t.Fatalf("expected capacity 3, got %d", l.Len())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := audit.New(10)
	l.Record(makeEvent(audit.EventDriftDetected, "svc-a"))
	events := l.All()
	events[0].Message = "mutated"
	if l.All()[0].Message == "mutated" {
		t.Error("All() should return an isolated copy")
	}
}

func TestClear_RemovesAll(t *testing.T) {
	l := audit.New(10)
	l.Record(makeEvent(audit.EventError))
	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("expected 0 after Clear, got %d", l.Len())
	}
}
