package notify

import (
	"fmt"
	"strings"
)

// ParseLevel converts a string to a Level, returning an error for unknown values.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case string(LevelSilent):
		return LevelSilent, nil
	case string(LevelSummary):
		return LevelSummary, nil
	case string(LevelVerbose):
		return LevelVerbose, nil
	default:
		return "", fmt.Errorf("unknown notify level %q: must be one of silent, summary, verbose", s)
	}
}

// String implements fmt.Stringer.
func (l Level) String() string {
	return string(l)
}
