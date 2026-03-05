package checkpoint

import (
	"sort"
	"time"
)

type OverviewStats struct {
	TotalCheckpoints int
	AvgTokens        int
	AvgSteps         int
	PeakDuration     string
	StreakDays       int
	ActivityByDay    map[string]int
	TotalTokens      int
	TotalAdditions   int
	TotalDeletions   int
}

func ComputeStats(summaries []CheckpointSummary) OverviewStats {
	stats := OverviewStats{
		TotalCheckpoints: len(summaries),
		ActivityByDay:    make(map[string]int),
	}

	if len(summaries) == 0 {
		return stats
	}

	var totalTokens int
	var totalAdds, totalDels int
	var maxDuration string
	for i := range summaries {
		totalTokens += summaries[i].TotalTokens
		totalAdds += summaries[i].Additions
		totalDels += summaries[i].Deletions
		day := summaries[i].CreatedAt.Format("2006-01-02")
		if day != "0001-01-01" {
			stats.ActivityByDay[day]++
		}
		if summaries[i].Duration > maxDuration {
			maxDuration = summaries[i].Duration
		}
	}

	stats.TotalTokens = totalTokens
	stats.TotalAdditions = totalAdds
	stats.TotalDeletions = totalDels
	if len(summaries) > 0 {
		stats.AvgTokens = totalTokens / len(summaries)
	}
	stats.PeakDuration = maxDuration
	stats.StreakDays = computeStreak(stats.ActivityByDay)

	return stats
}

func computeStreak(activity map[string]int) int {
	if len(activity) == 0 {
		return 0
	}

	days := make([]string, 0, len(activity))
	for d := range activity {
		days = append(days, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(days)))

	streak := 1
	for i := 1; i < len(days); i++ {
		prev, prevErr := time.Parse("2006-01-02", days[i-1])
		curr, currErr := time.Parse("2006-01-02", days[i])
		if prevErr != nil || currErr != nil {
			break
		}
		if prev.Sub(curr) == 24*time.Hour {
			streak++
		} else {
			break
		}
	}

	return streak
}
