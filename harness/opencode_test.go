package harness

import (
	"testing"
)

func TestNewOpenCodeHarness(t *testing.T) {
	h := NewOpenCodeHarness()
	if h == nil {
		t.Fatal("NewOpenCodeHarness() returned nil")
	}
	if h.command != "opencode" {
		t.Errorf("command = %q, want %q", h.command, "opencode")
	}
}

func TestOpenCodeHarness_Name(t *testing.T) {
	h := NewOpenCodeHarness()
	if got := h.Name(); got != "OpenCode" {
		t.Errorf("Name() = %q, want %q", got, "OpenCode")
	}
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
