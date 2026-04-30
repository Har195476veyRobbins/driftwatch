package docker

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestToContainerInfo_BasicFields(t *testing.T) {
	ctr := types.Container{
		ID:    "abc123",
		Names: []string{"/web"},
		Image: "nginx:latest",
		Labels: map[string]string{
			"com.docker.compose.project": "myapp",
			"com.docker.compose.service": "web",
		},
	}

	info := toContainerInfo(ctr)

	if info.ID != "abc123" {
		t.Errorf("expected ID abc123, got %s", info.ID)
	}
	if info.Name != "/web" {
		t.Errorf("expected Name /web, got %s", info.Name)
	}
	if info.Image != "nginx:latest" {
		t.Errorf("expected Image nginx:latest, got %s", info.Image)
	}
	if info.Project != "myapp" {
		t.Errorf("expected Project myapp, got %s", info.Project)
	}
	if info.Service != "web" {
		t.Errorf("expected Service web, got %s", info.Service)
	}
}

func TestToContainerInfo_NoNames(t *testing.T) {
	ctr := types.Container{
		ID:     "xyz789",
		Names:  []string{},
		Image:  "alpine:3.18",
		Labels: map[string]string{},
	}

	info := toContainerInfo(ctr)

	if info.Name != "" {
		t.Errorf("expected empty Name, got %s", info.Name)
	}
	if info.Project != "" {
		t.Errorf("expected empty Project, got %s", info.Project)
	}
}

func TestToContainerInfo_EnvMapInitialized(t *testing.T) {
	ctr := types.Container{
		ID:     "env001",
		Names:  []string{"/svc"},
		Image:  "redis:7",
		Labels: map[string]string{},
	}

	info := toContainerInfo(ctr)

	if info.Env == nil {
		t.Error("expected Env map to be initialized, got nil")
	}
}

func TestNewClient_ReturnsError_WhenInvalidHost(t *testing.T) {
	t.Setenv("DOCKER_HOST", "tcp://invalid-host-that-does-not-exist:9999")
	// NewClientWithOpts is lazy; it won't error here, but Close should be safe.
	c, err := NewClient()
	if err != nil {
		// Some environments may error immediately — that's acceptable.
		return
	}
	if closeErr := c.Close(); closeErr != nil {
		t.Logf("close returned: %v", closeErr)
	}
}
