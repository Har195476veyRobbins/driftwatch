// Package redact scrubs sensitive environment variable values before they
// appear in drift reports, audit logs, or outbound webhook payloads.
//
// A Redactor is constructed with a built-in list of common sensitive key
// substrings (password, secret, token, …) and can be extended with
// project-specific patterns at construction time.
//
// Usage:
//
//	red := redact.New("my_custom_secret")
//	safeEnv := red.Env(container.Env)
package redact
