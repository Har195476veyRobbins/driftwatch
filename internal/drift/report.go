package drift

import (
	"fmt"
	"io"
	"strings"
)

// Summary holds aggregate statistics from a drift detection run.
type Summary struct {
	Total   int
	Drifted int
	Clean   int
}

// Summarize computes aggregate stats from a slice of DriftResults.
func Summarize(results []DriftResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Drifted {
			s.Drifted++
		} else {
			s.Clean++
		}
	}
	return s
}

// WriteReport writes a human-readable drift report to w.
func WriteReport(w io.Writer, results []DriftResult) {
	for _, r := range results {
		if r.Drifted {
			fmt.Fprintf(w, "[DRIFT]  %s\n", r.Service)
			for _, reason := range r.Reasons {
				fmt.Fprintf(w, "         - %s\n", reason)
			}
		} else {
			fmt.Fprintf(w, "[OK]     %s\n", r.Service)
		}
	}

	s := Summarize(results)
	fmt.Fprintf(w, "\n%s\n", strings.Repeat("-", 40))
	fmt.Fprintf(w, "Total: %d  Clean: %d  Drifted: %d\n", s.Total, s.Clean, s.Drifted)
}
