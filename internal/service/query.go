package service

import (
	"context"

	"example.com/calibration-vault/internal/domain"
)

type ListOptions struct {
	Filter domain.Filter
	Offset int
	Limit  int
}

func (s *Service) Get(ctx context.Context, sampleID string) (domain.Sample, error) {
	if err := ctx.Err(); err != nil {
		return domain.Sample{}, err
	}
	return s.repository.Get(sampleID)
}

func (s *Service) List(ctx context.Context, options ListOptions) (domain.SamplePage, error) {
	if err := ctx.Err(); err != nil {
		return domain.SamplePage{}, err
	}
	items := s.repository.List(options.Filter)
	return domain.PaginateSamples(items, options.Offset, options.Limit), nil
}

func (s *Service) Timeline(ctx context.Context, sampleID string) ([]domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.repository.Timeline(sampleID)
}

func (s *Service) Summary(ctx context.Context) (domain.Summary, error) {
	if err := ctx.Err(); err != nil {
		return domain.Summary{}, err
	}
	return s.repository.Summary(), nil
}
