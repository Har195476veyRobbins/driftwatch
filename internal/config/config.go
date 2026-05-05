// Package config loads and validates driftwatch daemon configuration.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level daemon configuration.
type Config struct {
	ComposePath string        `yaml:"compose_path"`
	Interval    time.Duration `yaml:"interval"`
	Notify      NotifyConfig  `yaml:"notify"`
	Docker      DockerConfig  `yaml:"docker"`
}

// NotifyConfig controls how drift results are reported.
type NotifyConfig struct {
	Level  string `yaml:"level"`  // silent | summary | verbose
	Output string `yaml:"output"` // stdout | stderr
}

// DockerConfig holds Docker daemon connection settings.
type DockerConfig struct {
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
}

// Defaults applied when fields are zero-valued.
const (
	DefaultComposePath = "docker-compose.yml"
	DefaultInterval    = 30 * time.Second
	DefaultNotifyLevel = "summary"
	DefaultOutput      = "stdout"
)

// Load reads a YAML config file from path and applies defaults.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	applyDefaults(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.ComposePath == "" {
		cfg.ComposePath = DefaultComposePath
	}
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultInterval
	}
	if cfg.Notify.Level == "" {
		cfg.Notify.Level = DefaultNotifyLevel
	}
	if cfg.Notify.Output == "" {
		cfg.Notify.Output = DefaultOutput
	}
}

func validate(cfg *Config) error {
	valid := map[string]bool{"silent": true, "summary": true, "verbose": true}
	if !valid[cfg.Notify.Level] {
		return fmt.Errorf("config: invalid notify level %q (want silent|summary|verbose)", cfg.Notify.Level)
	}
	if cfg.Notify.Output != "stdout" && cfg.Notify.Output != "stderr" {
		return fmt.Errorf("config: invalid output %q (want stdout|stderr)", cfg.Notify.Output)
	}
	return nil
}
