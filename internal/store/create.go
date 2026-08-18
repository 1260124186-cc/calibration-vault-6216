package store

import "example.com/calibration-vault/internal/domain"

func (m *Memory) Create(sample domain.Sample, event domain.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.contains(sample.SampleID) {
		return domain.ErrDuplicateSample
	}
	// 入库时克隆，避免外部传入的 Tags 切片与存储共享底层数组
	m.samples[sample.SampleID] = sample.Clone()
	m.events[sample.SampleID] = []domain.Event{event.Clone()}
	return nil
}

func (m *Memory) Exists(sampleID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.contains(sampleID)
}

func (m *Memory) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.samples)
}
