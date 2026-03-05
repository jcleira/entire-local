package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/jcleira/entire-local/pkg/checkpoint"
)

const (
	roleUser      = "user"
	roleAssistant = "assistant"
)

var (
	userLabelStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	assistantLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Bold(true)

	userBorderStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(primary).
			PaddingLeft(1)

	assistantBorderStyle = lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(subtle).
				PaddingLeft(1)

	timestampStyle = lipgloss.NewStyle().
			Foreground(subtle)

	toolResultStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Italic(true)
)

type messageGroup struct {
	role    string
	content string
	tools   []string
	results int
	ts      string
}

func groupTranscriptEntries(entries []checkpoint.TranscriptEntry) []messageGroup {
	var groups []messageGroup
	var pendingTools []string
	pendingResults := 0

	flushTools := func() {
		if len(pendingTools) > 0 || pendingResults > 0 {
			groups = append(groups, messageGroup{
				role:    "tools",
				tools:   pendingTools,
				results: pendingResults,
			})
			pendingTools = nil
			pendingResults = 0
		}
	}

	for _, e := range entries {
		switch {
		case e.Role == "human" || e.Role == roleUser:
			if e.Content == "" {
				continue
			}
			flushTools()
			ts := ""
			if !e.Timestamp.IsZero() {
				ts = e.Timestamp.Format("3:04 PM")
			}
			groups = append(groups, messageGroup{
				role:    roleUser,
				content: e.Content,
				ts:      ts,
			})

		case e.Role == roleAssistant && e.ToolName != "":
			pendingTools = append(pendingTools, e.ToolName)

		case e.Role == roleAssistant:
			if e.Content == "" {
				continue
			}
			flushTools()
			ts := ""
			if !e.Timestamp.IsZero() {
				ts = e.Timestamp.Format("3:04 PM")
			}
			groups = append(groups, messageGroup{
				role:    roleAssistant,
				content: e.Content,
				ts:      ts,
			})

		case e.Type == "tool_result":
			pendingResults++

		default:
			if e.Content != "" {
				flushTools()
				groups = append(groups, messageGroup{
					role:    "other",
					content: e.Content,
				})
			}
		}
	}
	flushTools()
	return groups
}

func renderTabTranscript(entries []checkpoint.TranscriptEntry, width, height, scrollOffset int) string {
	if len(entries) == 0 {
		return dimStyle.Render("  Loading transcript...")
	}

	contentWidth := width - 8
	if contentWidth < 30 {
		contentWidth = 30
	}

	var allLines []string

	messageCount := 0
	for _, e := range entries {
		if (e.Role == "human" || e.Role == roleUser || e.Role == roleAssistant) && e.Content != "" && e.ToolName == "" {
			messageCount++
		}
	}

	allLines = append(allLines,
		sectionHeaderStyle.Render(fmt.Sprintf("  Transcript (%d messages)", messageCount)),
		"",
	)

	groups := groupTranscriptEntries(entries)

	for _, g := range groups {
		switch g.role {
		case roleUser:
			label := userLabelStyle.Render("You")
			if g.ts != "" {
				pad := contentWidth - lipgloss.Width(label) - len(g.ts)
				if pad < 2 {
					pad = 2
				}
				allLines = append(allLines, "  "+label+strings.Repeat(" ", pad)+timestampStyle.Render(g.ts))
			} else {
				allLines = append(allLines, "  "+label)
			}

			rendered := renderMarkdown(g.content, contentWidth-4)
			block := userBorderStyle.Width(contentWidth).Render(rendered)
			for _, line := range strings.Split(block, "\n") {
				allLines = append(allLines, "  "+line)
			}
			allLines = append(allLines, "")

		case roleAssistant:
			label := assistantLabelStyle.Render("Assistant")
			if g.ts != "" {
				pad := contentWidth - lipgloss.Width(label) - len(g.ts)
				if pad < 2 {
					pad = 2
				}
				allLines = append(allLines, "    "+label+strings.Repeat(" ", pad)+timestampStyle.Render(g.ts))
			} else {
				allLines = append(allLines, "    "+label)
			}

			rendered := renderMarkdown(g.content, contentWidth-6)
			block := assistantBorderStyle.Width(contentWidth).Render(rendered)
			for _, line := range strings.Split(block, "\n") {
				allLines = append(allLines, "    "+line)
			}
			allLines = append(allLines, "")

		case "tools":
			var badges []string
			for _, t := range g.tools {
				badges = append(badges, toolBadgeStyle.Render("["+t+"]"))
			}
			if len(badges) > 0 {
				allLines = append(allLines, "    "+strings.Join(badges, " "))
			}
			if g.results > 0 {
				allLines = append(allLines, "    "+toolResultStyle.Render(fmt.Sprintf("[%d results]", g.results)))
			}
			allLines = append(allLines, "")

		case "other":
			content := g.content
			if len(content) > contentWidth {
				content = content[:contentWidth-3] + "..."
			}
			allLines = append(allLines, "  "+content, "")
		}
	}

	maxScroll := len(allLines) - height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scrollOffset > maxScroll {
		scrollOffset = maxScroll
	}
	if scrollOffset < 0 {
		scrollOffset = 0
	}

	end := scrollOffset + height
	if end > len(allLines) {
		end = len(allLines)
	}

	visible := allLines[scrollOffset:end]

	if maxScroll > 0 {
		pos := scrollIndicatorStyle.Render(fmt.Sprintf("(%d/%d)", scrollOffset+1, maxScroll+1))
		visible = append(visible, "  "+pos)
	}

	return strings.Join(visible, "\n")
}
