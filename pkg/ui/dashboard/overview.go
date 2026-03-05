package dashboard

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/jcleira/entire-local/pkg/checkpoint"
)

func renderOverview(stats checkpoint.OverviewStats, width int) string {
	cardWidth := (width - 10) / 4
	if cardWidth < 18 {
		cardWidth = 18
	}
	if cardWidth > 40 {
		cardWidth = 40
	}

	cards := []string{
		renderStatCardDynamic("Checkpoints", fmt.Sprintf("%d", stats.TotalCheckpoints), "total", cardWidth),
		renderStatCardDynamic("Avg Tokens", fmt.Sprintf("%d", stats.AvgTokens), "per checkpoint", cardWidth),
		renderStatCardDynamic("Peak Duration", stats.PeakDuration, "longest session", cardWidth),
		renderStatCardDynamic("Streak", fmt.Sprintf("%d days", stats.StreakDays), "consecutive", cardWidth),
	}

	cardRow := lipgloss.JoinHorizontal(lipgloss.Top, cards...)

	chart := renderActivityChart(stats.ActivityByDay, width)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("Overview"),
		"",
		cardRow,
		"",
		chart,
	)
}

func renderStatCardDynamic(label, value, sub string, cardWidth int) string {
	content := lipgloss.JoinVertical(lipgloss.Left,
		statValueStyle.Render(value),
		statLabelStyle.Render(label),
		dimStyle.Render(sub),
	)
	return statCardStyle.Width(cardWidth).Render(content)
}

func renderActivityChart(activity map[string]int, width int) string {
	if len(activity) == 0 {
		return dimStyle.Render("  No activity data")
	}

	days := last30Days()
	maxCount := 0
	for _, d := range days {
		if c := activity[d]; c > maxCount {
			maxCount = c
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	blocks := []rune{'░', '▒', '▓', '█'}
	barWidth := width - 16
	if barWidth < 10 {
		barWidth = 10
	}

	var sb strings.Builder
	sb.WriteString(dimStyle.Render("  Activity (last 30 days)"))
	sb.WriteString("\n  ")

	shown := barWidth
	if shown > len(days) {
		shown = len(days)
	}
	start := len(days) - shown

	for i := start; i < len(days); i++ {
		count := activity[days[i]]
		if count == 0 {
			sb.WriteString(dimStyle.Render("·"))
		} else {
			idx := (count * len(blocks) / (maxCount + 1))
			if idx >= len(blocks) {
				idx = len(blocks) - 1
			}
			sb.WriteString(lipgloss.NewStyle().Foreground(primary).Render(string(blocks[idx])))
		}
	}

	return sb.String()
}

func last30Days() []string {
	days := make([]string, 30)
	now := time.Now()
	for i := range 30 {
		days[i] = now.AddDate(0, 0, -(29 - i)).Format("2006-01-02")
	}
	sort.Strings(days)
	return days
}
