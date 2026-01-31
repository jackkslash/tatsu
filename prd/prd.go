package prd

import "fmt"

// Task represents a single task in a PRD
type Task struct {
	Title     string
	Completed bool
	LineNum   int // 1-based line number in file (0 if unknown)
}

// PRD represents a Product Requirements Document containing tasks
type PRD struct {
	Tasks []Task
}

// Validate checks if a PRD is valid
func (p *PRD) Validate() error {
	if len(p.Tasks) == 0 {
		return fmt.Errorf("PRD must contain at least one task")
	}

	for i, task := range p.Tasks {
		if task.Title == "" {
			return fmt.Errorf("task at index %d has empty title", i)
		}
	}

	return nil
}

// IncompleteTasks returns a slice of tasks that are not completed
func (p *PRD) IncompleteTasks() []Task {
	var incomplete []Task
	for _, task := range p.Tasks {
		if !task.Completed {
			incomplete = append(incomplete, task)
		}
	}
	return incomplete
}

// CompletedCount returns the number of completed tasks
func (p *PRD) CompletedCount() int {
	count := 0
	for _, task := range p.Tasks {
		if task.Completed {
			count++
		}
	}
	return count
}

// TotalCount returns the total number of tasks
func (p *PRD) TotalCount() int {
	return len(p.Tasks)
}
