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

// Generate creates a tatsu.yaml file by detecting the project type
// If force is true, it will overwrite an existing tatsu.yaml file
func Generate(force bool) error {
	// Check if tatsu.yaml already exists
	if _, err := os.Stat("tatsu.yaml"); err == nil && !force {
		return fmt.Errorf("tatsu.yaml already exists (use --force to overwrite)")
	}

	// Detect project type and generate config
	cfg := detectProjectType()

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile("tatsu.yaml", data, 0644); err != nil {
		return fmt.Errorf("failed to write tatsu.yaml: %w", err)
	}

	return nil
}

func detectProjectType() *Config {
	var cfg Config

	// Default agent command
	cfg.Agent.Command = `opencode run "%s"`

	// Detect project type by looking for common files
	if _, err := os.Stat("go.mod"); err == nil {
		// Go project
		cfg.Validate.Command = "go test ./..."
		return &cfg
	}

	if _, err := os.Stat("package.json"); err == nil {
		// Node.js project
		cfg.Validate.Command = "npm test"
		return &cfg
	}

	if _, err := os.Stat("requirements.txt"); err == nil {
		// Python project (requirements.txt)
		cfg.Validate.Command = "pytest"
		return &cfg
	}

	if _, err := os.Stat("pyproject.toml"); err == nil {
		// Python project (pyproject.toml)
		cfg.Validate.Command = "pytest"
		return &cfg
	}

	if _, err := os.Stat("Cargo.toml"); err == nil {
		// Rust project
		cfg.Validate.Command = "cargo test"
		return &cfg
	}

	// Default fallback
	cfg.Validate.Command = "echo 'No tests configured. Update tatsu.yaml with your test command.'"
	return &cfg
}
