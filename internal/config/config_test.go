package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestLoad_Defaults(t *testing.T) {
	p := writeTempConfig(t, "{}\n")
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ComposePath != config.DefaultComposePath {
		t.Errorf("ComposePath = %q, want %q", cfg.ComposePath, config.DefaultComposePath)
	}
	if cfg.Interval != config.DefaultInterval {
		t.Errorf("Interval = %v, want %v", cfg.Interval, config.DefaultInterval)
	}
	if cfg.Notify.Level != config.DefaultNotifyLevel {
		t.Errorf("Notify.Level = %q, want %q", cfg.Notify.Level, config.DefaultNotifyLevel)
	}
	if cfg.Notify.Output != config.DefaultOutput {
		t.Errorf("Notify.Output = %q, want %q", cfg.Notify.Output, config.DefaultOutput)
	}
}

func TestLoad_ExplicitValues(t *testing.T) {
	content := `
compose_path: my-compose.yml
interval: 1m
notify:
  level: verbose
  output: stderr
docker:
  host: tcp://localhost:2375
`
	p := writeTempConfig(t, content)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ComposePath != "my-compose.yml" {
		t.Errorf("ComposePath = %q", cfg.ComposePath)
	}
	if cfg.Interval != time.Minute {
		t.Errorf("Interval = %v, want 1m", cfg.Interval)
	}
	if cfg.Notify.Level != "verbose" {
		t.Errorf("Level = %q", cfg.Notify.Level)
	}
	if cfg.Docker.Host != "tcp://localhost:2375" {
		t.Errorf("Docker.Host = %q", cfg.Docker.Host)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/driftwatch.yml")
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

func TestLoad_InvalidNotifyLevel(t *testing.T) {
	p := writeTempConfig(t, "notify:\n  level: loud\n")
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for invalid notify level")
	}
}

func TestLoad_InvalidOutput(t *testing.T) {
	p := writeTempConfig(t, "notify:\n  output: file\n")
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for invalid output")
	}
}
