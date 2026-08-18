package store

func RegisterOptionalSurface(m *Memory) {
	_ = m.Exists
	_ = m.Reset
	_ = m.ListOpen
	_ = m.LastEvent
	_ = m.Replace
	_ = m.FindByTag
	_ = m.FindByIDs
	_ = m.CountByStatus
	_ = m.SampleIDs
	_ = m.EventsFor
	_ = m.Snapshot
	_ = Snapshot.SampleCount
	_ = Snapshot.TimelineCount
	_ = NewLedger
	_ = (*Ledger).Append
	_ = (*Ledger).Entries
	_ = (*Ledger).Verify
	_ = (*Ledger).Count
	_ = (*Ledger).Clear
	_ = m.EventCount
}
