package service

import (
	"context"

	"example.com/calibration-vault/internal/domain"
)

type WorkflowResult struct {
	Sample   domain.Sample
	Timeline []domain.Event
}

func (s *Service) IntakeApproveRelease(
	ctx context.Context,
	intake domain.IntakeInput,
	review domain.ReviewInput,
	release domain.ReleaseInput,
) (WorkflowResult, error) {
	sample, err := s.Intake(ctx, intake)
	if err != nil {
		return WorkflowResult{}, err
	}
	sample, err = s.Review(ctx, sample.SampleID, review)
	if err != nil {
		return WorkflowResult{}, err
	}
	sample, err = s.Release(ctx, sample.SampleID, release)
	if err != nil {
		return WorkflowResult{}, err
	}
	timeline, err := s.Timeline(ctx, sample.SampleID)
	if err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{Sample: sample, Timeline: timeline}, nil
}
