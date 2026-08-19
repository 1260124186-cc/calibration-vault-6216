package service

import (
	"context"
	"testing"
	"time"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/store"
)

func TestBatchIntakeAssignsDistinctEventIDs(t *testing.T) {
	repository := store.NewMemory()
	application := New(repository, FixedClock{Time: time.Date(2026, 8, 18, 10, 0, 0, 0, time.UTC)})
	result, err := application.BatchIntake(context.Background(), domain.BatchIntakeInput{
		BatchReference: "batch-event-ids",
		Items: []domain.IntakeInput{
			{SampleID: "batch-event-1", Source: "line-a", Priority: domain.PriorityRoutine},
			{SampleID: "batch-event-2", Source: "line-a", Priority: domain.PriorityUrgent},
		},
	})
	if err != nil {
		t.Fatalf("BatchIntake() error = %v", err)
	}
	first, err := application.Timeline(context.Background(), result.Items[0].SampleID)
	if err != nil {
		t.Fatalf("Timeline(first) error = %v", err)
	}
	second, err := application.Timeline(context.Background(), result.Items[1].SampleID)
	if err != nil {
		t.Fatalf("Timeline(second) error = %v", err)
	}
	if first[0].ID == second[0].ID {
		t.Fatalf("batch event IDs are both %q, want distinct event IDs", first[0].ID)
	}
}
