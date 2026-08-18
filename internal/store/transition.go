package store

import (
	"example.com/calibration-vault/internal/domain"
)

func (m *Memory) ApplyReview(sampleID string, review domain.Review, event domain.Event) (domain.Sample, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sample, ok := m.samples[sampleID]
	if !ok {
		return domain.Sample{}, domain.ErrSampleNotFound
	}
	if !sample.CanReview() {
		return domain.Sample{}, domain.ErrInvalidState
	}
	if _, err := domain.TransitionForReview(sample, review); err != nil {
		return domain.Sample{}, err
	}
	reviewCopy := review.Clone()
	sample.Review = &reviewCopy
	sample.UpdatedAt = review.ReviewedAt
	sample.Status = domain.ReviewTarget(review.Decision)
	m.samples[sampleID] = sample
	m.events[sampleID] = append(m.events[sampleID], event.Clone())
	return sample.Clone(), nil
}

func (m *Memory) ApplyRelease(sampleID string, release domain.Release, event domain.Event) (domain.Sample, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sample, ok := m.samples[sampleID]
	if !ok {
		return domain.Sample{}, domain.ErrSampleNotFound
	}
	if !sample.CanRelease() {
		return domain.Sample{}, domain.ErrInvalidState
	}
	if _, err := domain.TransitionForRelease(sample, release); err != nil {
		return domain.Sample{}, err
	}
	releaseCopy := release.Clone()
	sample.Release = &releaseCopy
	sample.Status = domain.StatusReleased
	sample.UpdatedAt = release.ReleasedAt
	m.samples[sampleID] = sample
	m.events[sampleID] = append(m.events[sampleID], event.Clone())
	return sample.Clone(), nil
}

func (m *Memory) Replace(sample domain.Sample) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.contains(sample.SampleID) {
		return domain.ErrSampleNotFound
	}
	m.samples[sample.SampleID] = sample.Clone()
	return nil
}
