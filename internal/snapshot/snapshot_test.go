package snapshot_test

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/docker"
	"github.com/user/driftwatch/internal/snapshot"
)

func makeContainers(services ...string) []docker.ContainerInfo {
	out := make([]docker.ContainerInfo, len(services))
	for i, s := range services {
		out[i] = docker.ContainerInfo{Service: s, ID: s + "-id"}
	}
	return out
}

func TestNew_EmptyStore(t *testing.T) {
	st := snapshot.New()
	_, ok := st.Get()
	if ok {
		t.Fatal("expected empty store, got snapshot")
	}
}

func TestSave_And_Get(t *testing.T) {
	st := snapshot.New()
	containers := makeContainers("web", "db")

	before := time.Now()
	st.Save(containers)
	after := time.Now()

	snap, ok := st.Get()
	if !ok {
		t.Fatal("expected snapshot, got none")
	}
	if len(snap.Containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(snap.Containers))
	}
	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Error("CapturedAt timestamp out of expected range")
	}
}

func TestSave_IsolatesCopy(t *testing.T) {
	st := snapshot.New()
	containers := makeContainers("web")
	st.Save(containers)

	// mutate original slice — snapshot must be unaffected
	containers[0].Service = "mutated"

	snap, _ := st.Get()
	if snap.Containers[0].Service != "web" {
		t.Error("snapshot was not isolated from original slice mutation")
	}
}

func TestClear_RemovesSnapshot(t *testing.T) {
	st := snapshot.New()
	st.Save(makeContainers("web"))
	st.Clear()

	_, ok := st.Get()
	if ok {
		t.Error("expected no snapshot after Clear")
	}
}

func TestDiff_ReturnsRemoved(t *testing.T) {
	prev := &snapshot.Snapshot{Containers: makeContainers("web", "db", "cache")}
	next := &snapshot.Snapshot{Containers: makeContainers("web", "db")}

	removed := snapshot.Diff(prev, next)
	if len(removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(removed))
	}
	if removed[0].Service != "cache" {
		t.Errorf("expected 'cache' removed, got %q", removed[0].Service)
	}
}

func TestDiff_NilInputs(t *testing.T) {
	if got := snapshot.Diff(nil, nil); got != nil {
		t.Error("expected nil diff for nil inputs")
	}
}

func TestDiff_NoDrift(t *testing.T) {
	prev := &snapshot.Snapshot{Containers: makeContainers("web", "db")}
	next := &snapshot.Snapshot{Containers: makeContainers("web", "db")}

	if removed := snapshot.Diff(prev, next); len(removed) != 0 {
		t.Errorf("expected no removed containers, got %d", len(removed))
	}
}
