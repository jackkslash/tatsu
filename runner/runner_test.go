package runner

import (
	"testing"

	"github.com/jack/tatsu/config"
)

type mockHarness struct{}

func (m *mockHarness) Name() string        { return "MockHarness" }
func (m *mockHarness) IsAvailable() bool   { return true }

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	h := &mockHarness{}
	
	r := New(cfg, h)
	
	if r == nil {
		t.Fatal("New() returned nil")
	}
	if r.config != cfg {
		t.Error("New() did not set config correctly")
	}
	if r.harness != h {
		t.Error("New() did not set harness correctly")
	}
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
			if got != tt.expected {
				t.Errorf("escapeTask(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMaxIterations(t *testing.T) {
	if MaxIterations != 15 {
		t.Errorf("MaxIterations = %d, want 15", MaxIterations)
	}
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
	if err != nil {
		t.Errorf("Run() with passing validation returned error: %v", err)
	}
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
	if err == nil {
		t.Error("Run() with failing validation should return error")
	}
	if err.Error() != "max iterations reached" {
		t.Errorf("unexpected error: %v", err)
	}
}
