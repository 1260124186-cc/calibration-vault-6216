package store

import (
	"fmt"

	"example.com/calibration-vault/internal/domain"
)

func (m *Memory) NextEventID() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventIndex++
	return fmt.Sprintf("evt-%06d", m.eventIndex)
}

func (m *Memory) EventCount(sampleID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.events[sampleID])
}

func (m *Memory) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.samples = map[string]domain.Sample{}
	m.events = map[string][]domain.Event{}
	m.eventIndex = 0
}
