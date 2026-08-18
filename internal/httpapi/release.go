package httpapi

import (
	"net/http"

	"example.com/calibration-vault/internal/domain"
)

func (s *Server) handleRelease(writer http.ResponseWriter, request *http.Request, sampleID string) {
	var input domain.ReleaseInput
	if err := decodeBody(request, &input); err != nil {
		rejectMalformedBody(writer, err)
		return
	}
	sample, err := s.service.Release(request.Context(), sampleID, input)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writer.Header().Set("X-Release-Destination", input.Destination)
	writeJSON(writer, http.StatusOK, sample)
}
