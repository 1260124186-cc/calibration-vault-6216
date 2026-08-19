package domain

import (
	"sort"
	"time"
)

type Metrics struct {
	Summary          Summary        `json:"summary"`
	TagCounts        map[string]int `json:"tag_counts"`
	OldestOpenSample string         `json:"oldest_open_sample,omitempty"`
	OldestOpenAge    time.Duration  `json:"oldest_open_age"`
	GeneratedAt      time.Time      `json:"generated_at"`
}

func BuildMetrics(samples []Sample, now time.Time) Metrics {
	result := Metrics{
		Summary:     BuildSummary(samples),
		TagCounts:   TagCounts(samples),
		GeneratedAt: now,
	}
	open := make([]Sample, 0)
	for _, sample := range samples {
		if sample.IsOpen() {
			open = append(open, sample)
		}
	}
	sort.SliceStable(open, func(i, j int) bool {
		return open[i].CreatedAt.Before(open[j].CreatedAt)
	})
	if len(open) > 0 {
		result.OldestOpenSample = open[0].SampleID
		result.OldestOpenAge = now.Sub(open[0].CreatedAt)
		if result.OldestOpenAge < 0 {
			result.OldestOpenAge = 0
		}
	}
	return result
}

func PriorityRatio(summary Summary, priority Priority) float64 {
	if summary.Total == 0 {
		return 0
	}
	return float64(summary.ByPriority[priority]) / float64(summary.Total)
}

func StatusRatio(summary Summary, status Status) float64 {
	if summary.Total == 0 {
		return 0
	}
	return float64(summary.ByStatus[status]) / float64(summary.Total)
}

func OrderedTagCounts(counts map[string]int) []string {
	tags := make([]string, 0, len(counts))
	for tag := range counts {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
