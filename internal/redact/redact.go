// Package redact provides utilities for scrubbing sensitive values
// (e.g. passwords, tokens, secrets) from environment variable maps before
// they are written to logs, reports, or webhooks.
package redact

import "strings"

const placeholder = "[REDACTED]"

// defaultSensitiveKeys contains common key substrings that indicate a value
// should be redacted. Matching is case-insensitive.
var defaultSensitiveKeys = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"auth",
	"credential",
	"private_key",
	"access_key",
}

// Redactor scrubs sensitive environment variable values.
type Redactor struct {
	keys []string
}

// New returns a Redactor using the default sensitive key list merged with any
// additional keys supplied by the caller.
func New(extra ...string) *Redactor {
	keys := make([]string, len(defaultSensitiveKeys)+len(extra))
	copy(keys, defaultSensitiveKeys)
	for i, k := range extra {
		keys[len(defaultSensitiveKeys)+i] = strings.ToLower(k)
	}
	return &Redactor{keys: keys}
}

// Env returns a copy of the provided env map with sensitive values replaced by
// the redaction placeholder.
func (r *Redactor) Env(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.isSensitive(k) {
			out[k] = placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitive reports whether a key name is considered sensitive.
func (r *Redactor) IsSensitive(key string) bool {
	return r.isSensitive(key)
}

func (r *Redactor) isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range r.keys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}
