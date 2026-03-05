// Package commands provides CLI output helpers.
package commands

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
)

func PrintError(msg string) {
	fmt.Println(errorStyle.Render("Error: ") + msg)
}

func PrintInfo(msg string) {
	fmt.Println(infoStyle.Render(msg))
}
