package dashboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/jcleira/entire-local/pkg/checkpoint"
)

func renderCheckpointList(summaries []checkpoint.CheckpointSummary, cursor, height int, filter string, width int) string {
	filtered := filterSummaries(summaries, filter)

	var sb strings.Builder

	header := headerStyle.Render("Checkpoints")
	if filter != "" {
		header += " " + filterStyle.Render("[/"+filter+"]")
	}
	header += dimStyle.Render(fmt.Sprintf(" (%d)", len(filtered)))
	sb.WriteString(header)
	sb.WriteString("\n")

	if len(filtered) == 0 {
		sb.WriteString(dimStyle.Render("  No checkpoints found"))
		return sb.String()
	}

	visibleHeight := (height - 2) / 4
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	start := 0
	if cursor >= visibleHeight {
		start = cursor - visibleHeight + 1
	}
	end := start + visibleHeight
	if end > len(filtered) {
		end = len(filtered)
	}

	contentWidth := width - 4
	if contentWidth < 40 {
		contentWidth = 40
	}

	for i := start; i < end; i++ {
		s := filtered[i]
		block := renderCheckpointBlock(s, contentWidth)

		if i == cursor {
			sb.WriteString(selectedRowStyle.Width(contentWidth - 2).Render(block))
		} else {
			sb.WriteString(normalRowStyle.Width(contentWidth).Render(block))
		}
		sb.WriteString("\n")
	}

	if len(filtered) > visibleHeight {
		sb.WriteString(scrollIndicatorStyle.Render(fmt.Sprintf("  (%d/%d)", cursor+1, len(filtered))))
	}

	return sb.String()
}

func renderCheckpointBlock(s checkpoint.CheckpointSummary, width int) string {
	sep := dimStyle.Render("  ·  ")

	title := s.CommitMessage
	if title == "" {
		title = cleanPromptExcerpt(s.Prompt)
	}
	if title == "" {
		title = s.Context
	}
	if title == "" {
		title = s.Branch
	}
	if idx := strings.IndexByte(title, '\n'); idx >= 0 {
		title = title[:idx]
	}
	title = strings.TrimSpace(title)
	if len(title) > width {
		title = title[:width-1] + "…"
	}
	line1 := lipgloss.NewStyle().Bold(true).Render(title)

	subtitle := ""
	if s.CommitMessage != "" {
		excerpt := cleanPromptExcerpt(s.Prompt)
		if excerpt != "" && excerpt != title {
			if len(excerpt) > width {
				excerpt = excerpt[:width-1] + "…"
			}
			subtitle = dimStyle.Render(excerpt)
		}
	}

	var leftParts []string
	leftParts = append(leftParts, hashStyle.Render(truncateHash(s.ID)))
	timeStr := formatTimeAgo(s.CreatedAt)
	if timeStr != "" {
		leftParts = append(leftParts, dimStyle.Render(timeStr))
	}
	attribution := formatAttribution(s)
	if attribution != "" {
		leftParts = append(leftParts, agentBadgeStyle.Render(attribution))
	}
	leftSide := strings.Join(leftParts, sep)

	colDiff := lipgloss.NewStyle().Width(14).Align(lipgloss.Right)
	colFiles := lipgloss.NewStyle().Width(10).Align(lipgloss.Right)
	colSessions := lipgloss.NewStyle().Width(12).Align(lipgloss.Right)
	colTokens := lipgloss.NewStyle().Width(14).Align(lipgloss.Right)

	var rightCols []string
	if s.Additions > 0 || s.Deletions > 0 {
		diffStr := addStyle.Render(fmt.Sprintf("+%d", s.Additions)) + " / " +
			delStyle.Render(fmt.Sprintf("-%d", s.Deletions))
		rightCols = append(rightCols, colDiff.Render(diffStr))
	} else {
		rightCols = append(rightCols, colDiff.Render(""))
	}
	if s.FileCount > 0 {
		rightCols = append(rightCols, colFiles.Render(dimStyle.Render(fmt.Sprintf("%d files", s.FileCount))))
	} else {
		rightCols = append(rightCols, colFiles.Render(""))
	}
	if s.SessionCount > 0 {
		label := "session"
		if s.SessionCount > 1 {
			label = "sessions"
		}
		rightCols = append(rightCols, colSessions.Render(dimStyle.Render(fmt.Sprintf("%d %s", s.SessionCount, label))))
	} else {
		rightCols = append(rightCols, colSessions.Render(""))
	}
	if s.TotalTokens > 0 {
		rightCols = append(rightCols, colTokens.Render(dimStyle.Render(formatTokensShort(s.TotalTokens))))
	} else {
		rightCols = append(rightCols, colTokens.Render(""))
	}
	rightSide := strings.Join(rightCols, dimStyle.Render(" · "))

	metaLine := leftSide
	if rightSide != "" {
		padding := width - lipgloss.Width(leftSide) - lipgloss.Width(rightSide)
		if padding < 2 {
			padding = 2
		}
		metaLine = leftSide + strings.Repeat(" ", padding) + rightSide
	}

	if subtitle != "" {
		return line1 + "\n" + subtitle + "\n" + metaLine
	}
	return line1 + "\n" + metaLine
}

func cleanPromptExcerpt(prompt string) string {
	if prompt == "" {
		return ""
	}

	prefixes := []string{
		"Implement the following plan:",
		"implement the following plan:",
	}
	cleaned := strings.TrimSpace(prompt)
	for _, p := range prefixes {
		cleaned = strings.TrimPrefix(cleaned, p)
	}
	cleaned = strings.TrimSpace(cleaned)

	for _, line := range strings.SplitN(cleaned, "\n", 20) {
		line = strings.TrimSpace(line)
		if line == "" || line == "#" || line == "---" {
			continue
		}
		line = strings.TrimLeft(line, "# ")
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return ""
}

func filterSummaries(summaries []checkpoint.CheckpointSummary, filter string) []checkpoint.CheckpointSummary {
	if filter == "" {
		return summaries
	}

	f := strings.ToLower(filter)
	var result []checkpoint.CheckpointSummary
	for i := range summaries {
		if strings.Contains(strings.ToLower(summaries[i].Branch), f) ||
			strings.Contains(strings.ToLower(summaries[i].Agent), f) ||
			strings.Contains(strings.ToLower(summaries[i].Prompt), f) ||
			strings.Contains(strings.ToLower(summaries[i].Context), f) ||
			strings.Contains(strings.ToLower(summaries[i].ID), f) ||
			strings.Contains(strings.ToLower(summaries[i].CommitMessage), f) {
			result = append(result, summaries[i])
		}
	}
	return result
}

func formatTimeAgo(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}

func formatNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%d", n)
}

func formatTokensShort(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk tokens", float64(n)/1000)
	}
	return fmt.Sprintf("%d tokens", n)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

func formatAttribution(s checkpoint.CheckpointSummary) string {
	agent := truncate(s.Agent, 16)
	pct := ""
	if s.AgentPercent > 0 {
		pct = fmt.Sprintf(" %d%%", s.AgentPercent)
	}
	if s.Author != "" && agent != "" {
		return s.Author + " (" + agent + pct + ")"
	}
	if s.Author != "" {
		return s.Author
	}
	if agent != "" {
		return agent + pct
	}
	return ""
}

func truncateHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
