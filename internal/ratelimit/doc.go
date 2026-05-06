// Package ratelimit provides a token-bucket style rate limiter for
// controlling how frequently drift alerts and webhook notifications
// are dispatched per service or global key.
//
// The limiter tracks per-key counters within a rolling window, allowing
// bursts up to a configured maximum before suppressing further events
// until the window resets.
package ratelimit
