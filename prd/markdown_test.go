package prd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMarkdown_BasicTaskList(t *testing.T) {
	content := `## Tasks
- [ ] create authentication system
- [ ] add user dashboard
- [ ] implement dark mode
- [x] setup database (already done)
`

	prd, err := ParseMarkdown(content)
	require.NoError(t, err)

	assert.Equal(t, 4, prd.TotalCount())
	assert.Equal(t, 1, prd.CompletedCount())
	assert.Len(t, prd.IncompleteTasks(), 3)

	// Check first task
	assert.Equal(t, "create authentication system", prd.Tasks[0].Title)
	assert.False(t, prd.Tasks[0].Completed)

	// Check completed task
	assert.Equal(t, "setup database (already done)", prd.Tasks[3].Title)
	assert.True(t, prd.Tasks[3].Completed)
}

func TestParseMarkdown_UppercaseX(t *testing.T) {
	content := `- [X] completed task
- [ ] incomplete task
`

	prd, err := ParseMarkdown(content)
	require.NoError(t, err)

	assert.True(t, prd.Tasks[0].Completed, "uppercase X should be completed")
	assert.False(t, prd.Tasks[1].Completed)
}

func TestParseMarkdown_WithIndentation(t *testing.T) {
	content := `  - [ ] indented task
    - [x] nested task
- [ ] regular task
`

	prd, err := ParseMarkdown(content)
	require.NoError(t, err)

	assert.Equal(t, 3, prd.TotalCount())
	assert.Equal(t, "indented task", prd.Tasks[0].Title)
}

func TestParseMarkdown_EmptyContent(t *testing.T) {
	content := `## Tasks
No tasks here
`

	_, err := ParseMarkdown(content)
	require.Error(t, err)

	parseErr, ok := err.(*ParseError)
	require.True(t, ok, "expected ParseError")
	assert.Equal(t, "no tasks found in markdown", parseErr.Message)
}

func TestParseMarkdown_EmptyTaskTitle(t *testing.T) {
	content := `- [ ] 
- [x] valid task
`

	prd, err := ParseMarkdown(content)
	require.NoError(t, err)

	// Empty task should be skipped
	assert.Equal(t, 1, prd.TotalCount(), "empty task should be skipped")
	assert.Equal(t, "valid task", prd.Tasks[0].Title)
}

func TestParseMarkdown_StarAndPlusBullets(t *testing.T) {
	content := `* [ ] star bullet
+ [x] plus bullet completed
`

	prd, err := ParseMarkdown(content)
	require.NoError(t, err)

	assert.Equal(t, 2, prd.TotalCount())
	assert.Equal(t, "star bullet", prd.Tasks[0].Title)
	assert.Equal(t, "plus bullet completed", prd.Tasks[1].Title)
	assert.True(t, prd.Tasks[1].Completed)
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
	require.NoError(t, err)

	assert.Equal(t, 3, prd.TotalCount())
	// Should only extract task list items, ignoring other content
	assert.Equal(t, "first task", prd.Tasks[0].Title)
}

func TestMarkTaskCompleteInFile(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "prd.md")
	content := `# PRD

- [ ] first task
- [ ] second task
- [x] already done
`
	require.NoError(t, os.WriteFile(filename, []byte(content), 0644))

	// Line 3 = "- [ ] first task"
	err := MarkTaskCompleteInFile(filename, 3)
	require.NoError(t, err)

	data, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Contains(t, string(data), "- [x] first task")
	assert.Contains(t, string(data), "- [ ] second task")
	assert.Contains(t, string(data), "- [x] already done")

	// Line 4 = "second task"
	err = MarkTaskCompleteInFile(filename, 4)
	require.NoError(t, err)
	data, err = os.ReadFile(filename)
	require.NoError(t, err)
	assert.Contains(t, string(data), "- [x] second task")

	// Line 5 = "already done" - no change (already [x])
	err = MarkTaskCompleteInFile(filename, 5)
	require.NoError(t, err)

	// Out of range: no error, no change
	err = MarkTaskCompleteInFile(filename, 99)
	require.NoError(t, err)
}
