package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jack/tatsu/config"
	"github.com/jack/tatsu/harness"
)

const DefaultMaxIterations = 15

type Runner struct {
	config       *config.Config
	harness      harness.Harness
	maxIterations int
}

func New(cfg *config.Config, h harness.Harness) *Runner {
	return &Runner{
		config:       cfg,
		harness:      h,
		maxIterations: DefaultMaxIterations,
	}
}

// NewWithMaxIterations creates a Runner with custom max iterations
func NewWithMaxIterations(cfg *config.Config, h harness.Harness, maxIter int) *Runner {
	return &Runner{
		config:       cfg,
		harness:      h,
		maxIterations: maxIter,
	}
}

func (r *Runner) Run(task string) error {
	for i := 1; i <= r.maxIterations; i++ {
		fmt.Printf("🔁 Iteration %d/%d\n", i, r.maxIterations)

		// Run agent
		if err := r.runAgent(task); err != nil {
			fmt.Printf("⚠️  Agent error: %v\n", err)
		}

		// Validate
		if ok := r.validate(); ok {
			fmt.Println("\n✅ Task completed successfully!")
			return nil
		}

		fmt.Println("❌ Validation failed, retrying...")
		fmt.Println()
	}

	return fmt.Errorf("max iterations reached")
}

func (r *Runner) runAgent(task string) error {
	// Format command with task
	cmd := fmt.Sprintf(r.config.Agent.Command, EscapeTask(task))

	// Execute with non-interactive env (permission allow, CI)
	c := exec.Command("bash", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = nil // /dev/null - prevent blocking on stdin
	c.Env = harness.AgentEnv()

	return c.Run()
}

func (r *Runner) validate() bool {
	// Execute validation command
	c := exec.Command("bash", "-c", r.config.Validate.Command)
	output, err := c.CombinedOutput()

	if err != nil {
		fmt.Printf("\n📋 Validation output:\n%s\n", string(output))
		return false
	}

	return true
}

// EscapeTask escapes a task string for safe use in shell commands.
func EscapeTask(task string) string {
	return strings.ReplaceAll(task, `"`, `\"`)
}
