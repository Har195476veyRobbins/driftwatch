// Package config loads and validates driftwatch daemon configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the top-level daemon configuration.
type Config struct {
	ComposePath string        `yaml:"compose_path"`
	Interval    time.Duration `yaml:"interval"`
	DockerHost  string        `yaml:"docker_host"`
	Notify      NotifyConfig  `yaml:"notify"`
	Filter      FilterConfig  `yaml:"filter"`
}

// NotifyConfig controls how drift results are reported.
type NotifyConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

// FilterConfig mirrors filter.Config for YAML unmarshalling.
type FilterConfig struct {
	Include       []string `yaml:"include"`
	Exclude       []string `yaml:"exclude"`
	LabelSelector string   `yaml:"label_selector"`
}

const (
	defaultComposePath = "docker-compose.yml"
	defaultInterval    = 30 * time.Second
	defaultNotifyLevel = "summary"
	defaultOutput      = "stdout"
)

// Load reads a YAML config file from path and returns a validated Config.
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
		return nil, fmt.Errorf("config: invalid: %w", err)
	}
	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.ComposePath == "" {
		cfg.ComposePath = defaultComposePath
	}
	if cfg.Interval == 0 {
		cfg.Interval = defaultInterval
	}
	if cfg.Notify.Level == "" {
		cfg.Notify.Level = defaultNotifyLevel
	}
	if cfg.Notify.Output == "" {
		cfg.Notify.Output = defaultOutput
	}
}

func validate(cfg *Config) error {
	if cfg.Interval < time.Second {
		return errors.New("interval must be at least 1s")
	}
	validLevels := map[string]bool{"silent": true, "summary": true, "verbose": true}
	if !validLevels[cfg.Notify.Level] {
		return fmt.Errorf("notify.level %q is not valid (silent|summary|verbose)", cfg.Notify.Level)
	}
	return nil
}
