// Package healthcheck provides a simple HTTP health endpoint for driftwatch.
// It exposes liveness and readiness probes that report the daemon's current
// operational state and the timestamp of the last successful drift check.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status holds the current health state of the daemon.
type Status struct {
	Healthy   bool      `json:"healthy"`
	LastCheck time.Time `json:"last_check,omitempty"`
	DriftCount int      `json:"drift_count"`
	Message   string    `json:"message,omitempty"`
}

// Handler is an HTTP handler that serves health status.
type Handler struct {
	mu     sync.RWMutex
	status Status
}

// New creates a new Handler with a healthy initial state.
func New() *Handler {
	return &Handler{
		status: Status{
			Healthy: true,
			Message: "starting",
		},
	}
}

// Update records the result of the most recent drift check run.
func (h *Handler) Update(driftCount int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.status.LastCheck = time.Now().UTC()
	h.status.DriftCount = driftCount

	if err != nil {
		h.status.Healthy = false
		h.status.Message = err.Error()
	} else {
		h.status.Healthy = true
		h.status.Message = "ok"
	}
}

// ServeHTTP writes the current health status as JSON.
// Returns 200 when healthy, 503 when unhealthy.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	s := h.status
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	code := http.StatusOK
	if !s.Healthy {
		code = http.StatusServiceUnavailable
	}
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(s)
}
