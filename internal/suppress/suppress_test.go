package suppress

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_EmptyList(t *testing.T) {
	l := New()
	if l.IsSuppressed("svc") {
		t.Fatal("expected new list to have no suppressions")
	}
}

func TestSuppress_IsSuppressed(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)

	l.Suppress("web", 5*time.Minute)

	if !l.IsSuppressed("web") {
		t.Fatal("expected web to be suppressed")
	}
}

func TestSuppress_Expired(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)
	l.Suppress("web", 1*time.Minute)

	// advance time past expiry
	l.now = fixedNow(base.Add(2 * time.Minute))

	if l.IsSuppressed("web") {
		t.Fatal("expected suppression to have expired")
	}
}

func TestLift_RemovesSuppression(t *testing.T) {
	l := New()
	l.Suppress("db", 10*time.Minute)
	l.Lift("db")

	if l.IsSuppressed("db") {
		t.Fatal("expected suppression to be lifted")
	}
}

func TestSuppress_Extends(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)
	l.Suppress("web", 1*time.Minute)
	l.Suppress("web", 10*time.Minute) // extend

	l.now = fixedNow(base.Add(5 * time.Minute))
	if !l.IsSuppressed("web") {
		t.Fatal("expected extended suppression to still be active")
	}
}

func TestActive_ReturnsCurrent(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)
	l.Suppress("web", 5*time.Minute)
	l.Suppress("db", 5*time.Minute)

	actives := l.Active()
	if len(actives) != 2 {
		t.Fatalf("expected 2 active suppressions, got %d", len(actives))
	}
}

func TestActive_ExcludesExpired(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)
	l.Suppress("web", 1*time.Minute)
	l.Suppress("db", 10*time.Minute)

	l.now = fixedNow(base.Add(5 * time.Minute))
	actives := l.Active()
	if len(actives) != 1 {
		t.Fatalf("expected 1 active suppression, got %d", len(actives))
	}
	if actives[0].Service != "db" {
		t.Fatalf("expected db, got %s", actives[0].Service)
	}
}
