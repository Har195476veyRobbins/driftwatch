// Package filter provides service filtering for drift detection,
// allowing users to include or exclude specific services by name or label.
package filter

import "strings"

// Config holds include/exclude rules for service filtering.
type Config struct {
	// Include is an explicit allowlist of service names. Empty means all.
	Include []string
	// Exclude is a denylist of service names.
	Exclude []string
	// LabelSelector filters services by a compose label key=value pair.
	LabelSelector string
}

// Filter evaluates service names against a Config.
type Filter struct {
	cfg     Config
	labelKey string
	labelVal string
}

// New creates a Filter from the given Config.
func New(cfg Config) *Filter {
	f := &Filter{cfg: cfg}
	if cfg.LabelSelector != "" {
		parts := strings.SplitN(cfg.LabelSelector, "=", 2)
		f.labelKey = parts[0]
		if len(parts) == 2 {
			f.labelVal = parts[1]
		}
	}
	return f
}

// Allow returns true if the service should be included in drift detection.
func (f *Filter) Allow(name string) bool {
	for _, ex := range f.cfg.Exclude {
		if ex == name {
			return false
		}
	}
	if len(f.cfg.Include) == 0 {
		return true
	}
	for _, inc := range f.cfg.Include {
		if inc == name {
			return true
		}
	}
	return false
}

// AllowLabel returns true if the label map satisfies the LabelSelector.
// If no selector is configured, it always returns true.
func (f *Filter) AllowLabel(labels map[string]string) bool {
	if f.labelKey == "" {
		return true
	}
	v, ok := labels[f.labelKey]
	if !ok {
		return false
	}
	if f.labelVal == "" {
		return true
	}
	return v == f.labelVal
}
