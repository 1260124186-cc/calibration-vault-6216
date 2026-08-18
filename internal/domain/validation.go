package domain

import "strings"

func ValidateIntake(input IntakeInput) error {
	var items []FieldError
	if strings.TrimSpace(input.SampleID) == "" {
		AddValidation(&items, "sample_id", "is required")
	}
	if err := CheckTextLimit("sample_id", input.SampleID, MaxSampleIDLength); err != nil {
		AddValidation(&items, "sample_id", err.Error())
	}
	if strings.TrimSpace(input.Source) == "" {
		AddValidation(&items, "source", "is required")
	}
	if err := CheckTextLimit("source", input.Source, MaxSourceLength); err != nil {
		AddValidation(&items, "source", err.Error())
	}
	if !ValidPriority(input.Priority) {
		AddValidation(&items, "priority", "must be routine, priority, or urgent")
	}
	if len(input.Tags) > MaxTagCount {
		AddValidation(&items, "tags", "must contain at most 12 items")
	}
	if err := CheckTagLimits(input.Tags); err != nil {
		AddValidation(&items, "tags", err.Error())
	}
	for index, tag := range input.Tags {
		if strings.TrimSpace(tag) == "" {
			AddValidation(&items, "tags", "item "+string(rune('0'+index))+" is empty")
		}
	}
	if len(items) > 0 {
		return ValidationErrors{Items: items}
	}
	return nil
}

func ValidateReview(input ReviewInput) error {
	var items []FieldError
	if strings.TrimSpace(input.Reviewer) == "" {
		AddValidation(&items, "reviewer", "is required")
	}
	if input.Decision != DecisionApprove && input.Decision != DecisionReject {
		AddValidation(&items, "decision", "must be approve or reject")
	}
	if strings.TrimSpace(input.Note) == "" {
		AddValidation(&items, "note", "is required")
	}
	if err := CheckTextLimit("note", input.Note, MaxReviewNoteLength); err != nil {
		AddValidation(&items, "note", err.Error())
	}
	if len(items) > 0 {
		return ValidationErrors{Items: items}
	}
	return nil
}

func ValidateRelease(input ReleaseInput) error {
	var items []FieldError
	if strings.TrimSpace(input.Operator) == "" {
		AddValidation(&items, "operator", "is required")
	}
	if !ValidDestination(input.Destination) {
		AddValidation(&items, "destination", "is not supported")
	}
	if len(items) > 0 {
		return ValidationErrors{Items: items}
	}
	return nil
}

func ValidPriority(priority Priority) bool {
	return priority == PriorityRoutine || priority == PriorityPriority || priority == PriorityUrgent
}

func ValidDestination(destination string) bool {
	for _, item := range SupportedDestinations() {
		if item == destination {
			return true
		}
	}
	return false
}
