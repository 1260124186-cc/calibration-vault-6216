package service

import (
	"context"
	"sort"

	"example.com/calibration-vault/internal/domain"
)

type Inspection struct {
	SampleID       string        `json:"sample_id"`
	Status         domain.Status `json:"status"`
	StatusLabel    string        `json:"status_label"`
	Fingerprint    string        `json:"fingerprint"`
	EventCount     int           `json:"event_count"`
	Actors         []string      `json:"actors"`
	Complete       bool          `json:"complete"`
	MissingActions []string      `json:"missing_actions,omitempty"`
}

func (s *Service) Inspect(ctx context.Context, sampleID string) (Inspection, error) {
	if err := ctx.Err(); err != nil {
		return Inspection{}, err
	}
	sample, err := s.repository.Get(sampleID)
	if err != nil {
		return Inspection{}, err
	}
	events, err := s.repository.Timeline(sampleID)
	if err != nil {
		return Inspection{}, err
	}
	timeline := domain.NewTimeline(sampleID, events)
	result := Inspection{
		SampleID:    sample.SampleID,
		Status:      sample.Status,
		StatusLabel: domain.StateLabel(sample.Status),
		Fingerprint: domain.FingerprintPrefix(domain.SampleFingerprint(sample), 16),
		EventCount:  timeline.Count,
		Actors:      timeline.Actors(),
		Complete:    domain.IsTerminal(sample.Status),
	}
	if sample.Status == domain.StatusPendingReview {
		result.MissingActions = []string{"review"}
	}
	if sample.Status == domain.StatusApproved {
		result.MissingActions = []string{"release"}
	}
	sort.Strings(result.MissingActions)
	return result, nil
}

func (s *Service) InspectMany(ctx context.Context, sampleIDs []string) ([]Inspection, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	result := make([]Inspection, 0, len(sampleIDs))
	for _, sampleID := range sampleIDs {
		item, err := s.Inspect(ctx, sampleID)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func InspectionFor(sample domain.Sample, events []domain.Event) Inspection {
	timeline := domain.NewTimeline(sample.SampleID, events)
	return Inspection{
		SampleID:    sample.SampleID,
		Status:      sample.Status,
		StatusLabel: domain.StateLabel(sample.Status),
		Fingerprint: domain.FingerprintPrefix(domain.SampleFingerprint(sample), 16),
		EventCount:  timeline.Count,
		Actors:      timeline.Actors(),
		Complete:    domain.IsTerminal(sample.Status),
	}
}
