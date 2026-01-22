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

func TestGenerate_FileDoesNotExist(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")
	defer os.Remove("tatsu.yaml")

	// Generate config
	if err := Generate(false); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat("tatsu.yaml"); os.IsNotExist(err) {
		t.Fatal("tatsu.yaml was not created")
	}

	// Verify we can load it
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify default agent command
	if cfg.Agent.Command != `opencode run "%s"` {
		t.Errorf("Agent.Command = %q, want %q", cfg.Agent.Command, `opencode run "%s"`)
	}
}

func TestGenerate_FileExistsWithoutForce(t *testing.T) {
	// Create existing config file
	content := `agent:
  command: 'opencode run "%s"'
validate:
  command: 'go test ./...'
`
	if err := os.WriteFile("tatsu.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tatsu.yaml")

	// Attempt to generate without force
	err := Generate(false)
	if err == nil {
		t.Fatal("Generate() expected error when file exists, got nil")
	}

	// Verify error message mentions --force
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected error to contain 'already exists', got: %v", err)
	}
}

func TestGenerate_FileExistsWithForce(t *testing.T) {
	// Create existing config file
	oldContent := `agent:
  command: 'old command'
validate:
  command: 'old test'
`
	if err := os.WriteFile("tatsu.yaml", []byte(oldContent), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tatsu.yaml")

	// Generate with force
	if err := Generate(true); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify file was overwritten and can be loaded
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify it has the default agent command (not the old one)
	if cfg.Agent.Command != `opencode run "%s"` {
		t.Errorf("Agent.Command = %q, want %q", cfg.Agent.Command, `opencode run "%s"`)
	}
}

func TestGenerate_DetectsGoProject(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")
	defer os.Remove("tatsu.yaml")
	defer os.Remove("go.mod")

	// Create go.mod to simulate Go project
	if err := os.WriteFile("go.mod", []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Generate config
	if err := Generate(false); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Load and verify
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Validate.Command != "go test ./..." {
		t.Errorf("Validate.Command = %q, want %q", cfg.Validate.Command, "go test ./...")
	}
}

func TestGenerate_DetectsNodeProject(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")
	defer os.Remove("tatsu.yaml")
	defer os.Remove("package.json")

	// Create package.json to simulate Node project
	if err := os.WriteFile("package.json", []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Generate config
	if err := Generate(false); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Load and verify
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Validate.Command != "npm test" {
		t.Errorf("Validate.Command = %q, want %q", cfg.Validate.Command, "npm test")
	}
}

func TestGenerate_DefaultFallback(t *testing.T) {
	// Ensure no config file exists and no project markers
	os.Remove("tatsu.yaml")
	defer os.Remove("tatsu.yaml")
	os.Remove("go.mod")
	os.Remove("package.json")
	os.Remove("requirements.txt")
	os.Remove("pyproject.toml")
	os.Remove("Cargo.toml")

	// Generate config
	if err := Generate(false); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Load and verify
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Should have default fallback message
	expected := "echo 'No tests configured. Update tatsu.yaml with your test command.'"
	if cfg.Validate.Command != expected {
		t.Errorf("Validate.Command = %q, want %q", cfg.Validate.Command, expected)
	}
}
