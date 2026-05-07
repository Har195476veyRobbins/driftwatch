package fingerprint_test

import (
	"testing"

	"github.com/user/driftwatch/internal/docker"
	"github.com/user/driftwatch/internal/fingerprint"
)

func makeContainer(image string, env map[string]string) docker.ContainerInfo {
	return docker.ContainerInfo{
		ID:    "abc123",
		Image: image,
		Env:   env,
	}
}

func TestCompute_Deterministic(t *testing.T) {
	c := makeContainer("nginx:1.25", map[string]string{"PORT": "80", "ENV": "prod"})
	a := fingerprint.Compute(c)
	b := fingerprint.Compute(c)
	if a != b {
		t.Fatalf("expected identical fingerprints, got %q and %q", a, b)
	}
}

func TestCompute_DiffersOnImageChange(t *testing.T) {
	env := map[string]string{"PORT": "80"}
	a := fingerprint.Compute(makeContainer("nginx:1.25", env))
	b := fingerprint.Compute(makeContainer("nginx:1.26", env))
	if a == b {
		t.Fatal("expected different fingerprints for different images")
	}
}

func TestCompute_DiffersOnEnvChange(t *testing.T) {
	a := fingerprint.Compute(makeContainer("nginx:1.25", map[string]string{"PORT": "80"}))
	b := fingerprint.Compute(makeContainer("nginx:1.25", map[string]string{"PORT": "443"}))
	if a == b {
		t.Fatal("expected different fingerprints for different env values")
	}
}

func TestCompute_EnvOrderIndependent(t *testing.T) {
	a := fingerprint.Compute(makeContainer("app:1", map[string]string{"A": "1", "B": "2"}))
	b := fingerprint.Compute(makeContainer("app:1", map[string]string{"B": "2", "A": "1"}))
	if a != b {
		t.Fatal("expected same fingerprint regardless of env map iteration order")
	}
}

func TestChanged_FirstCallIsAlwaysChanged(t *testing.T) {
	s := fingerprint.New()
	fp := fingerprint.Compute(makeContainer("app:1", nil))
	if !s.Changed("web", fp) {
		t.Fatal("first call to Changed should return true")
	}
}

func TestChanged_SameFingerprintReturnsFalse(t *testing.T) {
	s := fingerprint.New()
	fp := fingerprint.Compute(makeContainer("app:1", nil))
	s.Changed("web", fp)
	if s.Changed("web", fp) {
		t.Fatal("second call with same fingerprint should return false")
	}
}

func TestChanged_DifferentFingerprintReturnsTrue(t *testing.T) {
	s := fingerprint.New()
	fp1 := fingerprint.Compute(makeContainer("app:1", nil))
	fp2 := fingerprint.Compute(makeContainer("app:2", nil))
	s.Changed("web", fp1)
	if !s.Changed("web", fp2) {
		t.Fatal("changed fingerprint should return true")
	}
}

func TestGet_ReturnsStoredFingerprint(t *testing.T) {
	s := fingerprint.New()
	fp := fingerprint.Compute(makeContainer("app:1", nil))
	s.Changed("web", fp)
	got, ok := s.Get("web")
	if !ok || got != fp {
		t.Fatalf("Get returned (%q, %v), want (%q, true)", got, ok, fp)
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	s := fingerprint.New()
	fp := fingerprint.Compute(makeContainer("app:1", nil))
	s.Changed("web", fp)
	s.Clear("web")
	_, ok := s.Get("web")
	if ok {
		t.Fatal("expected Get to return false after Clear")
	}
}
