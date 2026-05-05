package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/notify"
)

func driftResults(hasDrift bool) []drift.Result {
	if !hasDrift {
		return []drift.Result{
			{Service: "web", Drifted: false},
		}
	}
	return []drift.Result{
		{
			Service: "web",
			Drifted: true,
			Diffs: []drift.Diff{
				{Key: "PORT", Expected: "8080", Actual: "9090"},
			},
		},
	}
}

func TestNotify_Silent_WritesNothing(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.LevelSilent, &buf)
	if err := n.Notify(driftResults(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for silent level, got: %q", buf.String())
	}
}

func TestNotify_Summary_ContainsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.LevelSummary, &buf)
	if err := n.Notify(driftResults(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[") {
		t.Errorf("expected timestamp bracket in output, got: %q", out)
	}
}

func TestNotify_Summary_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.LevelSummary, &buf)
	_ = n.Notify(driftResults(false))
	if !strings.Contains(buf.String(), "no drift") {
		t.Errorf("expected 'no drift' in summary output, got: %q", buf.String())
	}
}

func TestNotify_Verbose_ContainsServiceName(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.LevelVerbose, &buf)
	_ = n.Notify(driftResults(true))
	if !strings.Contains(buf.String(), "web") {
		t.Errorf("expected service name 'web' in verbose output, got: %q", buf.String())
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Just ensure no panic when writer is nil.
	n := notify.New(notify.LevelSilent, nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
