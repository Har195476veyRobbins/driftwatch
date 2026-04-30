package compose_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/driftwatch/internal/compose"
)

const sampleCompose = `
version: "3.9"
services:
  web:
    image: nginx:1.25
    ports:
      - "80:80"
    environment:
      ENV: production
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: app
      POSTGRES_USER: user
`

func writeTempCompose(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp compose file: %v", err)
	}
	return path
}

func TestLoadSpec_Valid(t *testing.T) {
	path := writeTempCompose(t, sampleCompose)

	spec, err := compose.LoadSpec(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spec.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(spec.Services))
	}

	web, ok := spec.Services["web"]
	if !ok {
		t.Fatal("expected service 'web' not found")
	}
	if web.Image != "nginx:1.25" {
		t.Errorf("expected image nginx:1.25, got %s", web.Image)
	}
}

func TestLoadSpec_MissingFile(t *testing.T) {
	_, err := compose.LoadSpec("/nonexistent/docker-compose.yml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadSpec_InvalidYAML(t *testing.T) {
	path := writeTempCompose(t, ":::invalid yaml:::")
	_, err := compose.LoadSpec(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestServiceNames(t *testing.T) {
	path := writeTempCompose(t, sampleCompose)
	spec, _ := compose.LoadSpec(path)

	names := spec.ServiceNames()
	if len(names) != 2 {
		t.Errorf("expected 2 service names, got %d", len(names))
	}
}
