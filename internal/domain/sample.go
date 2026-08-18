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
	// 拷贝 Tags，避免复用调用方传入的 slice 导致后续被外部修改污染样本
	tags := append([]string(nil), input.Tags...)
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
	// Tags 是 slice，必须深拷贝底层数组，否则调用方修改返回值会污染仓储里的样本
	copy.Tags = append([]string(nil), s.Tags...)
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
