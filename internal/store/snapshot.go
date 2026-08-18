package store

import (
	"example.com/calibration-vault/internal/domain"
)

type Snapshot struct {
	Samples []domain.Sample
	Events  map[string][]domain.Event
}

func (m *Memory) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	snapshot := Snapshot{
		Samples: m.cloneSamples(domain.Filter{}),
		Events:  make(map[string][]domain.Event, len(m.events)),
	}
	for sampleID := range m.events {
		snapshot.Events[sampleID] = m.cloneEvents(sampleID)
	}
	return snapshot
}

func (s Snapshot) SampleCount() int {
	return len(s.Samples)
}

func (s Snapshot) TimelineCount() int {
	total := 0
	for _, timeline := range s.Events {
		total += len(timeline)
	}
	return total
}
