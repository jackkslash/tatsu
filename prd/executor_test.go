package prd

import (
	"os"
	"testing"

	"github.com/jack/tatsu/config"
	"github.com/jack/tatsu/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// quietTest redirects stdout/stderr to /dev/null for the test
func quietTest(t *testing.T, fn func()) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	devNull, err := os.Open("/dev/null")
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()
	os.Stdout = devNull
	os.Stderr = devNull
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()
	fn()
}

type mockHarness struct{}

func (m *mockHarness) Name() string      { return "MockHarness" }
func (m *mockHarness) IsAvailable() bool { return true }

func TestNewExecutor(t *testing.T) {
	cfg := &config.Config{}
	h := &mockHarness{}
	r := runner.New(cfg, h)

	executor := NewExecutor(r)
	require.NotNil(t, executor)
	assert.Equal(t, r, executor.runner)
}

func TestExecutePRD_AllCompleted(t *testing.T) {
	cfg := &config.Config{}
	cfg.Agent.Command = "echo 'Agent: %s' >/dev/null"
	cfg.Validate.Command = "exit 0"

	r := runner.New(cfg, &mockHarness{})
	executor := NewExecutor(r)

	prd := &PRD{
		Tasks: []Task{
			{Title: "task 1", Completed: true},
			{Title: "task 2", Completed: true},
		},
	}

	var err error
	quietTest(t, func() {
		err = executor.ExecutePRD(prd, "")
	})
	require.NoError(t, err)
}

func TestExecutePRD_ExecuteIncompleteTasks(t *testing.T) {
	cfg := &config.Config{}
	cfg.Agent.Command = "echo 'Agent: %s' >/dev/null"
	cfg.Validate.Command = "exit 0" // Always pass validation

	r := runner.New(cfg, &mockHarness{})
	executor := NewExecutor(r)

	prd := &PRD{
		Tasks: []Task{
			{Title: "completed task", Completed: true},
			{Title: "incomplete task 1", Completed: false},
			{Title: "incomplete task 2", Completed: false},
		},
	}

	var err error
	quietTest(t, func() {
		err = executor.ExecutePRD(prd, "")
	})
	require.NoError(t, err)
}

func TestExecutePRD_SkipsCompletedTasks(t *testing.T) {
	cfg := &config.Config{}
	cfg.Agent.Command = "echo 'Agent: %s' >/dev/null"
	cfg.Validate.Command = "exit 0"

	r := runner.New(cfg, &mockHarness{})
	executor := NewExecutor(r)

	prd := &PRD{
		Tasks: []Task{
			{Title: "task 1", Completed: true},
			{Title: "task 2", Completed: false},
			{Title: "task 3", Completed: true},
			{Title: "task 4", Completed: false},
		},
	}

	// Should only execute task 2 and task 4
	var err error
	quietTest(t, func() {
		err = executor.ExecutePRD(prd, "")
	})
	require.NoError(t, err)
}

func TestExecutePRD_TaskFailure(t *testing.T) {
	cfg := &config.Config{}
	cfg.Agent.Command = "echo 'Agent: %s' >/dev/null"
	cfg.Validate.Command = "exit 1" // Always fail validation

	r := runner.New(cfg, &mockHarness{})
	executor := NewExecutor(r)

	prd := &PRD{
		Tasks: []Task{
			{Title: "failing task", Completed: false},
		},
	}

	// Should return error after max iterations
	var err error
	quietTest(t, func() {
		err = executor.ExecutePRD(prd, "")
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failing task")
	assert.Contains(t, err.Error(), "failed")
}
