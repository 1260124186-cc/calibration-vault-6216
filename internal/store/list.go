package store

import "example.com/calibration-vault/internal/domain"

func (m *Memory) List(filter domain.Filter) []domain.Sample {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cloneSamples(filter)
}

func (m *Memory) Summary() domain.Summary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	all := m.cloneSamples(domain.Filter{})
	return domain.BuildSummary(all)
}

func (m *Memory) ListOpen() []domain.Sample {
	return m.List(domain.Filter{Status: domain.StatusPendingReview})
}
