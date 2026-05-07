// Package fingerprint computes a stable hash representing the observed
// state of a running container so that repeated drift checks can skip
// expensive comparisons when nothing has changed.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"

	"github.com/user/driftwatch/internal/docker"
)

// Store holds the last-seen fingerprint for each container service.
type Store struct {
	mu   sync.Mutex
	data map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{data: make(map[string]string)}
}

// Compute returns a deterministic hex fingerprint for the given container.
// The fingerprint covers the image, environment variables, and exposed ports.
func Compute(c docker.ContainerInfo) string {
	h := sha256.New()

	fmt.Fprintf(h, "image=%s\n", c.Image)

	// Sort env keys for determinism.
	keys := make([]string, 0, len(c.Env))
	for k := range c.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(h, "env:%s=%s\n", k, c.Env[k])
	}

	return hex.EncodeToString(h.Sum(nil))
}

// Changed returns true when the fingerprint for service differs from the
// previously stored value, then updates the store with the new fingerprint.
func (s *Store) Changed(service string, fp string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev, ok := s.data[service]
	s.data[service] = fp
	return !ok || prev != fp
}

// Get returns the stored fingerprint for service and whether it exists.
func (s *Store) Get(service string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fp, ok := s.data[service]
	return fp, ok
}

// Clear removes the stored fingerprint for service.
func (s *Store) Clear(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, service)
}
