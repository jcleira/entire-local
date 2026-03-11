package dashboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHelp(width, height int) string {
	bindings := []struct{ key, desc string }{
		{"j/k", "Navigate / scroll"},
		{"Enter/l", "Select checkpoint"},
		{"Esc/h", "Go back"},
		{"/", "Filter checkpoints"},
		{"Tab", "Next tab (detail view)"},
		{"Shift+Tab", "Previous tab"},
		{"1-3", "Jump to tab"},
		{"r", "Refresh data"},
		{"a", "Actions menu"},
		{"?", "Toggle help"},
		{"q", "Quit"},
	}

	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	for _, b := range bindings {
		sb.WriteString("  ")
		sb.WriteString(helpKeyStyle.Render(padRight(b.key, 14)))
		sb.WriteString(helpDescStyle.Render(b.desc))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  Press ? or Esc to close"))

	content := sb.String()

	boxWidth := 44
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

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
