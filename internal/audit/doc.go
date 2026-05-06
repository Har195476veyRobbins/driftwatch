// Package audit implements a bounded, thread-safe audit log for driftwatch.
//
// Each drift detection cycle records an Event describing the outcome:
// whether drift was found, which services were affected, the alert level
// that was raised, and any error message. Events are stored in memory in
// insertion order and the log is capped at a configurable capacity so that
// long-running daemons do not grow without bound.
package audit
