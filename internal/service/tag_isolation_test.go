package service_test

import (
	"context"
	"testing"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func TestSampleTagsDoNotAliasReturnedValues(t *testing.T) {
	application := service.New(store.NewMemory(), service.SystemClock{})
	sample, err := application.Intake(context.Background(), domain.IntakeInput{
		SampleID: "tag-isolation",
		Source:   "line-a",
		Priority: domain.PriorityRoutine,
		Tags:     []string{"calibration", "urgent-check"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sample.Tags[0] = "caller-mutated"
	loaded, err := application.Get(context.Background(), sample.SampleID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Tags[0] != "calibration" {
		t.Fatalf("stored tag changed to %q", loaded.Tags[0])
	}
	page, err := application.List(context.Background(), service.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	page.Items[0].Tags[1] = "list-mutated"
	loaded, err = application.Get(context.Background(), sample.SampleID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Tags[1] != "urgent-check" {
		t.Fatalf("list result changed stored tag to %q", loaded.Tags[1])
	}
}
