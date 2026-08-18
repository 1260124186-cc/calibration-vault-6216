package domain

import (
	"fmt"
	"strings"
	"time"
)

type BatchIntakeInput struct {
	BatchReference string        `json:"batch_reference"`
	Items          []IntakeInput `json:"items"`
}

type BatchIntakeResult struct {
	BatchReference string    `json:"batch_reference"`
	Items          []Sample  `json:"items"`
	Count          int       `json:"count"`
	CreatedAt      time.Time `json:"created_at"`
}

func ValidateBatchIntake(input BatchIntakeInput) error {
	var items []FieldError
	if strings.TrimSpace(input.BatchReference) == "" {
		AddValidation(&items, "batch_reference", "is required")
	}
	if len(input.BatchReference) > 80 {
		AddValidation(&items, "batch_reference", "must be at most 80 characters")
	}
	if len(input.Items) == 0 {
		AddValidation(&items, "items", "must contain at least one sample")
	}
	if len(input.Items) > 50 {
		AddValidation(&items, "items", "must contain at most 50 samples")
	}
	seen := map[string]struct{}{}
	for index, intake := range input.Items {
		sampleID := strings.TrimSpace(intake.SampleID)
		if _, exists := seen[sampleID]; exists && sampleID != "" {
			AddValidation(&items, "items", fmt.Sprintf("sample_id %q appears more than once", sampleID))
		}
		seen[sampleID] = struct{}{}
		if err := ValidateIntake(intake); err != nil {
			AddValidation(&items, "items", fmt.Sprintf("item %d is invalid: %s", index, err.Error()))
		}
	}
	if len(items) > 0 {
		return ValidationErrors{Items: items}
	}
	return nil
}

func NewBatchResult(reference string, samples []Sample, now time.Time) BatchIntakeResult {
	result := BatchIntakeResult{
		BatchReference: reference,
		Items:          append([]Sample(nil), samples...),
		Count:          len(samples),
		CreatedAt:      now,
	}
	return result
}

func BatchSampleIDs(result BatchIntakeResult) []string {
	ids := make([]string, 0, len(result.Items))
	for _, sample := range result.Items {
		ids = append(ids, sample.SampleID)
	}
	return ids
}

func BatchHasUrgent(result BatchIntakeResult) bool {
	for _, sample := range result.Items {
		if sample.Priority == PriorityUrgent {
			return true
		}
	}
	return false
}
