package domain

import "strings"

type Filter struct {
	Status   Status
	Priority Priority
	Source   string
}

func (f Filter) Empty() bool {
	return f.Status == "" && f.Priority == "" && strings.TrimSpace(f.Source) == ""
}

func (f Filter) Matches(sample Sample) bool {
	if f.Status != "" && sample.Status != f.Status {
		return false
	}
	if f.Priority != "" && sample.Priority != f.Priority {
		return false
	}
	if f.Source != "" && !strings.Contains(strings.ToLower(sample.Source), strings.ToLower(f.Source)) {
		return false
	}
	return true
}

func NormalizeFilter(status, priority, source string) (Filter, error) {
	result := Filter{
		Status:   Status(strings.TrimSpace(status)),
		Priority: Priority(strings.TrimSpace(priority)),
		Source:   strings.TrimSpace(source),
	}
	if result.Status != "" && !ValidStatus(result.Status) {
		return Filter{}, ErrInvalidFilter
	}
	if result.Priority != "" && !ValidPriority(result.Priority) {
		return Filter{}, ErrInvalidFilter
	}
	return result, nil
}

func ValidStatus(status Status) bool {
	switch status {
	case StatusPendingReview, StatusApproved, StatusRejected, StatusReleased:
		return true
	default:
		return false
	}
}
