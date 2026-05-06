package config

import "testing"

func TestDefaultSnapshotConfig(t *testing.T) {
	c := defaultSnapshotConfig()

	if c.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if !c.BaselineOnStart {
		t.Error("expected BaselineOnStart=true by default")
	}
	if c.RetainCount != 5 {
		t.Errorf("expected RetainCount=5, got %d", c.RetainCount)
	}
}

func TestApplySnapshotDefaults_ZeroRetain(t *testing.T) {
	c := SnapshotConfig{RetainCount: 0}
	applySnapshotDefaults(&c)
	if c.RetainCount != 5 {
		t.Errorf("expected RetainCount=5 after applying defaults, got %d", c.RetainCount)
	}
}

func TestApplySnapshotDefaults_PreservesExplicit(t *testing.T) {
	c := SnapshotConfig{RetainCount: 10}
	applySnapshotDefaults(&c)
	if c.RetainCount != 10 {
		t.Errorf("expected RetainCount=10 preserved, got %d", c.RetainCount)
	}
}

func TestApplySnapshotDefaults_NegativeRetain(t *testing.T) {
	c := SnapshotConfig{RetainCount: -3}
	applySnapshotDefaults(&c)
	if c.RetainCount != 5 {
		t.Errorf("expected RetainCount reset to 5 for negative value, got %d", c.RetainCount)
	}
}
