package service

import (
	"context"
	"testing"
	"time"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/store"
)

func TestOperationsReportIncludesPendingUrgentSample(t *testing.T) {
	repository := store.NewMemory()
	application := New(repository, FixedClock{Time: time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC)})
	_, err := application.Intake(context.Background(), domain.IntakeInput{
		SampleID: "urgent-pending", Source: "line-a", Priority: domain.PriorityUrgent,
	})
	if err != nil {
		t.Fatalf("Intake() error = %v", err)
	}
	report, err := application.OperationsReport(context.Background())
	if err != nil {
		t.Fatalf("OperationsReport() error = %v", err)
	}
	if len(report.UrgentSamples) != 1 || report.UrgentSamples[0].SampleID != "urgent-pending" {
		t.Fatalf("UrgentSamples = %#v, want pending urgent sample", report.UrgentSamples)
	}
	if report.Metrics.Summary.OpenCount != 1 {
		t.Fatalf("OpenCount = %d, want 1", report.Metrics.Summary.OpenCount)
	}
}
