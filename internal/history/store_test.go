package history_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/history"
)

func makeResults(drifted bool) []drift.Result {
	return []drift.Result{
		{Service: "web", Drifted: drifted},
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	s := history.New(0)
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestRecord_And_Len(t *testing.T) {
	s := history.New(10)
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Record(makeResults(false))
	s.Record(makeResults(true))
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestLast_Empty(t *testing.T) {
	s := history.New(5)
	_, ok := s.Last()
	if ok {
		t.Fatal("expected ok=false on empty store")
	}
}

func TestLast_ReturnsMostRecent(t *testing.T) {
	s := history.New(5)
	s.Record(makeResults(false))
	s.Record(makeResults(true))

	e, ok := s.Last()
	if !ok {
		t.Fatal("expected ok=true")
	}
	if !e.Drifted {
		t.Error("expected last entry to be drifted")
	}
	if e.Time.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if time.Since(e.Time) > 5*time.Second {
		t.Error("timestamp looks stale")
	}
}

func TestAll_Order(t *testing.T) {
	s := history.New(5)
	s.Record(makeResults(false))
	s.Record(makeResults(true))

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all[0].Drifted {
		t.Error("first entry should not be drifted")
	}
	if !all[1].Drifted {
		t.Error("second entry should be drifted")
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	cap := 3
	s := history.New(cap)
	for i := 0; i < cap+2; i++ {
		s.Record(makeResults(i%2 == 0))
	}
	if s.Len() != cap {
		t.Fatalf("expected %d entries, got %d", cap, s.Len())
	}
	// The ring should contain only the last `cap` records.
	all := s.All()
	if all[0].Drifted != (2%2 == 0) { // index 2 in original sequence
		t.Error("unexpected oldest entry after eviction")
	}
}
