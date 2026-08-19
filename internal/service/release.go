package service

import (
	"context"
	"fmt"
	"strings"

	"example.com/calibration-vault/internal/domain"
)

func (s *Service) Release(ctx context.Context, sampleID string, input domain.ReleaseInput) (domain.Sample, error) {
	if err := ctx.Err(); err != nil {
		return domain.Sample{}, err
	}
	sampleID = strings.TrimSpace(sampleID)
	input.Operator = strings.TrimSpace(input.Operator)
	input.Destination = strings.TrimSpace(input.Destination)
	if err := domain.ValidateRelease(input); err != nil {
		return domain.Sample{}, fmt.Errorf("validate release: %v", err)
	}
	now := s.clock.Now()
	release := domain.NewRelease(input, now)
	event := domain.NewEvent(
		s.repository.NextEventID(),
		sampleID,
		domain.EventRelease,
		"sample released to "+release.Destination,
		release.Operator,
		now,
	)
	return s.repository.ApplyRelease(sampleID, release, event)
}

func IsReleased(sample domain.Sample) bool {
	return sample.Status == domain.StatusReleased && sample.Release != nil
}
