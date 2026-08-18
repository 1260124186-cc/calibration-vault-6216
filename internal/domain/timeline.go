package domain

import "sort"

type Timeline struct {
	SampleID string  `json:"sample_id"`
	Events   []Event `json:"events"`
	Count    int     `json:"count"`
}

func NewTimeline(sampleID string, events []Event) Timeline {
	copy := append([]Event(nil), events...)
	sort.SliceStable(copy, func(i, j int) bool {
		if copy[i].CreatedAt.Equal(copy[j].CreatedAt) {
			return copy[i].ID < copy[j].ID
		}
		return copy[i].CreatedAt.Before(copy[j].CreatedAt)
	})
	return Timeline{SampleID: sampleID, Events: copy, Count: len(copy)}
}

func (t Timeline) LastKind() EventKind {
	if len(t.Events) == 0 {
		return ""
	}
	return t.Events[len(t.Events)-1].Kind
}

func (t Timeline) HasRelease() bool {
	return t.LastKind() == EventRelease
}

func (t Timeline) HasReview() bool {
	for _, event := range t.Events {
		if event.Kind == EventReview || event.Kind == EventRejection {
			return true
		}
	}
	return false
}

func (t Timeline) Actors() []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(t.Events))
	for _, event := range t.Events {
		if _, exists := seen[event.Actor]; exists {
			continue
		}
		seen[event.Actor] = struct{}{}
		result = append(result, event.Actor)
	}
	sort.Strings(result)
	return result
}
