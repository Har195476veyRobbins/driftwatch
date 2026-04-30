package drift

import (
	"testing"

	"github.com/user/driftwatch/internal/docker"
)

func makeContainer(service string, env map[string]string) docker.ContainerInfo {
	return docker.ContainerInfo{
		ID:     "abc123",
		Names:  []string{"/" + service},
		Image:  "myimage:latest",
		Status: "running",
		Env:    env,
		Labels: map[string]string{
			"com.docker.compose.service": service,
		},
	}
}

func TestDetect_NoDrift(t *testing.T) {
	spec := map[string][]string{
		"web": {"PORT=8080", "DEBUG=false"},
	}
	containers := []docker.ContainerInfo{
		makeContainer("web", map[string]string{"PORT": "8080", "DEBUG": "false"}),
	}
	results := Detect(spec, containers)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, got reasons: %v", results[0].Reasons)
	}
}

func TestDetect_MissingService(t *testing.T) {
	spec := map[string][]string{
		"worker": {"QUEUE=default"},
	}
	results := Detect(spec, []docker.ContainerInfo{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Drifted {
		t.Error("expected drift for missing service")
	}
	if len(results[0].Reasons) == 0 {
		t.Error("expected at least one reason")
	}
}

func TestDetect_EnvValueMismatch(t *testing.T) {
	spec := map[string][]string{
		"api": {"LOG_LEVEL=info"},
	}
	containers := []docker.ContainerInfo{
		makeContainer("api", map[string]string{"LOG_LEVEL": "debug"}),
	}
	results := Detect(spec, containers)
	if !results[0].Drifted {
		t.Error("expected drift due to env mismatch")
	}
	if len(results[0].Reasons) != 1 {
		t.Errorf("expected 1 reason, got %d", len(results[0].Reasons))
	}
}

func TestDetect_MissingEnvKey(t *testing.T) {
	spec := map[string][]string{
		"db": {"POSTGRES_PASSWORD=secret"},
	}
	containers := []docker.ContainerInfo{
		makeContainer("db", map[string]string{}),
	}
	results := Detect(spec, containers)
	if !results[0].Drifted {
		t.Error("expected drift for missing env key")
	}
}

func TestDetect_MultipleServices_SortedOutput(t *testing.T) {
	spec := map[string][]string{
		"zebra": {},
		"alpha": {},
	}
	containers := []docker.ContainerInfo{
		makeContainer("zebra", map[string]string{}),
		makeContainer("alpha", map[string]string{}),
	}
	results := Detect(spec, containers)
	if results[0].Service != "alpha" || results[1].Service != "zebra" {
		t.Errorf("expected sorted results, got %v %v", results[0].Service, results[1].Service)
	}
}
