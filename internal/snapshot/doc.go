// Package snapshot provides a thread-safe store for capturing and retrieving
// point-in-time container state. It enables driftwatch to compare the current
// running state against a previously recorded baseline, surfacing containers
// that have been added, removed, or changed between observation windows.
package snapshot
