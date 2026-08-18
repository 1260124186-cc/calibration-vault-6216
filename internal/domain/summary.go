package domain

type Summary struct {
	Total         int              `json:"total"`
	ByStatus      map[Status]int   `json:"by_status"`
	ByPriority    map[Priority]int `json:"by_priority"`
	OpenCount     int              `json:"open_count"`
	ReleasedCount int              `json:"released_count"`
}

func BuildSummary(samples []Sample) Summary {
	result := Summary{
		ByStatus:   map[Status]int{},
		ByPriority: map[Priority]int{},
	}
	for _, sample := range samples {
		result.Total++
		result.ByStatus[sample.Status]++
		result.ByPriority[sample.Priority]++
		if sample.IsOpen() {
			result.OpenCount++
		}
		if sample.Status == StatusReleased {
			result.ReleasedCount++
		}
	}
	return result
}
