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

// 紧急样本只要仍处于开放流转状态就应纳入报告；到达终态（已退回/已放行）后移出
func TestOperationsReportUrgentSamplesCoversOpenStates(t *testing.T) {
	clock := FixedClock{Time: time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC)}
	application := New(store.NewMemory(), clock)

	// 登记四个紧急样本，分别走向 pending_review / approved / rejected / released
	ids := map[string]string{
		"pending":  "urgent-pending",
		"approved": "urgent-approved",
		"rejected": "urgent-rejected",
		"released": "urgent-released",
	}
	for _, id := range ids {
		if _, err := application.Intake(context.Background(), domain.IntakeInput{
			SampleID: id, Source: "line-a", Priority: domain.PriorityUrgent,
		}); err != nil {
			t.Fatalf("Intake(%s) error = %v", id, err)
		}
	}

	approve := func(id string) {
		if _, err := application.Review(context.Background(), id, domain.ReviewInput{
			Reviewer: "r1", Decision: domain.DecisionApprove, Note: "ok",
		}); err != nil {
			t.Fatalf("Review(%s) error = %v", id, err)
		}
	}
	approve(ids["approved"])
	approve(ids["released"])
	if _, err := application.Review(context.Background(), ids["rejected"], domain.ReviewInput{
		Reviewer: "r1", Decision: domain.DecisionReject, Note: "bad",
	}); err != nil {
		t.Fatalf("Review(rejected) error = %v", err)
	}
	if _, err := application.Release(context.Background(), ids["released"], domain.ReleaseInput{
		Operator: "o1", Destination: "lab-west",
	}); err != nil {
		t.Fatalf("Release() error = %v", err)
	}

	report, err := application.OperationsReport(context.Background())
	if err != nil {
		t.Fatalf("OperationsReport() error = %v", err)
	}

	got := make(map[string]bool, len(report.UrgentSamples))
	for _, sample := range report.UrgentSamples {
		got[sample.SampleID] = true
	}
	want := map[string]bool{ids["pending"]: true, ids["approved"]: true}
	if len(got) != 2 || got[ids["pending"]] != want[ids["pending"]] || got[ids["approved"]] != want[ids["approved"]] {
		t.Fatalf("UrgentSamples = %v, want only pending and approved urgent samples", got)
	}
	if got[ids["rejected"]] || got[ids["released"]] {
		t.Fatalf("terminal urgent samples should not appear, got = %v", got)
	}

	// 开放流转计数应只含 pending 与 approved，排除两个终态
	if report.Metrics.Summary.OpenCount != 2 {
		t.Fatalf("OpenCount = %d, want 2", report.Metrics.Summary.OpenCount)
	}

	// pending urgent 也应出现在待复核清单中，供复核员处理
	found := false
	for _, sample := range report.PendingReview {
		if sample.SampleID == ids["pending"] {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("pending urgent sample missing from PendingReview")
	}
}
