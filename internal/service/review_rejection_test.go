package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func TestRejectedReviewStaysRejected(t *testing.T) {
	application := service.New(store.NewMemory(), service.FixedClock{Time: time.Unix(1700000000, 0)})
	sample, err := application.Intake(context.Background(), domain.IntakeInput{
		SampleID: "reject-me",
		Source:   "line-a",
		Priority: domain.PriorityRoutine,
	})
	if err != nil {
		t.Fatal(err)
	}
	reviewed, err := application.Review(context.Background(), sample.SampleID, domain.ReviewInput{
		Reviewer: "reviewer-a",
		Decision: domain.DecisionReject,
		Note:     "metadata is incomplete",
	})
	if err != nil {
		t.Fatal(err)
	}
	if reviewed.Status != domain.StatusRejected {
		t.Fatalf("rejected sample status = %s, want %s", reviewed.Status, domain.StatusRejected)
	}
	_, err = application.Release(context.Background(), sample.SampleID, domain.ReleaseInput{
		Operator:    "operator-a",
		Destination: "lab-west",
	})
	if !errors.Is(err, domain.ErrInvalidState) {
		t.Fatalf("release error = %v, want invalid state", err)
	}
}
