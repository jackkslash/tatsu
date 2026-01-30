package tui

import (
	"fmt"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jack/tatsu/config"
	"github.com/jack/tatsu/harness"
)

// Mode is either "task" or "prd"
type Mode string

const (
	ModeTask Mode = "task"
	ModePRD  Mode = "prd"
)

type appState int

const (
	stateInput appState = iota
	stateRunning
	stateDone
)

// Messages from run goroutine
type iterationStartMsg struct {
	iter    int
	maxIter int
}

type agentOutputMsg struct {
	line string
}

type agentErrorMsg struct {
	err string
}

type validationStartMsg struct{}

type validationResultMsg struct {
	success bool
	output  string
}

type runCompleteMsg struct {
	success bool
	errMsg  string
}

type prdTaskStartMsg struct {
	current int
	total   int
	title   string
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	outputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
)

type model struct {
	// input state
	mode   Mode
	input  string
	width  int
	height int
	state  appState

	// run context (set when starting run)
	cfg        *config.Config
	harness    harness.Harness
	maxIter    int
	send       func(tea.Msg)
	runMode    Mode
	runInput   string
	prdCurrent int
	prdTotal   int
	prdTitle   string

	// running state
	currentIter     int
	maxIterations   int
	agentOutput     []string
	validationOutput string
	agentError     string
	status         string

	// done state
	runSuccess bool
	runErr     string
	scrollOffset int // for scrolling through output
}

// NewModel creates a TUI model. send is program.Send; set before Run().
func NewModel(cfg *config.Config, h harness.Harness, maxIter int) *model {
	return &model{
		mode:         ModeTask,
		input:        "",
		state:        stateInput,
		cfg:          cfg,
		harness:      h,
		maxIter:      maxIter,
		agentOutput:  make([]string, 0),
	}
}

func (m *model) setSend(send func(tea.Msg)) {
	m.send = send
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case iterationStartMsg:
		m.currentIter = msg.iter
		m.maxIterations = msg.maxIter
		m.status = "running agent"
		m.agentOutput = nil
		m.agentError = ""
		return m, nil

	case agentOutputMsg:
		m.agentOutput = append(m.agentOutput, msg.line)
		if len(m.agentOutput) > 100 {
			m.agentOutput = m.agentOutput[len(m.agentOutput)-100:]
		}
		return m, nil

	case agentErrorMsg:
		m.agentError = msg.err
		return m, nil

	case validationStartMsg:
		m.status = "validating"
		return m, nil

	case validationResultMsg:
		m.validationOutput = msg.output
		if msg.success {
			m.status = "success"
		} else {
			m.status = "validation failed"
		}
		return m, nil

	case runCompleteMsg:
		m.runSuccess = msg.success
		m.runErr = msg.errMsg
		m.state = stateDone
		m.scrollOffset = 0
		return m, nil

	case prdTaskStartMsg:
		m.prdCurrent = msg.current
		m.prdTotal = msg.total
		m.prdTitle = msg.title
		return m, nil

	case tea.KeyMsg:
		s := msg.String()
		if m.state == stateInput {
			return m.handleInputKey(s)
		}
		if m.state == stateDone {
			switch s {
			case "enter", "r":
				m.state = stateInput
				m.runSuccess = false
				m.runErr = ""
				m.agentOutput = nil
				m.validationOutput = ""
				m.prdTitle = ""
				m.scrollOffset = 0
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.scrollOffset > 0 {
					m.scrollOffset--
				}
				return m, nil
			case "down", "j":
				m.scrollOffset++
				return m, nil
			}
			return m, nil
		}
		// running: only allow quit
		if s == "q" || s == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m *model) handleInputKey(s string) (tea.Model, tea.Cmd) {
	switch s {
	case "tab", "left", "right":
		if m.mode == ModeTask {
			m.mode = ModePRD
		} else {
			m.mode = ModeTask
		}
		return m, nil

	case "enter":
		in := strings.TrimSpace(m.input)
		if in == "" {
			return m, nil
		}
		m.runMode = m.mode
		m.runInput = in
		m.state = stateRunning
		m.currentIter = 0
		m.maxIterations = m.maxIter
		m.agentOutput = nil
		m.validationOutput = ""
		m.agentError = ""
		m.status = "starting..."
		if m.runMode == ModeTask {
			go RunTaskInTUI(m.send, m.cfg, m.harness, m.maxIter, in)
		} else {
			go RunPRDInTUI(m.send, m.cfg, m.harness, m.maxIter, in)
		}
		return m, nil

	case "q", "ctrl+c":
		return m, tea.Quit

	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
		return m, nil

	default:
		runes := []rune(s)
		if len(runes) == 1 {
			r := runes[0]
			if unicode.IsPrint(r) || r == ' ' {
				m.input += s
			}
		}
		return m, nil
	}
}

func (m *model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}
	switch m.state {
	case stateInput:
		return m.viewInput()
	case stateRunning:
		return m.viewRunning()
	case stateDone:
		return m.viewDone()
	}
	return ""
}

