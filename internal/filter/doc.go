// Package filter provides include/exclude filtering for compose services
// during drift detection.
//
// A Filter can restrict which services are evaluated based on:
//   - An explicit include list (allowlist)
//   - An explicit exclude list (denylist)
//   - A label selector of the form "key" or "key=value"
//
// Exclude rules take precedence over include rules.
// When no rules are configured, all services are allowed.
package filter
