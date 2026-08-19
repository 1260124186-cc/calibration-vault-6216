package store

import (
	"example.com/calibration-vault/internal/domain"
)

type Repository interface {
	Create(sample domain.Sample, event domain.Event) error
	Get(sampleID string) (domain.Sample, error)
	List(filter domain.Filter) []domain.Sample
	Timeline(sampleID string) ([]domain.Event, error)
	ApplyReview(sampleID string, review domain.Review, event domain.Event) (domain.Sample, error)
	ApplyRelease(sampleID string, release domain.Release, event domain.Event) (domain.Sample, error)
	NextEventID() string
	NextEventIDs(count int) []string
	Summary() domain.Summary
}
