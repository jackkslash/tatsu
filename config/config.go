package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent struct {
		Command string `yaml:"command"`
	} `yaml:"agent"`
	Validate struct {
		Command string `yaml:"command"`
	} `yaml:"validate"`
}

func Load() (*Config, error) {
	// Check file exists
	if _, err := os.Stat("tatsu.yaml"); os.IsNotExist(err) {
		return nil, fmt.Errorf("tatsu.yaml not found. Create one with:\n\nagent:\n  command: 'opencode run \"%%s\"'\nvalidate:\n  command: 'go test ./...'")
	}

	// Read file
	data, err := os.ReadFile("tatsu.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read tatsu.yaml: %w", err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate required fields
	if cfg.Agent.Command == "" {
		return nil, fmt.Errorf("agent.command is required in tatsu.yaml")
	}
	if cfg.Validate.Command == "" {
		return nil, fmt.Errorf("validate.command is required in tatsu.yaml")
	}

	return &cfg, nil
}
