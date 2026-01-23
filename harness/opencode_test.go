package harness

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenCodeHarness(t *testing.T) {
	h := NewOpenCodeHarness()
	require.NotNil(t, h)
	assert.Equal(t, "opencode", h.command)
}

func TestOpenCodeHarness_Name(t *testing.T) {
	h := NewOpenCodeHarness()
	assert.Equal(t, "OpenCode", h.Name())
}

func TestOpenCodeHarness_IsAvailable(t *testing.T) {
	h := NewOpenCodeHarness()
	
	// Note: This test will pass or fail based on whether opencode is actually installed
	// We're just testing that the method runs without panic
	available := h.IsAvailable()
	
	// Log the result for debugging
	t.Logf("OpenCode available: %v", available)
	
	// The test passes regardless - we're just checking it doesn't crash
	// In a real test suite, you might mock the exec.Command call
}

func TestOpenCodeHarness_ImplementsInterface(t *testing.T) {
	var _ Harness = (*OpenCodeHarness)(nil)
}
