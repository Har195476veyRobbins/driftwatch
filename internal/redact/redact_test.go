package redact

import (
	"testing"
)

func TestNew_DefaultKeys(t *testing.T) {
	r := New()
	if len(r.keys) == 0 {
		t.Fatal("expected default keys to be populated")
	}
}

func TestNew_ExtraKeys_Appended(t *testing.T) {
	r := New("MY_CUSTOM")
	if !r.IsSensitive("MY_CUSTOM") {
		t.Error("expected MY_CUSTOM to be sensitive")
	}
}

func TestIsSensitive_MatchesSubstring(t *testing.T) {
	r := New()
	cases := []struct {
		key       string
		wantMatch bool
	}{
		{"DB_PASSWORD", true},
		{"api_key", true},
		{"AUTH_TOKEN", true},
		{"SECRET_VALUE", true},
		{"POSTGRES_HOST", false},
		{"LOG_LEVEL", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := r.IsSensitive(tc.key)
		if got != tc.wantMatch {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.wantMatch)
		}
	}
}

func TestEnv_RedactsSensitiveValues(t *testing.T) {
	r := New()
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"LOG_LEVEL":   "info",
		"API_TOKEN":   "tok-abc123",
	}

	got := r.Env(env)

	if got["DB_PASSWORD"] != placeholder {
		t.Errorf("DB_PASSWORD: got %q, want %q", got["DB_PASSWORD"], placeholder)
	}
	if got["API_TOKEN"] != placeholder {
		t.Errorf("API_TOKEN: got %q, want %q", got["API_TOKEN"], placeholder)
	}
	if got["LOG_LEVEL"] != "info" {
		t.Errorf("LOG_LEVEL: got %q, want %q", got["LOG_LEVEL"], "info")
	}
}

func TestEnv_DoesNotMutateOriginal(t *testing.T) {
	r := New()
	orig := map[string]string{"DB_PASSWORD": "secret123"}

	_ = r.Env(orig)

	if orig["DB_PASSWORD"] != "secret123" {
		t.Error("original map was mutated")
	}
}

func TestEnv_EmptyMap(t *testing.T) {
	r := New()
	got := r.Env(map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestEnv_CaseInsensitiveKey(t *testing.T) {
	r := New()
	env := map[string]string{"db_password": "val"}
	got := r.Env(env)
	if got["db_password"] != placeholder {
		t.Errorf("expected redaction for lowercase key, got %q", got["db_password"])
	}
}
