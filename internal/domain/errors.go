package domain

import "errors"

var (
	ErrInvalidSample   = errors.New("invalid sample")
	ErrDuplicateSample = errors.New("sample already exists")
	ErrSampleNotFound  = errors.New("sample not found")
	ErrInvalidState    = errors.New("invalid sample state")
	ErrInvalidReview   = errors.New("invalid review")
	ErrInvalidRelease  = errors.New("invalid release")
	ErrEventNotFound   = errors.New("timeline not found")
	ErrInvalidFilter   = errors.New("invalid filter")
)

type FieldError struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

func (e FieldError) Error() string {
	return e.Field + ": " + e.Issue
}

type ValidationErrors struct {
	Items []FieldError `json:"items"`
}

func (e ValidationErrors) Error() string {
	if len(e.Items) == 0 {
		return ErrInvalidSample.Error()
	}
	return e.Items[0].Error()
}

func AddValidation(items *[]FieldError, field, issue string) {
	*items = append(*items, FieldError{Field: field, Issue: issue})
}
