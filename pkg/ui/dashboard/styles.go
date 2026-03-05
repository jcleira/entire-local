package dashboard

import "github.com/charmbracelet/lipgloss"

var (
	primary = lipgloss.Color("208")
	accent  = lipgloss.Color("214")
	success = lipgloss.Color("46")
	errClr  = lipgloss.Color("196")
	info    = lipgloss.Color("75")
	subtle  = lipgloss.Color("250")

	headerStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	statCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	statValueStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	statLabelStyle = lipgloss.NewStyle().
			Foreground(subtle)

	agentBadgeStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	branchStyle = lipgloss.NewStyle().
			Foreground(info).
			Bold(true)

	hashStyle = lipgloss.NewStyle().
			Foreground(subtle)

	addStyle = lipgloss.NewStyle().
			Foreground(success)

	delStyle = lipgloss.NewStyle().
			Foreground(errClr)

	toolBadgeStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Italic(true)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(subtle)

	filterStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(subtle)

	tabActiveStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(subtle)

	sectionHeaderStyle = lipgloss.NewStyle().
				Foreground(accent).
				Bold(true)

	scrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(subtle)

	selectedRowStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primary).
				PaddingLeft(1).
				PaddingRight(1)

	normalRowStyle = lipgloss.NewStyle().
			PaddingLeft(3)

	hunkHeaderStyle = lipgloss.NewStyle().
			Foreground(info)

	diffFileHeaderStyle = lipgloss.NewStyle().
				Foreground(subtle).
				Bold(true)
)
