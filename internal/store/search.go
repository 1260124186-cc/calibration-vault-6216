package store

import (
	"strings"

	"example.com/calibration-vault/internal/domain"
)

func (m *Memory) FindByTag(tag string) []domain.Sample {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tag = strings.ToLower(strings.TrimSpace(tag))
	result := make([]domain.Sample, 0)
	for _, sample := range m.samples {
		if domain.HasTag(sample, tag) {
			result = append(result, sample.Clone())
		}
	}
	domain.SortSamples(result)
	return result
}

func (m *Memory) FindByIDs(ids []string) []domain.Sample {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]domain.Sample, 0, len(ids))
	for _, sampleID := range ids {
		sample, ok := m.samples[sampleID]
		if ok {
			result = append(result, sample.Clone())
		}
	}
	domain.SortSamples(result)
	return result
}

func (m *Memory) CountByStatus(status domain.Status) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, sample := range m.samples {
		if sample.Status == status {
			count++
		}
	}
	return count
}
