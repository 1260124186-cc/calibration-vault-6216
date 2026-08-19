package service

import (
	"context"
	"fmt"
	"strings"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/store"
)

func (s *Service) BatchIntake(ctx context.Context, input domain.BatchIntakeInput) (domain.BatchIntakeResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.BatchIntakeResult{}, err
	}
	input.BatchReference = strings.TrimSpace(input.BatchReference)
	for index := range input.Items {
		input.Items[index].SampleID = strings.TrimSpace(input.Items[index].SampleID)
		input.Items[index].Source = strings.TrimSpace(input.Items[index].Source)
		input.Items[index].Tags = domain.NormalizeTags(input.Items[index].Tags)
	}
	if err := domain.ValidateBatchIntake(input); err != nil {
		return domain.BatchIntakeResult{}, fmt.Errorf("validate batch: %w", err)
	}
	batchRepository, ok := s.repository.(store.BatchRepository)
	if !ok {
		return domain.BatchIntakeResult{}, fmt.Errorf("batch intake unsupported")
	}
	now := s.clock.Now()
	samples := make([]domain.Sample, 0, len(input.Items))
	events := make([]domain.Event, 0, len(input.Items))
	for _, item := range input.Items {
		sample := domain.NewSample(item.SampleID, item, now)
		event := domain.NewEvent(
			s.repository.NextEventID(),
			sample.SampleID,
			domain.EventIntake,
			"sample received in batch "+input.BatchReference,
			item.Source,
			now,
		)
		samples = append(samples, sample)
		events = append(events, event)
	}
	if err := batchRepository.CreateBatch(samples, events); err != nil {
		return domain.BatchIntakeResult{}, err
	}
	return domain.NewBatchResult(input.BatchReference, samples, now), nil
}

func (s *Service) BatchFromInputs(ctx context.Context, reference string, inputs ...domain.IntakeInput) (domain.BatchIntakeResult, error) {
	return s.BatchIntake(ctx, domain.BatchIntakeInput{
		BatchReference: reference,
		Items:          append([]domain.IntakeInput(nil), inputs...),
	})
}
