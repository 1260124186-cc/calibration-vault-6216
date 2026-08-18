package store_test

import (
	"context"
	"testing"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func TestIntakeWorksAfterMemoryReset(t *testing.T) {
	repository := store.NewMemory()
	repository.Reset()
	application := service.New(repository, service.SystemClock{})
	_, err := application.Intake(context.Background(), domain.IntakeInput{
		SampleID: "after-reset",
		Source:   "line-a",
		Priority: domain.PriorityRoutine,
	})
	if err != nil {
		t.Fatalf("intake after reset failed: %v", err)
	}
}

func TestBatchIntakeWorksAfterMemoryReset(t *testing.T) {
	repository := store.NewMemory()
	repository.Reset()
	application := service.New(repository, service.SystemClock{})
	_, err := application.BatchIntake(context.Background(), domain.BatchIntakeInput{
		BatchReference: "after-reset",
		Items: []domain.IntakeInput{
			{SampleID: "reset-batch-one", Source: "line-a", Priority: domain.PriorityRoutine},
			{SampleID: "reset-batch-two", Source: "line-b", Priority: domain.PriorityUrgent},
		},
	})
	if err != nil {
		t.Fatalf("batch intake after reset failed: %v", err)
	}
}
