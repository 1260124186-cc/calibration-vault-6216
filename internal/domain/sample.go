package domain

import "time"

type Status string

const (
	StatusPendingReview Status = "pending_review"
	StatusApproved      Status = "approved"
	StatusRejected      Status = "rejected"
	StatusReleased      Status = "released"
)

type Priority string

const (
	PriorityRoutine  Priority = "routine"
	PriorityPriority Priority = "priority"
	PriorityUrgent   Priority = "urgent"
)

type Sample struct {
	ID        string    `json:"id"`
	SampleID  string    `json:"sample_id"`
	Source    string    `json:"source"`
	Priority  Priority  `json:"priority"`
	Tags      []string  `json:"tags"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Review    *Review   `json:"review,omitempty"`
	Release   *Release  `json:"release,omitempty"`
}

type IntakeInput struct {
	SampleID string   `json:"sample_id"`
	Source   string   `json:"source"`
	Priority Priority `json:"priority"`
	Tags     []string `json:"tags"`
}

func NewSample(id string, input IntakeInput, now time.Time) Sample {
	tags := input.Tags
	return Sample{
		ID:        id,
		SampleID:  input.SampleID,
		Source:    input.Source,
		Priority:  input.Priority,
		Tags:      tags,
		Status:    StatusPendingReview,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s Sample) Clone() Sample {
	copy := s
	copy.Tags = s.Tags
	if s.Review != nil {
		review := *s.Review
		copy.Review = &review
	}
	if s.Release != nil {
		release := *s.Release
		copy.Release = &release
	}
	return copy
}

func (s Sample) IsOpen() bool {
	return s.Status == StatusPendingReview || s.Status == StatusApproved
}

func (s Sample) CanReview() bool {
	return s.Status == StatusPendingReview
}

func (s Sample) CanRelease() bool {
	return s.Status == StatusApproved
}
