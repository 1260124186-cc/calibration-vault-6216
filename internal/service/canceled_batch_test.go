package service_test

import (
	"context"
	"errors"
	"testing"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func TestCanceledBatchIntakeDoesNotPersist(t *testing.T) {
	repository := store.NewMemory()
	application := service.New(repository, service.SystemClock{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := application.BatchIntake(ctx, domain.BatchIntakeInput{
		BatchReference: "canceled-request",
		Items: []domain.IntakeInput{
			{SampleID: "canceled-one", Source: "line-a", Priority: domain.PriorityRoutine},
			{SampleID: "canceled-two", Source: "line-b", Priority: domain.PriorityPriority},
		},
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("batch error = %v, want context canceled", err)
	}
	if got := repository.Count(); got != 0 {
		t.Fatalf("stored sample count = %d, want 0", got)
	}
}
