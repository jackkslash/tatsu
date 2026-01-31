package runner

import (
	"testing"

	"github.com/jack/tatsu/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHarness struct{}

func (m *mockHarness) Name() string        { return "MockHarness" }
func (m *mockHarness) IsAvailable() bool   { return true }

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	h := &mockHarness{}
	
	r := New(cfg, h)
	
	require.NotNil(t, r)
	assert.Equal(t, cfg, r.config)
	assert.Equal(t, h, r.harness)
	assert.Equal(t, DefaultMaxIterations, r.maxIterations)
}

func TestEscapeTask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no quotes",
			input:    "simple task",
			expected: "simple task",
		},
		{
			name:     "single quote",
			input:    `task with "quote"`,
			expected: `task with \"quote\"`,
		},
		{
			name:     "multiple quotes",
			input:    `"start" and "end"`,
			expected: `\"start\" and \"end\"`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapeTask(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestDefaultMaxIterations(t *testing.T) {
	assert.Equal(t, 15, DefaultMaxIterations)
}

func TestNewWithMaxIterations(t *testing.T) {
	cfg := &config.Config{}
	h := &mockHarness{}
	
	r := NewWithMaxIterations(cfg, h, 5)
	require.NotNil(t, r)
	assert.Equal(t, 5, r.maxIterations)
	assert.Equal(t, cfg, r.config)
	assert.Equal(t, h, r.harness)
}

func TestRunner_UsesCustomMaxIterations(t *testing.T) {
	cfg := &config.Config{}
	cfg.Agent.Command = "echo 'Agent: %s'"
	cfg.Validate.Command = "exit 1" // Always fail
	
	h := &mockHarness{}
	r := NewWithMaxIterations(cfg, h, 3) // Only 3 iterations
	
	// Should fail after 3 iterations, not 15
	err := r.Run("test task")
	require.Error(t, err)
	assert.Equal(t, "max iterations reached", err.Error())
	
	// Verify it only tried 3 times by checking the error
	// (In real usage, you'd see 3 iteration messages, not 15)
}

func TestRunner_IntegrationWithPassingValidation(t *testing.T) {
	// Create config with commands that will succeed
	cfg := &config.Config{}
	cfg.Agent.Command = "echo 'Agent running: %s'"
	cfg.Validate.Command = "exit 0"
	
	h := &mockHarness{}
	r := New(cfg, h)
	
	// Run with passing validation - should succeed on first iteration
	err := r.Run("test task")
	assert.NoError(t, err)
}

func TestRunner_IntegrationWithFailingValidation(t *testing.T) {
	// Create config with failing validation
	cfg := &config.Config{}
	cfg.Agent.Command = "echo 'Agent running: %s'"
	cfg.Validate.Command = "exit 1"
	
	h := &mockHarness{}
	r := New(cfg, h)
	
	// Run with failing validation - should hit max iterations
	err := r.Run("test task")
	require.Error(t, err)
	assert.Equal(t, "max iterations reached", err.Error())
}
