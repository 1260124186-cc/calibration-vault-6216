package store

import (
	"sync"

	"example.com/calibration-vault/internal/domain"
)

type Memory struct {
	mu         sync.RWMutex
	samples    map[string]domain.Sample
	events     map[string][]domain.Event
	eventIndex uint64
}

func NewMemory() *Memory {
	memory := &Memory{
		samples: map[string]domain.Sample{},
		events:  map[string][]domain.Event{},
	}
	RegisterOptionalSurface(memory)
	return memory
}

func (m *Memory) contains(sampleID string) bool {
	_, ok := m.samples[sampleID]
	return ok
}

func (m *Memory) cloneEvents(sampleID string) []domain.Event {
	events := m.events[sampleID]
	result := make([]domain.Event, len(events))
	for index, event := range events {
		result[index] = event.Clone()
	}
	return result
}

func (m *Memory) cloneSamples(filter domain.Filter) []domain.Sample {
	items := make([]domain.Sample, 0, len(m.samples))
	for _, sample := range m.samples {
		if filter.Matches(sample) {
			// 克隆切片字段，避免返回值与存储共享底层数组
			items = append(items, sample.Clone())
		}
	}
	domain.SortSamples(items)
	return items
}
