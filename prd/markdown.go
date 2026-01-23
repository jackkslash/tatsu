package prd

import (
	"regexp"
	"strings"
)

var (
	// taskListItemRegex matches markdown task list items: "- [ ] task" or "- [x] task"
	taskListItemRegex = regexp.MustCompile(`^[\s]*[-*+][\s]+\[([\sxX])\][\s]+(.+)$`)
)

// ParseMarkdown parses a markdown PRD file and returns a PRD struct
func ParseMarkdown(content string) (*PRD, error) {
	lines := strings.Split(content, "\n")
	var tasks []Task

	for _, line := range lines {
		matches := taskListItemRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		// matches[1] is the checkbox state (space, 'x', or 'X')
		// matches[2] is the task title
		checkbox := strings.TrimSpace(matches[1])
		title := strings.TrimSpace(matches[2])

		// Skip if title is empty
		if title == "" {
			continue
		}

		// Determine if task is completed
		completed := checkbox == "x" || checkbox == "X"

		tasks = append(tasks, Task{
			Title:     title,
			Completed: completed,
		})
	}

	if len(tasks) == 0 {
		return nil, &ParseError{Message: "no tasks found in markdown"}
	}

	prd := &PRD{
		Tasks: tasks,
	}

	if err := prd.Validate(); err != nil {
		return nil, &ParseError{Message: err.Error()}
	}

	return prd, nil
}

// ParseError represents an error during PRD parsing
type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}
