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
			got := escapeTask(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestMaxIterations(t *testing.T) {
	assert.Equal(t, 15, MaxIterations)
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
