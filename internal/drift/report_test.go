package drift_test

import (
	"strings"
	"testing"

	"driftwatch/internal/drift"
)

func TestSummarize_NoDrift(t *testing.T) {
	results := []drift.Result{
		{Service: "web", Diffs: nil},
		{Service: "db", Diffs: nil},
	}
	s := drift.Summarize(results)
	if s.DriftCount != 0 {
		t.Errorf("expected 0 drift, got %d", s.DriftCount)
	}
	if s.Total != 2 {
		t.Errorf("expected total 2, got %d", s.Total)
	}
}

func TestSummarize_WithDrift(t *testing.T) {
	results := []drift.Result{
		{Service: "web", Diffs: []string{"env PORT: want 8080 got 9090"}},
		{Service: "db", Diffs: nil},
	}
	s := drift.Summarize(results)
	if s.DriftCount != 1 {
		t.Errorf("expected 1 drifted service, got %d", s.DriftCount)
	}
}

func TestOneLiner_NoDrift(t *testing.T) {
	s := drift.Summary{Total: 3, DriftCount: 0}
	line := s.OneLiner()
	if !strings.Contains(line, "no drift") {
		t.Errorf("expected 'no drift' in output, got: %s", line)
	}
}

func TestOneLiner_WithDrift(t *testing.T) {
	s := drift.Summary{Total: 3, DriftCount: 2}
	line := s.OneLiner()
	if !strings.Contains(line, "2/3") {
		t.Errorf("expected '2/3' in output, got: %s", line)
	}
}

func TestWriteReport_ContainsServiceName(t *testing.T) {
	results := []drift.Result{
		{Service: "api", Diffs: []string{"env KEY: want foo got bar"}},
	}
	s := drift.Summarize(results)
	var buf strings.Builder
	drift.WriteReport(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "api") {
		t.Errorf("expected service name 'api' in report, got:\n%s", out)
	}
	if !strings.Contains(out, "env KEY") {
		t.Errorf("expected diff detail in report, got:\n%s", out)
	}
}
