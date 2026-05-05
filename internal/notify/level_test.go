package notify_test

import (
	"testing"

	"github.com/user/driftwatch/internal/notify"
)

func TestParseLevel_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected notify.Level
	}{
		{"silent", notify.LevelSilent},
		{"summary", notify.LevelSummary},
		{"verbose", notify.LevelVerbose},
		{"SILENT", notify.LevelSilent},
		{"  Verbose  ", notify.LevelVerbose},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := notify.ParseLevel(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("ParseLevel(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	_, err := notify.ParseLevel("debug")
	if err == nil {
		t.Error("expected error for unknown level, got nil")
	}
}

func TestLevel_String(t *testing.T) {
	if notify.LevelSummary.String() != "summary" {
		t.Errorf("expected 'summary', got %q", notify.LevelSummary.String())
	}
}
