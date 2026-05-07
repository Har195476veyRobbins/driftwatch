// Package fingerprint provides lightweight SHA-256 fingerprinting of
// container state so that the drift-detection pipeline can skip
// unchanged services and reduce unnecessary comparisons.
//
// A Compute call hashes the container image reference plus all
// environment key=value pairs (sorted for determinism). The Store
// tracks the last-seen fingerprint per service and exposes a Changed
// helper that returns true the first time a fingerprint is seen or
// whenever it differs from the previous value.
package fingerprint
