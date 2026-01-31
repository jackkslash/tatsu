package prd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	checkboxIncomplete = " [ ] "
	checkboxComplete   = " [x] "
)

var (
	// taskListItemRegex matches markdown task list items: "- [ ] task" or "- [x] task"
	taskListItemRegex = regexp.MustCompile(`^[\s]*[-*+][\s]+\[([\sxX])\][\s]+(.+)$`)
)

// ParseMarkdown parses a markdown PRD file and returns a PRD struct
func ParseMarkdown(content string) (*PRD, error) {
	lines := strings.Split(content, "\n")
	var tasks []Task

	for i, line := range lines {
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
			LineNum:   i + 1,
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

// MarkTaskCompleteInFile marks a task as complete in the PRD file by changing [ ] to [x].
// lineNum is 1-based (same as editors).
func MarkTaskCompleteInFile(filename string, lineNum int) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read PRD file: %w", err)
	}
	lines := strings.Split(string(data), "\n")
	lineIndex := lineNum - 1

	if lineIndex >= 0 && lineIndex < len(lines) {
		lines[lineIndex] = strings.Replace(lines[lineIndex], checkboxIncomplete, checkboxComplete, 1)
		return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	}
	return nil
}
