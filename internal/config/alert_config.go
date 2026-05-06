package config

// AlertConfig holds alerting threshold settings loaded from the config file.
type AlertConfig struct {
	// WarnThreshold is the minimum number of drifted services to trigger a
	// warn-level alert. 0 disables warn alerts.
	WarnThreshold int `yaml:"warn_threshold"`
	// CritThreshold is the minimum number of drifted services to trigger a
	// crit-level alert. 0 disables crit alerts.
	CritThreshold int `yaml:"crit_threshold"`
}

// defaultAlertConfig returns conservative defaults that flag any drift.
func defaultAlertConfig() AlertConfig {
	return AlertConfig{
		WarnThreshold: 1,
		CritThreshold: 3,
	}
}

// applyAlertDefaults fills zero-value fields with defaults.
func applyAlertDefaults(a *AlertConfig) {
	def := defaultAlertConfig()
	if a.WarnThreshold == 0 {
		a.WarnThreshold = def.WarnThreshold
	}
	if a.CritThreshold == 0 {
		a.CritThreshold = def.CritThreshold
	}
}
