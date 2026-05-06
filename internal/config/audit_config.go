package config

// AuditConfig controls the in-memory audit log retained by driftwatch.
type AuditConfig struct {
	// Enabled controls whether the audit log is active.
	Enabled bool `yaml:"enabled"`
	// Capacity is the maximum number of events kept in memory.
	// When the log is full the oldest event is evicted.
	Capacity int `yaml:"capacity"`
}

func defaultAuditConfig() AuditConfig {
	return AuditConfig{
		Enabled:  true,
		Capacity: 256,
	}
}

func applyAuditDefaults(c *AuditConfig) {
	d := defaultAuditConfig()
	if c.Capacity <= 0 {
		c.Capacity = d.Capacity
	}
}
