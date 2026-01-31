package prd

import (
	"fmt"

	"github.com/jack/tatsu/runner"
)

// Executor executes tasks from a PRD
type Executor struct {
	runner *runner.Runner
}

// NewExecutor creates a new PRD executor
func NewExecutor(r *runner.Runner) *Executor {
	return &Executor{
		runner: r,
	}
}

// ExecutePRD executes all incomplete tasks from a PRD sequentially.
// If filename is non-empty, the PRD file is updated to mark each task complete as it succeeds.
func (e *Executor) ExecutePRD(prd *PRD, filename string) error {
	incomplete := prd.IncompleteTasks()

	if len(incomplete) == 0 {
		fmt.Println("✅ All tasks are already completed!")
		return nil
	}

	fmt.Printf("📋 PRD Summary:\n")
	fmt.Printf("   Total tasks: %d\n", prd.TotalCount())
	fmt.Printf("   Completed: %d\n", prd.CompletedCount())
	fmt.Printf("   Remaining: %d\n\n", len(incomplete))

	// Execute each incomplete task
	for i, task := range incomplete {
		fmt.Printf("📌 Task %d/%d: %s\n\n", i+1, len(incomplete), task.Title)

		// Execute task using runner
		if err := e.runner.Run(task.Title); err != nil {
			return fmt.Errorf("task '%s' failed: %w", task.Title, err)
		}

		// Mark task complete in PRD file
		if filename != "" && task.LineNum > 0 {
			if err := MarkTaskCompleteInFile(filename, task.LineNum); err != nil {
				fmt.Printf("⚠️  Failed to update PRD file: %v\n", err)
			}
		}

		fmt.Println()
	}

	fmt.Println("✅ All PRD tasks completed successfully!")
	return nil
}
