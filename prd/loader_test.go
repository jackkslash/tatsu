package prd

import (
	"os"
	"testing"
)

func TestLoadPRD_Success(t *testing.T) {
	// Create temporary PRD file
	content := `## Tasks
- [ ] first task
- [x] completed task
- [ ] second task
`
	filename := "test_prd.md"
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename)

	// Load PRD
	prd, err := LoadPRD(filename)
	if err != nil {
		t.Fatalf("LoadPRD() error = %v", err)
	}

	// Verify PRD was loaded correctly
	if prd.TotalCount() != 3 {
		t.Errorf("TotalCount() = %d, want 3", prd.TotalCount())
	}

	if prd.CompletedCount() != 1 {
		t.Errorf("CompletedCount() = %d, want 1", prd.CompletedCount())
	}

	if prd.Tasks[0].Title != "first task" {
		t.Errorf("Tasks[0].Title = %q, want %q", prd.Tasks[0].Title, "first task")
	}

	if !prd.Tasks[1].Completed {
		t.Error("Tasks[1].Completed = false, want true")
	}
}

func TestLoadPRD_FileNotFound(t *testing.T) {
	// Ensure file doesn't exist
	os.Remove("nonexistent.md")

	// Attempt to load
	_, err := LoadPRD("nonexistent.md")
	if err == nil {
		t.Fatal("LoadPRD() expected error for missing file, got nil")
	}

	// Verify error message mentions file not found
	if err.Error() != "PRD file not found: nonexistent.md" {
		t.Errorf("error message = %q, want %q", err.Error(), "PRD file not found: nonexistent.md")
	}
}

func TestLoadPRD_InvalidMarkdown(t *testing.T) {
	// Create file with no tasks
	content := `## Tasks
No tasks here
`
	filename := "empty_prd.md"
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename)

	// Attempt to load
	_, err := LoadPRD(filename)
	if err == nil {
		t.Fatal("LoadPRD() expected error for invalid PRD, got nil")
	}

	// Verify error message mentions parse error
	if err.Error() != "failed to parse PRD file: no tasks found in markdown" {
		t.Errorf("error message = %q, want %q", err.Error(), "failed to parse PRD file: no tasks found in markdown")
	}
}

func TestLoadPRD_EmptyFile(t *testing.T) {
	// Create empty file
	filename := "empty.md"
	if err := os.WriteFile(filename, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename)

	// Attempt to load
	_, err := LoadPRD(filename)
	if err == nil {
		t.Fatal("LoadPRD() expected error for empty file, got nil")
	}

	// Should get parse error for no tasks
	if err.Error() != "failed to parse PRD file: no tasks found in markdown" {
		t.Errorf("error message = %q, want %q", err.Error(), "failed to parse PRD file: no tasks found in markdown")
	}
}
