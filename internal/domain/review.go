package domain

import "time"

type Decision string

const (
	DecisionApprove Decision = "approve"
	DecisionReject  Decision = "reject"
)

type Review struct {
	Reviewer   string    `json:"reviewer"`
	Decision   Decision  `json:"decision"`
	Note       string    `json:"note"`
	ReviewedAt time.Time `json:"reviewed_at"`
}

type ReviewInput struct {
	Reviewer string   `json:"reviewer"`
	Decision Decision `json:"decision"`
	Note     string   `json:"note"`
}

func NewReview(input ReviewInput, now time.Time) Review {
	return Review{
		Reviewer:   input.Reviewer,
		Decision:   input.Decision,
		Note:       input.Note,
		ReviewedAt: now,
	}
}

func (r Review) Accepted() bool {
	return r.Decision == DecisionApprove
}

func (r Review) Summary() string {
	if r.Accepted() {
		return "sample approved"
	}
	return "sample rejected"
}

func (r Review) Clone() Review {
	return Review{
		Reviewer:   r.Reviewer,
		Decision:   r.Decision,
		Note:       r.Note,
		ReviewedAt: r.ReviewedAt,
	}
}
