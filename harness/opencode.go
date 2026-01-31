package harness

import (
	"os"
	"os/exec"
)

// AgentEnv returns environment variables for non-interactive OpenCode runs.
// OPENCODE_CONFIG_CONTENT with explicit permission rules ensures edit, write,
// bash, etc. run without approval. CI=true signals headless/automated mode.
func AgentEnv() []string {
	return append(os.Environ(),
		`OPENCODE_CONFIG_CONTENT={"permission":{"*":"allow","edit":"allow","write":"allow","bash":"allow","read":"allow","list":"allow","glob":"allow"}}`,
		"CI=true",
	)
}

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
