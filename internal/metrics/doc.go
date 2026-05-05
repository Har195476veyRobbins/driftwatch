// Package metrics provides a lightweight, thread-safe collector for
// driftwatch operational statistics.
//
// Usage:
//
//	col := metrics.New()
//
//	// after each detection run:
//	col.RecordRun(len(services), hasDrift)
//
//	// read a consistent snapshot at any time:
//	snap := col.Snapshot()
//	fmt.Printf("runs=%d drift=%d\n", snap.TotalRuns, snap.DriftDetected)
//
// The Collector is safe for concurrent use from multiple goroutines.
package metrics
