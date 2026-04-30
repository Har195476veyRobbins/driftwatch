package compose

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceSpec represents a single service definition from a compose file.
type ServiceSpec struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
	Command     string            `yaml:"command"`
}

// ComposeSpec represents the top-level structure of a docker-compose file.
type ComposeSpec struct {
	Version  string                 `yaml:"version"`
	Services map[string]ServiceSpec `yaml:"services"`
}

// LoadSpec reads and parses a docker-compose YAML file from the given path.
func LoadSpec(path string) (*ComposeSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file %q: %w", path, err)
	}

	var spec ComposeSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parsing compose file %q: %w", path, err)
	}

	if spec.Services == nil {
		spec.Services = make(map[string]ServiceSpec)
	}

	return &spec, nil
}

// ServiceNames returns a sorted list of service names declared in the spec.
func (c *ComposeSpec) ServiceNames() []string {
	names := make([]string, 0, len(c.Services))
	for name := range c.Services {
		names = append(names, name)
	}
	return names
}
