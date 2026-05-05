package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// Level controls the verbosity of notifications.
type Level string

const (
	LevelSilent  Level = "silent"
	LevelSummary Level = "summary"
	LevelVerbose Level = "verbose"
)

// Notifier sends drift reports to one or more outputs.
type Notifier struct {
	level  Level
	writer io.Writer
}

// New creates a Notifier that writes to the given writer at the given level.
// If writer is nil, os.Stdout is used.
func New(level Level, writer io.Writer) *Notifier {
	if writer == nil {
		writer = os.Stdout
	}
	return &Notifier{level: level, writer: writer}
}

// Notify emits a notification for the given drift results.
// It is a no-op when Level is LevelSilent.
func (n *Notifier) Notify(results []drift.Result) error {
	if n.level == LevelSilent {
		return nil
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)

	if n.level == LevelSummary {
		summary := drift.Summarize(results)
		_, err := fmt.Fprintf(n.writer, "[%s] %s\n", timestamp, summary)
		return err
	}

	// LevelVerbose
	_, err := fmt.Fprintf(n.writer, "[%s] drift check\n", timestamp)
	if err != nil {
		return err
	}
	return drift.WriteReport(n.writer, results)
}
