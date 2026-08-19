package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"example.com/calibration-vault/internal/domain"
)

type problem struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Fields  []domain.FieldError `json:"fields,omitempty"`
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeProblem(writer http.ResponseWriter, status int, code, message string) {
	writeJSON(writer, status, problem{Code: code, Message: message})
}

func writeServiceError(writer http.ResponseWriter, err error) {
	var validation domain.ValidationErrors
	switch {
	case errors.As(err, &validation):
		writeJSON(writer, http.StatusBadRequest, problem{
			Code:    "validation_failed",
			Message: "request validation failed",
			Fields:  validation.Items,
		})
	case errors.Is(err, domain.ErrDuplicateSample):
		writeProblem(writer, http.StatusConflict, "duplicate_sample", err.Error())
	case errors.Is(err, domain.ErrSampleNotFound):
		writeProblem(writer, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, domain.ErrInvalidState):
		writeProblem(writer, http.StatusConflict, "invalid_state", err.Error())
	case errors.Is(err, domain.ErrInvalidFilter):
		writeProblem(writer, http.StatusBadRequest, "invalid_filter", err.Error())
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		writeProblem(writer, http.StatusRequestTimeout, "request_cancelled", err.Error())
	default:
		writeProblem(writer, http.StatusInternalServerError, "internal_error", "unexpected service error")
	}
}
