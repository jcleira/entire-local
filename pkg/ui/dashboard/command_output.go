package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) handleCommandOutputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	switch {
	case msg.Type == tea.KeyEscape, msg.String() == "q":
		m.showCommandOutput = false
		m.commandOutput = ""
		m.commandOutputScroll = 0
		return m, nil

	case key.Matches(msg, keys.Down):
		m.commandOutputScroll++
		return m, nil

	case key.Matches(msg, keys.Up):
		if m.commandOutputScroll > 0 {
			m.commandOutputScroll--
		}
		return m, nil
	}

	return m, nil
}

func renderCommandOutput(output string, scroll, width, height int) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Command Output"))
	sb.WriteString("\n")
	sb.WriteString(actionSepStyle.Render(strings.Repeat("─", 40)))
	sb.WriteString("\n\n")

	lines := strings.Split(output, "\n")

	boxWidth := width - 10
	if boxWidth > 80 {
		boxWidth = 80
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	contentHeight := height - 14
	if contentHeight < 5 {
		contentHeight = 5
	}

	maxScroll := len(lines) - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}

	end := scroll + contentHeight
	if end > len(lines) {
		end = len(lines)
	}

	visible := lines[scroll:end]
	for _, line := range visible {
		sb.WriteString("  " + line + "\n")
	}

	if maxScroll > 0 {
		sb.WriteString("\n")
		sb.WriteString(scrollIndicatorStyle.Render(fmt.Sprintf("  (%d/%d)", scroll+1, maxScroll+1)))
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  j/k scroll  esc/q close"))

	content := sb.String()

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
