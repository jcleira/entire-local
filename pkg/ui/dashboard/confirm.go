package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jcleira/entire-local/pkg/checkpoint"
)

func (m model) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	switch {
	case key.Matches(msg, keys.Confirm):
		action := *m.confirmAction
		m.confirmAction = nil
		return m.runAction(action)

	case key.Matches(msg, keys.Deny), msg.Type == tea.KeyEscape:
		m.confirmAction = nil
		return m, nil
	}

	return m, nil
}

func renderConfirmDialog(action actionItem, detail *checkpoint.Checkpoint, width, height int) string {
	var sb strings.Builder
	sb.WriteString(warningStyle.Render("Confirm Action"))
	sb.WriteString("\n\n")

	msg := confirmMessage(action, detail)
	for _, line := range strings.Split(msg, "\n") {
		sb.WriteString("  " + line + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString("  ")
	sb.WriteString(helpKeyStyle.Render("y"))
	sb.WriteString(" confirm    ")
	sb.WriteString(helpKeyStyle.Render("n"))
	sb.WriteString("/")
	sb.WriteString(helpKeyStyle.Render("esc"))
	sb.WriteString(" cancel")

	content := sb.String()

	boxWidth := 52
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(1, 2).
		Width(boxWidth).
		Render(content)

	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

func confirmMessage(action actionItem, detail *checkpoint.Checkpoint) string {
	switch action.action {
	case actionRewind:
		id := "current"
		if detail != nil {
			id = truncateHash(detail.CheckpointID)
		}
		return fmt.Sprintf("Rewind to checkpoint %s?\nThis will reset your working tree.", id)
	case actionResume:
		branch := ""
		if detail != nil {
			branch = detail.Branch
		}
		return fmt.Sprintf("Resume session on %s?\nThis will switch branches.", branch)
	case actionReset:
		return "Reset session state?\nThis will delete current session data."
	case actionClean:
		return "Clean orphaned data?\nThis will remove unreferenced checkpoints."
	case actionStatus, actionExplain, actionDoctor:
		return "Proceed with this action?"
	}
	return "Proceed with this action?"
}
