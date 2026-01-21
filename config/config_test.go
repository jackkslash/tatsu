package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	// Create a temporary config file
	content := `agent:
  command: 'opencode run "%s"'
validate:
  command: 'go test ./...'
`
	if err := os.WriteFile("tatsu.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tatsu.yaml")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify config values
	if cfg.Agent.Command != `opencode run "%s"` {
		t.Errorf("Agent.Command = %q, want %q", cfg.Agent.Command, `opencode run "%s"`)
	}
	if cfg.Validate.Command != "go test ./..." {
		t.Errorf("Validate.Command = %q, want %q", cfg.Validate.Command, "go test ./...")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for missing file, got nil")
	}

	// Verify error message mentions file not found
	if !strings.Contains(err.Error(), "tatsu.yaml not found") {
		t.Errorf("expected error to contain 'tatsu.yaml not found', got: %v", err)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	// Create invalid YAML
	content := `agent:
  command: 'test
invalid yaml here
`
	if err := os.WriteFile("tatsu.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for invalid YAML, got nil")
	}
}

func TestLoad_MissingAgentCommand(t *testing.T) {
	// Create config with missing agent command
	content := `validate:
  command: 'go test ./...'
`
	if err := os.WriteFile("tatsu.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for missing agent.command, got nil")
	}
	if err.Error() != "agent.command is required in tatsu.yaml" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoad_MissingValidateCommand(t *testing.T) {
	// Create config with missing validate command
	content := `agent:
  command: 'opencode run "%s"'
`
	if err := os.WriteFile("tatsu.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for missing validate.command, got nil")
	}
	if err.Error() != "validate.command is required in tatsu.yaml" {
		t.Errorf("unexpected error: %v", err)
	}
}
