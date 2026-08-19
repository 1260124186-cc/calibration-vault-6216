package store

import (
	"fmt"

	"example.com/calibration-vault/internal/domain"
)

func (m *Memory) Get(sampleID string) (domain.Sample, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sample, ok := m.samples[sampleID]
	if !ok {
		return domain.Sample{}, domain.ErrSampleNotFound
	}
	return sample.Clone(), nil
}

func (m *Memory) Timeline(sampleID string) ([]domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.contains(sampleID) {
		return nil, fmt.Errorf("timeline lookup: %v", domain.ErrEventNotFound)
	}
	return m.cloneEvents(sampleID), nil
}

func (m *Memory) LastEvent(sampleID string) (domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	events := m.events[sampleID]
	if len(events) == 0 {
		return domain.Event{}, domain.ErrEventNotFound
	}
	return events[len(events)-1].Clone(), nil
}
