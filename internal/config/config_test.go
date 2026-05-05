package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_Defaults(t *testing.T) {
	p := writeTempConfig(t, "{}\n")
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ComposePath != "docker-compose.yml" {
		t.Errorf("ComposePath default: got %q", cfg.ComposePath)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("Interval default: got %v", cfg.Interval)
	}
	if cfg.Notify.Level != "summary" {
		t.Errorf("Notify.Level default: got %q", cfg.Notify.Level)
	}
}

func TestLoad_ExplicitValues(t *testing.T) {
	p := writeTempConfig(t, `
compose_path: /srv/compose.yml
interval: 60s
notify:
  level: verbose
  output: stderr
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ComposePath != "/srv/compose.yml" {
		t.Errorf("got %q", cfg.ComposePath)
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("got %v", cfg.Interval)
	}
	if cfg.Notify.Level != "verbose" {
		t.Errorf("got %q", cfg.Notify.Level)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeTempConfig(t, ": bad: yaml: [")
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_InvalidLevel(t *testing.T) {
	p := writeTempConfig(t, "notify:\n  level: loud\n")
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for invalid level")
	}
}

func TestLoad_FilterConfig(t *testing.T) {
	p := writeTempConfig(t, `
filter:
  include:
    - web
    - api
  exclude:
    - debug
  label_selector: env=prod
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Filter.Include) != 2 {
		t.Errorf("expected 2 include entries, got %d", len(cfg.Filter.Include))
	}
	if cfg.Filter.LabelSelector != "env=prod" {
		t.Errorf("got label_selector %q", cfg.Filter.LabelSelector)
	}
}
