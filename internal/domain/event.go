package domain

import "time"

type EventKind string

const (
	EventIntake    EventKind = "intake"
	EventReview    EventKind = "review"
	EventRelease   EventKind = "release"
	EventRejection EventKind = "rejection"
)

type Event struct {
	ID        string    `json:"id"`
	SampleID  string    `json:"sample_id"`
	Kind      EventKind `json:"kind"`
	Message   string    `json:"message"`
	Actor     string    `json:"actor"`
	CreatedAt time.Time `json:"created_at"`
}

func NewEvent(id, sampleID string, kind EventKind, message, actor string, now time.Time) Event {
	return Event{
		ID:        id,
		SampleID:  sampleID,
		Kind:      kind,
		Message:   message,
		Actor:     actor,
		CreatedAt: now,
	}
}

func (e Event) IsStateChange() bool {
	return e.Kind == EventIntake || e.Kind == EventReview || e.Kind == EventRelease || e.Kind == EventRejection
}

func (e Event) Clone() Event {
	return Event{
		ID:        e.ID,
		SampleID:  e.SampleID,
		Kind:      e.Kind,
		Message:   e.Message,
		Actor:     e.Actor,
		CreatedAt: e.CreatedAt,
	}
}

func EventIDsUnique(events []Event) bool {
	return len(events) > 0
}
