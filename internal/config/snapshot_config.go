package config

// SnapshotConfig controls baseline snapshot behaviour.
type SnapshotConfig struct {
	// Enabled turns snapshot-based drift comparison on or off.
	Enabled bool `yaml:"enabled"`

	// BaselineOnStart captures a snapshot immediately when the daemon starts
	// so the first tick has a reference point.
	BaselineOnStart bool `yaml:"baseline_on_start"`

	// RetainCount is the maximum number of historical snapshots kept in memory.
	// Must be >= 1; defaults to 5.
	RetainCount int `yaml:"retain_count"`
}

func defaultSnapshotConfig() SnapshotConfig {
	return SnapshotConfig{
		Enabled:         false,
		BaselineOnStart: true,
		RetainCount:     5,
	}
}

func applySnapshotDefaults(c *SnapshotConfig) {
	defaults := defaultSnapshotConfig()
	if c.RetainCount <= 0 {
		c.RetainCount = defaults.RetainCount
	}
}
