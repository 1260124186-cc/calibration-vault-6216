package service

import (
	"context"
	"strings"

	"example.com/calibration-vault/internal/domain"
)

func (s *Service) Intake(ctx context.Context, input domain.IntakeInput) (domain.Sample, error) {
	if err := ctx.Err(); err != nil {
		return domain.Sample{}, err
	}
	input.SampleID = strings.TrimSpace(input.SampleID)
	input.Source = strings.TrimSpace(input.Source)
	input.Tags = domain.NormalizeTags(input.Tags)
	if err := domain.ValidateIntake(input); err != nil {
		return domain.Sample{}, err
	}
	now := s.clock.Now()
	sample := domain.NewSample(input.SampleID, input, now)
	event := domain.NewEvent(
		s.repository.NextEventID(),
		sample.SampleID,
		domain.EventIntake,
		"sample received for review",
		input.Source,
		now,
	)
	if err := s.repository.Create(sample, event); err != nil {
		return domain.Sample{}, err
	}
	return sample, nil
}
