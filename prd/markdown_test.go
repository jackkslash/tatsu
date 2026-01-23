package prd

import (
	"testing"
)

func TestParseMarkdown_BasicTaskList(t *testing.T) {
	content := `## Tasks
- [ ] create authentication system
- [ ] add user dashboard
- [ ] implement dark mode
- [x] setup database (already done)
`

	prd, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if prd.TotalCount() != 4 {
		t.Errorf("TotalCount() = %d, want 4", prd.TotalCount())
	}

	if prd.CompletedCount() != 1 {
		t.Errorf("CompletedCount() = %d, want 1", prd.CompletedCount())
	}

	if len(prd.IncompleteTasks()) != 3 {
		t.Errorf("IncompleteTasks() count = %d, want 3", len(prd.IncompleteTasks()))
	}

	// Check first task
	if prd.Tasks[0].Title != "create authentication system" {
		t.Errorf("Tasks[0].Title = %q, want %q", prd.Tasks[0].Title, "create authentication system")
	}
	if prd.Tasks[0].Completed {
		t.Error("Tasks[0].Completed = true, want false")
	}

	// Check completed task
	if prd.Tasks[3].Title != "setup database (already done)" {
		t.Errorf("Tasks[3].Title = %q, want %q", prd.Tasks[3].Title, "setup database (already done)")
	}
	if !prd.Tasks[3].Completed {
		t.Error("Tasks[3].Completed = false, want true")
	}
}

func TestParseMarkdown_UppercaseX(t *testing.T) {
	content := `- [X] completed task
- [ ] incomplete task
`

	prd, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if !prd.Tasks[0].Completed {
		t.Error("Tasks[0].Completed = false, want true (uppercase X should be completed)")
	}

	if prd.Tasks[1].Completed {
		t.Error("Tasks[1].Completed = true, want false")
	}
}

func TestParseMarkdown_WithIndentation(t *testing.T) {
	content := `  - [ ] indented task
    - [x] nested task
- [ ] regular task
`

	prd, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if prd.TotalCount() != 3 {
		t.Errorf("TotalCount() = %d, want 3", prd.TotalCount())
	}

	// Check that indentation is preserved in title
	if prd.Tasks[0].Title != "indented task" {
		t.Errorf("Tasks[0].Title = %q, want %q", prd.Tasks[0].Title, "indented task")
	}
}

func TestParseMarkdown_EmptyContent(t *testing.T) {
	content := `## Tasks
No tasks here
`

	_, err := ParseMarkdown(content)
	if err == nil {
		t.Fatal("ParseMarkdown() expected error for empty task list, got nil")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected ParseError, got %T", err)
	}

	if parseErr.Message != "no tasks found in markdown" {
		t.Errorf("ParseError.Message = %q, want %q", parseErr.Message, "no tasks found in markdown")
	}
}

func TestParseMarkdown_EmptyTaskTitle(t *testing.T) {
	content := `- [ ] 
- [x] valid task
`

	prd, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	// Empty task should be skipped
	if prd.TotalCount() != 1 {
		t.Errorf("TotalCount() = %d, want 1 (empty task should be skipped)", prd.TotalCount())
	}

	if prd.Tasks[0].Title != "valid task" {
		t.Errorf("Tasks[0].Title = %q, want %q", prd.Tasks[0].Title, "valid task")
	}
}

func TestParseMarkdown_StarAndPlusBullets(t *testing.T) {
	content := `* [ ] star bullet
+ [x] plus bullet completed
`

	prd, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if prd.TotalCount() != 2 {
		t.Errorf("TotalCount() = %d, want 2", prd.TotalCount())
	}

	if prd.Tasks[0].Title != "star bullet" {
		t.Errorf("Tasks[0].Title = %q, want %q", prd.Tasks[0].Title, "star bullet")
	}

	if prd.Tasks[1].Title != "plus bullet completed" {
		t.Errorf("Tasks[1].Title = %q, want %q", prd.Tasks[1].Title, "plus bullet completed")
	}

	if !prd.Tasks[1].Completed {
		t.Error("Tasks[1].Completed = false, want true")
	}
}

func TestParseMarkdown_MixedContent(t *testing.T) {
	content := `# Project PRD

This is a description.

## Tasks
- [ ] first task
- [x] second task (done)

Some other content here.

- [ ] third task
`

	prd, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if prd.TotalCount() != 3 {
		t.Errorf("TotalCount() = %d, want 3", prd.TotalCount())
	}

	// Should only extract task list items, ignoring other content
	if prd.Tasks[0].Title != "first task" {
		t.Errorf("Tasks[0].Title = %q, want %q", prd.Tasks[0].Title, "first task")
	}
}
