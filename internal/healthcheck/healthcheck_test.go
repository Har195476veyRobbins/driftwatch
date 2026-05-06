package healthcheck_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/driftwatch/internal/healthcheck"
)

func TestNew_InitialStateHealthy(t *testing.T) {
	h := healthcheck.New()
	if h == nil {
		t.Fatal("expected non-nil handler")
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUpdate_NoDrift_ReturnsHealthy(t *testing.T) {
	h := healthcheck.New()
	h.Update(0, nil)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var s healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !s.Healthy {
		t.Error("expected healthy=true")
	}
	if s.Message != "ok" {
		t.Errorf("expected message=ok, got %q", s.Message)
	}
}

func TestUpdate_WithError_ReturnsUnhealthy(t *testing.T) {
	h := healthcheck.New()
	h.Update(0, errors.New("docker connection refused"))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var s healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s.Healthy {
		t.Error("expected healthy=false")
	}
}

func TestUpdate_DriftCountReflected(t *testing.T) {
	h := healthcheck.New()
	h.Update(3, nil)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	var s healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s.DriftCount != 3 {
		t.Errorf("expected drift_count=3, got %d", s.DriftCount)
	}
}

func TestUpdate_LastCheckSet(t *testing.T) {
	h := healthcheck.New()
	h.Update(0, nil)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	var s healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s.LastCheck.IsZero() {
		t.Error("expected last_check to be set after Update")
	}
}
