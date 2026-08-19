package service

import (
	"context"
	"sort"

	"example.com/calibration-vault/internal/domain"
)

type OperationsReport struct {
	Metrics       domain.Metrics  `json:"metrics"`
	UrgentSamples []domain.Sample `json:"urgent_samples"`
	PendingReview []domain.Sample `json:"pending_review"`
}

func (s *Service) OperationsReport(ctx context.Context) (OperationsReport, error) {
	if err := ctx.Err(); err != nil {
		return OperationsReport{}, err
	}
	all := s.repository.List(domain.Filter{})
	report := OperationsReport{
		Metrics: domain.BuildMetrics(all, s.clock.Now()),
	}
	for _, sample := range all {
		switch {
		case sample.Priority == domain.PriorityUrgent && sample.IsOpen():
			report.UrgentSamples = append(report.UrgentSamples, sample)
		case sample.Status == domain.StatusPendingReview:
			report.PendingReview = append(report.PendingReview, sample)
		}
	}
	sort.SliceStable(report.UrgentSamples, func(i, j int) bool {
		return report.UrgentSamples[i].CreatedAt.Before(report.UrgentSamples[j].CreatedAt)
	})
	sort.SliceStable(report.PendingReview, func(i, j int) bool {
		return report.PendingReview[i].CreatedAt.Before(report.PendingReview[j].CreatedAt)
	})
	return report, nil
}

func (s *Service) SearchTag(ctx context.Context, tag string) ([]domain.Sample, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	items := s.repository.List(domain.Filter{})
	result := make([]domain.Sample, 0)
	for _, sample := range items {
		if domain.HasTag(sample, tag) {
			result = append(result, sample)
		}
	}
	return result, nil
}

func (s *Service) Metadata(ctx context.Context) (map[string]any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return map[string]any{
		"priorities":   []domain.Priority{domain.PriorityRoutine, domain.PriorityPriority, domain.PriorityUrgent},
		"statuses":     []domain.Status{domain.StatusPendingReview, domain.StatusApproved, domain.StatusRejected, domain.StatusReleased},
		"destinations": domain.SupportedDestinations(),
	}, nil
}
