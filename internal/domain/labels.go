package domain

import "strings"

type DisplayLabels struct {
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Terminal bool   `json:"terminal"`
}

func LabelsFor(sample Sample) DisplayLabels {
	return DisplayLabels{
		Status:   StateLabel(sample.Status),
		Priority: PriorityLabel(sample.Priority),
		Terminal: IsTerminal(sample.Status),
	}
}

func PriorityLabel(priority Priority) string {
	switch priority {
	case PriorityRoutine:
		return "常规"
	case PriorityPriority:
		return "优先"
	case PriorityUrgent:
		return "紧急"
	default:
		return "未设置"
	}
}

func NormalizeStatus(raw string) Status {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(StatusPendingReview):
		return StatusPendingReview
	case string(StatusApproved):
		return StatusApproved
	case string(StatusRejected):
		return StatusRejected
	case string(StatusReleased):
		return StatusReleased
	default:
		return ""
	}
}

func IsActionable(sample Sample) bool {
	return sample.Status == StatusPendingReview || sample.Status == StatusApproved
}
