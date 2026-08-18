package domain

import "sort"

type SamplePage struct {
	Items  []Sample                 `json:"items"`
	Labels map[string]DisplayLabels `json:"labels"`
	Count  int                      `json:"count"`
	Offset int                      `json:"offset"`
	Limit  int                      `json:"limit"`
}

func SortSamples(items []Sample) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
}

func PaginateSamples(items []Sample, offset, limit int) SamplePage {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	total := len(items)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	window := append([]Sample(nil), items[offset:end]...)
	labels := make(map[string]DisplayLabels, len(window))
	for _, sample := range window {
		labels[sample.SampleID] = LabelsFor(sample)
	}
	return SamplePage{Items: window, Labels: labels, Count: total, Offset: offset, Limit: limit}
}
