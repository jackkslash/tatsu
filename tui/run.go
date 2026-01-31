package tui

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jack/tatsu/config"
	"github.com/jack/tatsu/harness"
	"github.com/jack/tatsu/prd"
)

// RunTaskInTUI runs a single task and sends progress messages to the TUI.
// send is program.Send; call from a goroutine.
func RunTaskInTUI(send func(tea.Msg), cfg *config.Config, maxIter int, task string) {
	for i := 1; i <= maxIter; i++ {
		send(iterationStartMsg{iter: i, maxIter: maxIter})

		// Run agent
		if err := runAgentCapture(send, cfg, task); err != nil {
			send(agentErrorMsg{err: err.Error()})
		}

		// Validate
		send(validationStartMsg{})
		success, output := runValidate(cfg)
		send(validationResultMsg{success: success, output: output})

		if success {
			send(runCompleteMsg{success: true})
			return
		}
	}
	send(runCompleteMsg{success: false, errMsg: "max iterations reached"})
}

// RunPRDInTUI runs a PRD file and sends progress messages to the TUI.
func RunPRDInTUI(send func(tea.Msg), cfg *config.Config, maxIter int, prdPath string) {
	doc, err := prd.LoadPRD(prdPath)
	if err != nil {
		send(runCompleteMsg{success: false, errMsg: err.Error()})
		return
	}
	incomplete := doc.IncompleteTasks()
	if len(incomplete) == 0 {
		send(runCompleteMsg{success: true})
		return
	}

	for idx, task := range incomplete {
		send(prdTaskStartMsg{current: idx + 1, total: len(incomplete), title: task.Title})
		if err := runTaskLoop(send, cfg, maxIter, task.Title); err != nil {
			send(runCompleteMsg{success: false, errMsg: err.Error()})
			return
		}
	}
	send(runCompleteMsg{success: true})
}

func runTaskLoop(send func(tea.Msg), cfg *config.Config, maxIter int, task string) error {
	for i := 1; i <= maxIter; i++ {
		send(iterationStartMsg{iter: i, maxIter: maxIter})
		if err := runAgentCapture(send, cfg, task); err != nil {
			send(agentErrorMsg{err: err.Error()})
		}
		send(validationStartMsg{})
		success, output := runValidate(cfg)
		send(validationResultMsg{success: success, output: output})
		if success {
			return nil
		}
	}
	return fmt.Errorf("max iterations reached")
}

func runAgentCapture(send func(tea.Msg), cfg *config.Config, task string) error {
	cmdStr := fmt.Sprintf(cfg.Agent.Command, escapeTask(task))
	c := exec.Command("bash", "-c", cmdStr)
	c.Env = harness.AgentEnv()
	stdout, _ := c.StdoutPipe()
	stderr, _ := c.StderrPipe()
	if err := c.Start(); err != nil {
		return err
	}
	go captureLines(send, stdout)
	go captureLines(send, stderr)
	return c.Wait()
}

func captureLines(send func(tea.Msg), r io.ReadCloser) {
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		send(agentOutputMsg{line: scanner.Text()})
	}
}

func runValidate(cfg *config.Config) (bool, string) {
	c := exec.Command("bash", "-c", cfg.Validate.Command)
	out, err := c.CombinedOutput()
	s := string(out)
	if err != nil {
		return false, s
	}
	return true, s
}
func escapeTask(task string) string {
	return strings.ReplaceAll(task, `"`, `\"`)
}

