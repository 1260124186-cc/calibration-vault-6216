package domain

import (
	"fmt"
	"strings"
)

const (
	MaxSampleIDLength     = 80
	MaxSourceLength       = 120
	MaxReviewNoteLength   = 500
	MaxBatchSize          = 50
	MaxTagCount           = 12
	MaxTagLength          = 40
	MaxBatchReferenceSize = 80
)

func CheckTextLimit(field, value string, limit int) error {
	if len(value) > limit {
		return FieldError{Field: field, Issue: fmt.Sprintf("must be at most %d characters", limit)}
	}
	return nil
}

func CheckTagLimits(tags []string) error {
	var items []FieldError
	for index, raw := range tags {
		tag := strings.TrimSpace(raw)
		if len(tag) > MaxTagLength {
			AddValidation(&items, "tags", fmt.Sprintf("item %d is too long", index))
		}
		if strings.ContainsAny(tag, "\r\n") {
			AddValidation(&items, "tags", fmt.Sprintf("item %d contains a line break", index))
		}
	}
	if len(items) > 0 {
		return ValidationErrors{Items: items}
	}
	return nil
}

func EffectivePriority(priority Priority) int {
	switch priority {
	case PriorityUrgent:
		return 3
	case PriorityPriority:
		return 2
	case PriorityRoutine:
		return 1
	default:
		return 0
	}
}

func ComparePriority(left, right Priority) int {
	leftWeight := EffectivePriority(left)
	rightWeight := EffectivePriority(right)
	switch {
	case leftWeight < rightWeight:
		return -1
	case leftWeight > rightWeight:
		return 1
	default:
		return 0
	}
}

func HigherPriority(left, right Priority) Priority {
	if ComparePriority(left, right) >= 0 {
		return left
	}
	return right
}

func ValidatePage(offset, limit int) error {
	if offset < 0 {
		return FieldError{Field: "offset", Issue: "must not be negative"}
	}
	if limit < 0 || limit > 100 {
		return FieldError{Field: "limit", Issue: "must be between 1 and 100"}
	}
	return nil
}
