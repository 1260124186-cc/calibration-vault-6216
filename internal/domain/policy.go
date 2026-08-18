package domain

import (
	"fmt"
	"strings"
)

type Transition struct {
	From Status `json:"from"`
	To   Status `json:"to"`
	Why  string `json:"why"`
}

func AllowedTransition(from, to Status) bool {
	switch from {
	case StatusPendingReview:
		return to == StatusApproved || to == StatusRejected
	case StatusApproved:
		return to == StatusReleased
	default:
		return false
	}
}

func CheckTransition(from, to Status) error {
	if !AllowedTransition(from, to) {
		return fmt.Errorf("%w: cannot move from %s to %s", ErrInvalidState, from, to)
	}
	return nil
}

func ReviewTarget(decision Decision) Status {
	if decision == DecisionReject {
		return StatusRejected
	}
	return StatusApproved
}

func TransitionForReview(sample Sample, review Review) (Transition, error) {
	target := ReviewTarget(review.Decision)
	if err := CheckTransition(sample.Status, target); err != nil {
		return Transition{}, err
	}
	return Transition{From: sample.Status, To: target, Why: strings.TrimSpace(review.Note)}, nil
}

func TransitionForRelease(sample Sample, release Release) (Transition, error) {
	if err := CheckTransition(sample.Status, StatusReleased); err != nil {
		return Transition{}, err
	}
	return Transition{From: sample.Status, To: StatusReleased, Why: release.Destination}, nil
}

func StateLabel(status Status) string {
	switch status {
	case StatusPendingReview:
		return "待复核"
	case StatusApproved:
		return "已批准"
	case StatusRejected:
		return "已退回"
	case StatusReleased:
		return "已放行"
	default:
		return "未知"
	}
}

func IsTerminal(status Status) bool {
	return status == StatusRejected || status == StatusReleased
}
