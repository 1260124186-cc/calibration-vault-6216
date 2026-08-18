package service

import (
	"context"
	"strings"

	"example.com/calibration-vault/internal/domain"
)

func (s *Service) Review(ctx context.Context, sampleID string, input domain.ReviewInput) (domain.Sample, error) {
	sampleID = strings.TrimSpace(sampleID)
	input.Reviewer = strings.TrimSpace(input.Reviewer)
	input.Note = strings.TrimSpace(input.Note)
	if err := domain.ValidateReview(input); err != nil {
		return domain.Sample{}, err
	}
	now := s.clock.Now()
	review := domain.NewReview(input, now)
	eventKind := domain.EventReview
	if !review.Accepted() {
		eventKind = domain.EventRejection
	}
	event := domain.NewEvent(
		s.repository.NextEventID(),
		sampleID,
		eventKind,
		review.Summary(),
		review.Reviewer,
		now,
	)
	return s.repository.ApplyReview(sampleID, review, event)
}

func IsApproval(sample domain.Sample) bool {
	return sample.Status == domain.StatusApproved
}
