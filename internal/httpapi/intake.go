package httpapi

import (
	"net/http"
	"strconv"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/service"
)

func (s *Server) handleIntakes(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		s.handleCreateIntake(writer, request)
	case http.MethodGet:
		s.handleListIntakes(writer, request)
	default:
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use GET or POST")
	}
}

func (s *Server) handleCreateIntake(writer http.ResponseWriter, request *http.Request) {
	var input domain.IntakeInput
	if err := decodeBody(request, &input); err != nil {
		rejectMalformedBody(writer, err)
		return
	}
	sample, err := s.service.Intake(request.Context(), input)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, sample)
}

func (s *Server) handleListIntakes(writer http.ResponseWriter, request *http.Request) {
	filter, err := domain.NormalizeFilter(
		request.URL.Query().Get("status"),
		request.URL.Query().Get("priority"),
		request.URL.Query().Get("source"),
	)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	options := service.ListOptions{
		Filter: filter,
		Offset: parseNonNegative(request.URL.Query().Get("offset")),
		Limit:  parsePositive(request.URL.Query().Get("limit")),
	}
	page, err := s.service.List(request.Context(), options)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, page)
}

func parseNonNegative(raw string) int {
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0
	}
	return value
}

func parsePositive(raw string) int {
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 50
	}
	return value
}
