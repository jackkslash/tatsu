package prd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadPRD_Success(t *testing.T) {
	// Create temporary PRD file
	content := `## Tasks
- [ ] first task
- [x] completed task
- [ ] second task
`
	filename := "test_prd.md"
	require.NoError(t, os.WriteFile(filename, []byte(content), 0644))
	defer os.Remove(filename)

	// Load PRD
	prd, err := LoadPRD(filename)
	require.NoError(t, err)

	// Verify PRD was loaded correctly
	assert.Equal(t, 3, prd.TotalCount())
	assert.Equal(t, 1, prd.CompletedCount())
	assert.Equal(t, "first task", prd.Tasks[0].Title)
	assert.True(t, prd.Tasks[1].Completed)
}

func TestLoadPRD_FileNotFound(t *testing.T) {
	// Ensure file doesn't exist
	os.Remove("nonexistent.md")

	// Attempt to load
	_, err := LoadPRD("nonexistent.md")
	require.Error(t, err)
	assert.Equal(t, "PRD file not found: nonexistent.md", err.Error())
}

func TestLoadPRD_InvalidMarkdown(t *testing.T) {
	// Create file with no tasks
	content := `## Tasks
No tasks here
`
	filename := "empty_prd.md"
	require.NoError(t, os.WriteFile(filename, []byte(content), 0644))
	defer os.Remove(filename)

	// Attempt to load
	_, err := LoadPRD(filename)
	require.Error(t, err)
	assert.Equal(t, "failed to parse PRD file: no tasks found in markdown", err.Error())
}

func TestLoadPRD_EmptyFile(t *testing.T) {
	// Create empty file
	filename := "empty.md"
	require.NoError(t, os.WriteFile(filename, []byte(""), 0644))
	defer os.Remove(filename)

	// Attempt to load
	_, err := LoadPRD(filename)
	require.Error(t, err)
	// Should get parse error for no tasks
	assert.Equal(t, "failed to parse PRD file: no tasks found in markdown", err.Error())
}
