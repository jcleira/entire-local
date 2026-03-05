package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/jcleira/entire-local/pkg/checkpoint"
)

func renderTabBar(active detailTab) string {
	tabs := []struct {
		label string
		tab   detailTab
	}{
		{"Transcript", tabTranscript},
		{"Files", tabFiles},
		{"Plan", tabPlan},
	}

	var parts []string
	for _, t := range tabs {
		if t.tab == active {
			parts = append(parts, tabActiveStyle.Render(" ["+t.label+"] "))
		} else {
			parts = append(parts, tabInactiveStyle.Render("  "+t.label+"  "))
		}
	}

	return "  " + strings.Join(parts, "")
}

func renderDetailHeader(cp *checkpoint.Checkpoint, _ int) string {
	parts := make([]string, 0, 2)

	parts = append(parts, headerStyle.Render("Checkpoint Detail"))

	agentName := cp.Session.Agent
	if agentName == "" {
		agentName = "unknown"
	}
	agentPct := int(cp.Session.Attribution.AgentPercentage)

	row1 := fmt.Sprintf(
		"%s  %s  %s",
		hashStyle.Render("checkpoint "+truncateHash(cp.CheckpointID)),
		branchStyle.Render(cp.Branch),
		agentBadgeStyle.Render(fmt.Sprintf("%s %d%%", agentName, agentPct)),
	)

	var metaParts []string
	totalTokens := cp.TokenUsage.OutputTokens
	if totalTokens > 0 {
		metaParts = append(metaParts, dimStyle.Render(fmt.Sprintf("tokens: %s", formatNumber(totalTokens))))
	}

	if len(cp.FilesTouched) > 0 {
		metaParts = append(metaParts, dimStyle.Render(fmt.Sprintf("%d files", len(cp.FilesTouched))))
	}

	diffStats := parseDiffStats(cp.Diff)
	if diffStats.totalAdd > 0 || diffStats.totalDel > 0 {
		metaParts = append(metaParts,
			addStyle.Render(fmt.Sprintf("+%d", diffStats.totalAdd))+" "+
				delStyle.Render(fmt.Sprintf("-%d", diffStats.totalDel)))
	}

	if len(metaParts) > 0 {
		row1 += "  " + strings.Join(metaParts, "  ")
	}

	parts = append(parts, "  "+row1)

	return strings.Join(parts, "\n")
}

type fileDiffStat struct {
	name string
	add  int
	del  int
}

type diffStatsResult struct {
	files    []fileDiffStat
	totalAdd int
	totalDel int
}

func parseDiffStats(diff string) diffStatsResult {
	var result diffStatsResult
	if diff == "" {
		return result
	}

	var current *fileDiffStat
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "diff --git") {
			if current != nil {
				result.files = append(result.files, *current)
			}
			parts := strings.Fields(line)
			name := ""
			if len(parts) >= 4 {
				name = strings.TrimPrefix(parts[3], "b/")
			}
			current = &fileDiffStat{name: name}
		} else if current != nil {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				current.add++
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				current.del++
			}
		}
	}
	if current != nil {
		result.files = append(result.files, *current)
	}

	for _, f := range result.files {
		result.totalAdd += f.add
		result.totalDel += f.del
	}

	return result
}

func renderTabFiles(diff string, width, _ int) string {
	stats := parseDiffStats(diff)

	if len(stats.files) == 0 {
		return dimStyle.Render("  No files changed")
	}

	var sb strings.Builder

	headerLine := fmt.Sprintf("  Files (%d)", len(stats.files))
	totalStat := addStyle.Render(fmt.Sprintf("+%d", stats.totalAdd)) + " " +
		delStyle.Render(fmt.Sprintf("-%d", stats.totalDel)) + dimStyle.Render(" total")
	padding := width - lipgloss.Width(headerLine) - lipgloss.Width(totalStat) - 4
	if padding < 2 {
		padding = 2
	}
	sb.WriteString(sectionHeaderStyle.Render(headerLine) + strings.Repeat(" ", padding) + totalStat)
	sb.WriteString("\n")

	nameWidth := width - 24
	if nameWidth < 20 {
		nameWidth = 20
	}

	for _, f := range stats.files {
		name := f.name
		if len(name) > nameWidth {
			name = "…" + name[len(name)-nameWidth+1:]
		}

		fileStat := addStyle.Render(fmt.Sprintf("+%d", f.add)) + " " +
			delStyle.Render(fmt.Sprintf("-%d", f.del))
		pad := width - len(name) - lipgloss.Width(fileStat) - 8
		if pad < 2 {
			pad = 2
		}
		sb.WriteString("    " + name + strings.Repeat(" ", pad) + fileStat + "\n")
	}

	sb.WriteString("\n")

	for _, line := range strings.Split(diff, "\n") {
		var styled string
		switch {
		case strings.HasPrefix(line, "diff --git"),
			strings.HasPrefix(line, "---"),
			strings.HasPrefix(line, "+++"),
			strings.HasPrefix(line, "index "):
			styled = diffFileHeaderStyle.Render(line)
		case strings.HasPrefix(line, "@@"):
			styled = hunkHeaderStyle.Render(line)
		case strings.HasPrefix(line, "+"):
			styled = addStyle.Render(line)
		case strings.HasPrefix(line, "-"):
			styled = delStyle.Render(line)
		default:
			styled = line
		}
		sb.WriteString("  " + styled + "\n")
	}

	return sb.String()
}

func renderTabPlan(plan string, width, _ int) string {
	if plan == "" {
		return dimStyle.Render("  No plan available")
	}

	contentWidth := width - 6
	if contentWidth < 30 {
		contentWidth = 30
	}

	var sb strings.Builder
	sb.WriteString(sectionHeaderStyle.Render("  Plan"))
	sb.WriteString("\n\n")

	rendered := renderMarkdown(plan, contentWidth)
	for _, line := range strings.Split(rendered, "\n") {
		sb.WriteString("  " + line + "\n")
	}

	return sb.String()
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}
	var result strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if len(line) <= width {
			result.WriteString(line + "\n")
			continue
		}
		remaining := line
		for len(remaining) > width {
			breakAt := width
			if idx := strings.LastIndex(remaining[:width], " "); idx > 0 {
				breakAt = idx
			}
			result.WriteString(remaining[:breakAt] + "\n")
			remaining = strings.TrimLeft(remaining[breakAt:], " ")
		}
		if remaining != "" {
			result.WriteString(remaining + "\n")
		}
	}
	return strings.TrimRight(result.String(), "\n")
}
