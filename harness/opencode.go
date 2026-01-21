package harness

import "os/exec"

type OpenCodeHarness struct {
	command string
}

func NewOpenCodeHarness() *OpenCodeHarness {
	return &OpenCodeHarness{
		command: "opencode",
	}
}

func (h *OpenCodeHarness) Name() string {
	return "OpenCode"
}

func (h *OpenCodeHarness) IsAvailable() bool {
	cmd := exec.Command(h.command, "--version")
	return cmd.Run() == nil
}
