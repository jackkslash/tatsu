package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	// Create a temporary config file
	content := `agent:
  command: 'opencode run "%s"'
validate:
  command: 'go test ./...'
`
	require.NoError(t, os.WriteFile("tatsu.yaml", []byte(content), 0644))
	defer os.Remove("tatsu.yaml")

	// Load config
	cfg, err := Load()
	require.NoError(t, err)

	// Verify config values
	assert.Equal(t, `opencode run "%s"`, cfg.Agent.Command)
	assert.Equal(t, "go test ./...", cfg.Validate.Command)
}

func TestLoad_FileNotFound(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tatsu.yaml not found")
}

func TestLoad_InvalidYAML(t *testing.T) {
	// Create invalid YAML
	content := `agent:
  command: 'test
invalid yaml here
`
	require.NoError(t, os.WriteFile("tatsu.yaml", []byte(content), 0644))
	defer os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	require.Error(t, err)
}

func TestLoad_MissingAgentCommand(t *testing.T) {
	// Create config with missing agent command
	content := `validate:
  command: 'go test ./...'
`
	require.NoError(t, os.WriteFile("tatsu.yaml", []byte(content), 0644))
	defer os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	require.Error(t, err)
	assert.Equal(t, "agent.command is required in tatsu.yaml", err.Error())
}

func TestLoad_MissingValidateCommand(t *testing.T) {
	// Create config with missing validate command
	content := `agent:
  command: 'opencode run "%s"'
`
	require.NoError(t, os.WriteFile("tatsu.yaml", []byte(content), 0644))
	defer os.Remove("tatsu.yaml")

	// Attempt to load
	_, err := Load()
	require.Error(t, err)
	assert.Equal(t, "validate.command is required in tatsu.yaml", err.Error())
}

func TestGenerate_FileDoesNotExist(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")
	defer os.Remove("tatsu.yaml")

	// Generate config
	require.NoError(t, Generate(false))

	// Verify file was created
	_, err := os.Stat("tatsu.yaml")
	require.NoError(t, err, "tatsu.yaml was not created")

	// Verify we can load it
	cfg, err := Load()
	require.NoError(t, err)

	// Verify default agent command
	assert.Equal(t, `opencode run "%s"`, cfg.Agent.Command)
}

func TestGenerate_FileExistsWithoutForce(t *testing.T) {
	// Create existing config file
	content := `agent:
  command: 'opencode run "%s"'
validate:
  command: 'go test ./...'
`
	require.NoError(t, os.WriteFile("tatsu.yaml", []byte(content), 0644))
	defer os.Remove("tatsu.yaml")

	// Attempt to generate without force
	err := Generate(false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestGenerate_FileExistsWithForce(t *testing.T) {
	// Create existing config file
	oldContent := `agent:
  command: 'old command'
validate:
  command: 'old test'
`
	require.NoError(t, os.WriteFile("tatsu.yaml", []byte(oldContent), 0644))
	defer os.Remove("tatsu.yaml")

	// Generate with force
	require.NoError(t, Generate(true))

	// Verify file was overwritten and can be loaded
	cfg, err := Load()
	require.NoError(t, err)

	// Verify it has the default agent command (not the old one)
	assert.Equal(t, `opencode run "%s"`, cfg.Agent.Command)
}

func TestGenerate_DetectsGoProject(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")
	defer os.Remove("tatsu.yaml")
	defer os.Remove("go.mod")

	// Create go.mod to simulate Go project
	require.NoError(t, os.WriteFile("go.mod", []byte("module test\n"), 0644))

	// Generate config
	require.NoError(t, Generate(false))

	// Load and verify
	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "go test ./...", cfg.Validate.Command)
}

func TestGenerate_DetectsNodeProject(t *testing.T) {
	// Ensure no config file exists
	os.Remove("tatsu.yaml")
	defer os.Remove("tatsu.yaml")
	defer os.Remove("package.json")

	// Create package.json to simulate Node project
	require.NoError(t, os.WriteFile("package.json", []byte(`{"name": "test"}`), 0644))

	// Generate config
	require.NoError(t, Generate(false))

	// Load and verify
	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "npm test", cfg.Validate.Command)
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
	require.NoError(t, Generate(false))

	// Load and verify
	cfg, err := Load()
	require.NoError(t, err)

	// Should have default fallback message
	expected := "echo 'No tests configured. Update tatsu.yaml with your test command.'"
	assert.Equal(t, expected, cfg.Validate.Command)
}
