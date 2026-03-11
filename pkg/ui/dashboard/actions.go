package dashboard

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jcleira/entire-local/pkg/checkpoint"
)

type entireAction int

const (
	actionStatus entireAction = iota
	actionExplain
	actionDoctor
	actionRewind
	actionResume
	actionReset
	actionClean
)

type actionItem struct {
	key                string
	label              string
	desc               string
	action             entireAction
	requiresCheckpoint bool
	destructive        bool
	execProcess        bool
}

var actionItems = []actionItem{
	{"s", "Status", "Show entire.io state", actionStatus, false, false, false},
	{"e", "Explain", "AI analysis of checkpoint", actionExplain, true, false, true},
	{"d", "Doctor", "Diagnose stuck sessions", actionDoctor, false, false, true},
	{"w", "Rewind", "Rewind to checkpoint", actionRewind, true, true, true},
	{"m", "Resume", "Resume session on branch", actionResume, true, true, true},
	{"x", "Reset", "Delete session state", actionReset, false, true, true},
	{"c", "Clean", "Remove orphaned data", actionClean, false, true, true},
}

func detectEntireCLI() tea.Cmd {
	return func() tea.Msg {
		_, err := exec.LookPath("entire")
		return entireCLIDetectedMsg{available: err == nil}
	}
}

func buildEntireArgs(action entireAction, detail *checkpoint.Checkpoint) []string {
	switch action {
	case actionStatus:
		return []string{"status", "--detailed"}
	case actionExplain:
		if detail != nil {
			return []string{"explain", "-c", detail.CheckpointID}
		}
		return []string{"explain"}
	case actionDoctor:
		return []string{"doctor"}
	case actionRewind:
		if detail != nil {
			return []string{"rewind", "--to", detail.CheckpointID}
		}
		return []string{"rewind"}
	case actionResume:
		if detail != nil {
			return []string{"resume", detail.Branch}
		}
		return []string{"resume"}
	case actionReset:
		return []string{"reset"}
	case actionClean:
		return []string{"clean"}
	}
	return nil
}

func executeEntireProcess(action entireAction, detail *checkpoint.Checkpoint) tea.Cmd {
	args := buildEntireArgs(action, detail)
	c := exec.Command("entire", args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return entireCmdFinishedMsg{err: err}
	})
}

func captureEntireOutput(action entireAction, detail *checkpoint.Checkpoint) tea.Cmd {
	return func() tea.Msg {
		args := buildEntireArgs(action, detail)
		c := exec.Command("entire", args...)
		out, err := c.CombinedOutput()
		if err != nil {
			return entireCmdOutputMsg{output: fmt.Sprintf("Error: %v\n%s", err, string(out))}
		}
		return entireCmdOutputMsg{output: string(out)}
	}
}

func (m model) handleActionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	switch {
	case msg.Type == tea.KeyEscape:
		m.showActions = false
		return m, nil

	case key.Matches(msg, keys.Down):
		if m.actionCursor < len(actionItems)-1 {
			m.actionCursor++
		}
		return m, nil

	case key.Matches(msg, keys.Up):
		if m.actionCursor > 0 {
			m.actionCursor--
		}
		return m, nil

	case key.Matches(msg, keys.Enter):
		return m.selectAction(actionItems[m.actionCursor])
	}

	for i, item := range actionItems {
		if msg.String() == item.key {
			m.actionCursor = i
			return m.selectAction(item)
		}
	}

	return m, nil
}

func (m model) selectAction(item actionItem) (tea.Model, tea.Cmd) {
	if !m.entireCLIAvailable {
		return m, nil
	}
	if item.requiresCheckpoint && m.detail == nil {
		return m, nil
	}
	if item.destructive {
		confirmed := item
		m.confirmAction = &confirmed
		m.showActions = false
		return m, nil
	}
	m.showActions = false
	return m.runAction(item)
}

func (m model) runAction(item actionItem) (tea.Model, tea.Cmd) {
	if item.execProcess {
		return m, executeEntireProcess(item.action, m.detail)
	}
	return m, captureEntireOutput(item.action, m.detail)
}

func renderActions(cursor int, cliAvailable, hasDetail bool, width, height int) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Actions"))
	sb.WriteString("\n")
	sb.WriteString(actionSepStyle.Render(strings.Repeat("─", 35)))
	sb.WriteString("\n")

	if !cliAvailable {
		sb.WriteString(warningStyle.Render("  ⚠ entire CLI not found in PATH"))
		sb.WriteString("\n")
	}

	for i, item := range actionItems {
		disabled := !cliAvailable || (item.requiresCheckpoint && !hasDetail)

		prefix := "   "
		if i == cursor {
			prefix = " ▸ "
		}

		keyPart := "[" + item.key + "]"

		sb.WriteString(prefix)
		if disabled {
			sb.WriteString(disabledStyle.Render(keyPart))
			sb.WriteString("  ")
			sb.WriteString(disabledStyle.Render(padRight(item.label, 10)))
			sb.WriteString(disabledStyle.Render(item.desc))
		} else {
			sb.WriteString(helpKeyStyle.Render(keyPart))
			sb.WriteString("  ")
			sb.WriteString(padRight(item.label, 10))
			sb.WriteString(dimStyle.Render(item.desc))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(actionSepStyle.Render(strings.Repeat("─", 35)))
	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  esc  Close"))

	content := sb.String()

	boxWidth := 48
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 2).
		Width(boxWidth).
		Render(content)

	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}
