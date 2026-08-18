package httpapi

import (
	"net/http"

	"example.com/calibration-vault/internal/domain"
)

func (s *Server) handleBatchIntake(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	var input domain.BatchIntakeInput
	if err := decodeBody(request, &input); err != nil {
		rejectMalformedBody(writer, err)
		return
	}
	result, err := s.service.BatchIntake(request.Context(), input)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, result)
}
