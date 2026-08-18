package store

import (
	"example.com/calibration-vault/internal/domain"
)

type BatchRepository interface {
	CreateBatch(samples []domain.Sample, events []domain.Event) error
}

func (m *Memory) CreateBatch(samples []domain.Sample, events []domain.Event) error {
	if len(samples) != len(events) {
		return domain.ErrInvalidSample
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	seen := map[string]struct{}{}
	for _, sample := range samples {
		if m.contains(sample.SampleID) {
			return domain.ErrDuplicateSample
		}
		if _, exists := seen[sample.SampleID]; exists {
			return domain.ErrDuplicateSample
		}
		seen[sample.SampleID] = struct{}{}
	}
	for index, sample := range samples {
		m.ensureEventStoreLocked()
		m.samples[sample.SampleID] = sample.Clone()
		m.events[sample.SampleID] = []domain.Event{events[index].Clone()}
	}
	return nil
}

func (m *Memory) SampleIDs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]string, 0, len(m.samples))
	for sampleID := range m.samples {
		result = append(result, sampleID)
	}
	return result
}

func (m *Memory) EventsFor(sampleIDs []string) map[string][]domain.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string][]domain.Event, len(sampleIDs))
	for _, sampleID := range sampleIDs {
		result[sampleID] = m.cloneEvents(sampleID)
	}
	return result
}