func (m *model) viewInput() string {
	var sections []string
	title := titleStyle.Render("Tatsu")
	sections = append(sections, title)
	sections = append(sections, "")
	if m.mode == ModeTask {
		sections = append(sections, labelStyle.Render("Task description:"))
	} else {
		sections = append(sections, labelStyle.Render("PRD file path:"))
	}
	sections = append(sections, helpStyle.Render("Tab to switch mode"))
	sections = append(sections, "")
	sections = append(sections, "  "+m.input+"▌")
	sections = append(sections, "")
	sections = append(sections, helpStyle.Render("Enter to run • q to quit"))
	content := lipgloss.JoinVertical(lipgloss.Center, sections...)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m *model) viewRunning() string {
	var sections []string
	sections = append(sections, titleStyle.Render("Tatsu"))
	sections = append(sections, "")
	if m.prdTotal > 0 {
		sections = append(sections, labelStyle.Render(fmt.Sprintf("PRD task %d/%d: %s", m.prdCurrent, m.prdTotal, m.prdTitle)))
		sections = append(sections, "")
	}
	sections = append(sections, labelStyle.Render(fmt.Sprintf("🔁 Iteration %d/%d • %s", m.currentIter, m.maxIterations, m.status)))
	sections = append(sections, "")
	if m.agentError != "" {
		sections = append(sections, errorStyle.Render("Agent error: "+m.agentError))
		sections = append(sections, "")
	}
	if len(m.agentOutput) > 0 {
		start := len(m.agentOutput) - 15
		if start < 0 {
			start = 0
		}
		lines := strings.Join(m.agentOutput[start:], "\n")
		sections = append(sections, outputBoxStyle.Width(m.width-4).Render("Agent output:\n"+lines))
		sections = append(sections, "")
	}
	if m.validationOutput != "" {
		valLines := m.validationOutput
		if len(valLines) > 500 {
			valLines = valLines[len(valLines)-500:]
		}
		sections = append(sections, outputBoxStyle.Width(m.width-4).Render("Validation:\n"+valLines))
	}
	sections = append(sections, "")
	sections = append(sections, helpStyle.Render("q to quit"))
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m *model) viewDone() string {
	// Build full output content (same as running view) so user can scroll
	var sections []string
	sections = append(sections, titleStyle.Render("Tatsu")+" — completed")
	sections = append(sections, "")
	if m.prdTotal > 0 {
		sections = append(sections, labelStyle.Render(fmt.Sprintf("PRD task %d/%d: %s", m.prdCurrent, m.prdTotal, m.prdTitle)))
		sections = append(sections, "")
	}
	sections = append(sections, labelStyle.Render(fmt.Sprintf("🔁 Iteration %d/%d • %s", m.currentIter, m.maxIterations, m.status)))
	sections = append(sections, "")
	if m.agentError != "" {
		sections = append(sections, errorStyle.Render("Agent error: "+m.agentError))
		sections = append(sections, "")
	}
	if len(m.agentOutput) > 0 {
		agentLines := strings.Join(m.agentOutput, "\n")
		sections = append(sections, outputBoxStyle.Width(m.width-4).Render("Agent output:\n"+agentLines))
		sections = append(sections, "")
	}
	if m.validationOutput != "" {
		sections = append(sections, outputBoxStyle.Width(m.width-4).Render("Validation:\n"+m.validationOutput))
	}
	// Join and split into lines for scrolling (output only; footer is fixed below)
	fullOutput := lipgloss.JoinVertical(lipgloss.Left, sections...)
	lines := strings.Split(fullOutput, "\n")
	footerHeight := 4
	visibleHeight := m.height - footerHeight
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	maxScroll := len(lines) - visibleHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	offset := m.scrollOffset
	if offset > maxScroll {
		offset = maxScroll
	}
	if offset < 0 {
		offset = 0
	}
	window := lines
	if len(lines) > visibleHeight {
		end := offset + visibleHeight
		if end > len(lines) {
			end = len(lines)
		}
		window = lines[offset:end]
	}
	scrollContent := strings.Join(window, "\n")
	// Fixed footer: result + key bindings
	var footer []string
	footer = append(footer, "")
	if m.runSuccess {
		footer = append(footer, successStyle.Render("✅ Done"))
	} else {
		footer = append(footer, errorStyle.Render("❌ "+m.runErr))
	}
	footer = append(footer, helpStyle.Render("↑/↓ j/k scroll • r run again • q quit"))
	content := scrollContent + "\n" + lipgloss.JoinVertical(lipgloss.Left, footer...)
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, content)
}

// Run starts the TUI. Config and harness must be loaded; execution happens inside the TUI.
func Run(cfg *config.Config, h harness.Harness, maxIter int) error {
	m := NewModel(cfg, h, maxIter)
	p := tea.NewProgram(m, tea.WithAltScreen())
	m.setSend(p.Send)
	_, err := p.Run()
	return err
}
