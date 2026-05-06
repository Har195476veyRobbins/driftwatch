package config

import "testing"

func TestDefaultAlertConfig(t *testing.T) {
	cfg := defaultAlertConfig()
	if cfg.WarnThreshold != 1 {
		t.Fatalf("expected WarnThreshold=1, got %d", cfg.WarnThreshold)
	}
	if cfg.CritThreshold != 3 {
		t.Fatalf("expected CritThreshold=3, got %d", cfg.CritThreshold)
	}
}

func TestApplyAlertDefaults_ZeroValues(t *testing.T) {
	a := AlertConfig{}
	applyAlertDefaults(&a)
	if a.WarnThreshold != 1 {
		t.Fatalf("expected WarnThreshold=1 after defaults, got %d", a.WarnThreshold)
	}
	if a.CritThreshold != 3 {
		t.Fatalf("expected CritThreshold=3 after defaults, got %d", a.CritThreshold)
	}
}

func TestApplyAlertDefaults_PreservesExplicitValues(t *testing.T) {
	a := AlertConfig{WarnThreshold: 2, CritThreshold: 5}
	applyAlertDefaults(&a)
	if a.WarnThreshold != 2 {
		t.Fatalf("expected WarnThreshold=2, got %d", a.WarnThreshold)
	}
	if a.CritThreshold != 5 {
		t.Fatalf("expected CritThreshold=5, got %d", a.CritThreshold)
	}
}

func TestApplyAlertDefaults_PartialZero(t *testing.T) {
	a := AlertConfig{WarnThreshold: 4, CritThreshold: 0}
	applyAlertDefaults(&a)
	if a.WarnThreshold != 4 {
		t.Fatalf("expected WarnThreshold=4, got %d", a.WarnThreshold)
	}
	if a.CritThreshold != 3 {
		t.Fatalf("expected CritThreshold=3 (default), got %d", a.CritThreshold)
	}
}
