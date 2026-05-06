// Package cache provides a simple in-memory TTL cache used to deduplicate
// repeated drift events and avoid redundant notifications within a
// configurable time window.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached value along with its expiry time.
type Entry struct {
	Value     string
	ExpiresAt time.Time
}

// Cache is a thread-safe in-memory key/value store with per-entry TTL.
type Cache struct {
	mu      sync.Mutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given TTL applied to every Set call.
func New(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Set stores value under key, overwriting any previous entry.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get returns the value for key and whether it exists and has not expired.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return "", false
	}
	if time.Now().After(e.ExpiresAt) {
		delete(c.entries, key)
		return "", false
	}
	return e.Value, true
}

// Delete removes the entry for key if present.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Purge removes all expired entries and returns the number removed.
func (c *Cache) Purge() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	removed := 0
	for k, e := range c.entries {
		if now.After(e.ExpiresAt) {
			delete(c.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
