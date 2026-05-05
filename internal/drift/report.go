package drift

import (
	"fmt"
	"io"
	"strings"
)

// Summary holds aggregated drift information.
type Summary struct {
	Total      int
	DriftCount int
	Results    []Result
}

// OneLiner returns a short human-readable summary string.
func (s Summary) OneLiner() string {
	if s.DriftCount == 0 {
		return fmt.Sprintf("no drift detected (%d services checked)", s.Total)
	}
	return fmt.Sprintf("%d/%d services have drift", s.DriftCount, s.Total)
}

// Summarize aggregates a slice of Results into a Summary.
func Summarize(results []Result) Summary {
	s := Summary{Total: len(results), Results: results}
	for _, r := range results {
		if r.HasDrift() {
			s.DriftCount++
		}
	}
	return s
}

// WriteReport writes a detailed drift report to w.
func WriteReport(w io.Writer, summary Summary) {
	fmt.Fprintf(w, "Drift Report\n")
	fmt.Fprintf(w, "============\n")
	fmt.Fprintf(w, "Services checked : %d\n", summary.Total)
	fmt.Fprintf(w, "Services drifted : %d\n\n", summary.DriftCount)

	for _, r := range summary.Results {
		if !r.HasDrift() {
			continue
		}
		fmt.Fprintf(w, "Service: %s\n", r.Service)
		for _, d := range r.Diffs {
			fmt.Fprintf(w, "  - %s\n", strings.TrimSpace(d))
		}
	}
}
