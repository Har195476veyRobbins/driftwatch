package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/alert"
	"github.com/yourorg/driftwatch/internal/drift"
)

func makeResults(drifted, clean int) []drift.Result {
	var out []drift.Result
	for i := 0; i < drifted; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < clean; i++ {
		out = append(out, drift.Result{Service: "ok", Drifted: false})
	}
	return out
}

func TestEvaluate_NoDrift_LevelNone(t *testing.T) {
	var buf bytes.Buffer
	e := alert.New(alert.DefaultConfig(), &buf)
	a := e.Evaluate(makeResults(0, 3))
	if a.Level != alert.LevelNone {
		t.Fatalf("expected none, got %s", a.Level)
	}
	if buf.Len() != 0 {
		t.Fatal("expected no output for no drift")
	}
}

func TestEvaluate_OneDrift_LevelWarn(t *testing.T) {
	var buf bytes.Buffer
	e := alert.New(alert.DefaultConfig(), &buf)
	a := e.Evaluate(makeResults(1, 2))
	if a.Level != alert.LevelWarn {
		t.Fatalf("expected warn, got %s", a.Level)
	}
	if a.DriftCount != 1 {
		t.Fatalf("expected count 1, got %d", a.DriftCount)
	}
}

func TestEvaluate_CritThreshold(t *testing.T) {
	var buf bytes.Buffer
	e := alert.New(alert.DefaultConfig(), &buf)
	a := e.Evaluate(makeResults(3, 0))
	if a.Level != alert.LevelCrit {
		t.Fatalf("expected crit, got %s", a.Level)
	}
	if !strings.Contains(buf.String(), "crit") {
		t.Fatal("expected crit in output")
	}
}

func TestEvaluate_MessageContainsDriftCount(t *testing.T) {
	var buf bytes.Buffer
	e := alert.New(alert.DefaultConfig(), &buf)
	a := e.Evaluate(makeResults(2, 1))
	if !strings.Contains(a.Message, "2") {
		t.Fatalf("expected message to mention count, got: %s", a.Message)
	}
}

func TestDefaultConfig_Thresholds(t *testing.T) {
	cfg := alert.DefaultConfig()
	if cfg.WarnThreshold != 1 {
		t.Fatalf("expected warn=1, got %d", cfg.WarnThreshold)
	}
	if cfg.CritThreshold != 3 {
		t.Fatalf("expected crit=3, got %d", cfg.CritThreshold)
	}
}

func TestEvaluate_ZeroThreshold_Disabled(t *testing.T) {
	var buf bytes.Buffer
	cfg := alert.Config{WarnThreshold: 0, CritThreshold: 0}
	e := alert.New(cfg, &buf)
	a := e.Evaluate(makeResults(5, 0))
	if a.Level != alert.LevelNone {
		t.Fatalf("expected none when thresholds are zero, got %s", a.Level)
	}
}
