// Package alert implements threshold-based alerting for driftwatch.
//
// An Evaluator inspects a slice of drift.Result values and compares the
// count of drifted services against configurable warn and crit thresholds.
//
// Usage:
//
//	e := alert.New(alert.DefaultConfig(), os.Stderr)
//	a := e.Evaluate(results)
//	if a.Level != alert.LevelNone {
//		log.Println(a.Message)
//	}
package alert
