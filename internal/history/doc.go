// Package history provides a lightweight in-memory ring-buffer store for
// drift detection results produced by driftwatch.
//
// The store is safe for concurrent use and is intended to be shared between
// the scheduler loop (writer) and any reporting or notification path (reader).
//
// Usage:
//
//	store := history.New(50)          // keep last 50 runs
//	store.Record(results)             // called after each Detect run
//	last, ok := store.Last()          // inspect the most recent result
//	all  := store.All()               // iterate over all retained entries
package history
