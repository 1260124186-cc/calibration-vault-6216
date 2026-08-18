package store

import (
	"sync"

	"example.com/calibration-vault/internal/domain"
)

type LedgerEntry struct {
	SampleID     string `json:"sample_id"`
	EventID      string `json:"event_id"`
	SampleDigest string `json:"sample_digest"`
	EventDigest  string `json:"event_digest"`
	Sequence     int    `json:"sequence"`
}

type Ledger struct {
	mu      sync.RWMutex
	entries map[string][]LedgerEntry
}

func NewLedger() *Ledger {
	return &Ledger{entries: map[string][]LedgerEntry{}}
}

func (l *Ledger) Append(sample domain.Sample, event domain.Event) LedgerEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	items := l.entries[sample.SampleID]
	entry := LedgerEntry{
		SampleID:     sample.SampleID,
		EventID:      event.ID,
		SampleDigest: domain.SampleFingerprint(sample),
		EventDigest:  domain.EventFingerprint(event),
		Sequence:     len(items) + 1,
	}
	l.entries[sample.SampleID] = append(items, entry)
	return entry
}

func (l *Ledger) Entries(sampleID string) []LedgerEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	items := l.entries[sampleID]
	return append([]LedgerEntry(nil), items...)
}

func (l *Ledger) Verify(sampleID string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	items := l.entries[sampleID]
	for index, entry := range items {
		if entry.Sequence != index+1 || entry.SampleID != sampleID || entry.EventID == "" {
			return false
		}
	}
	return true
}

func (l *Ledger) Count(sampleID string) int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries[sampleID])
}

func (l *Ledger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = map[string][]LedgerEntry{}
}
