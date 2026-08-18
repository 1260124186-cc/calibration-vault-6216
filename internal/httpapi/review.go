package httpapi

import (
	"net/http"

	"example.com/calibration-vault/internal/domain"
)

func (s *Server) handleReview(writer http.ResponseWriter, request *http.Request, sampleID string) {
	var input domain.ReviewInput
	if err := decodeBody(request, &input); err != nil {
		rejectMalformedBody(writer, err)
		return
	}
	sample, err := s.service.Review(request.Context(), sampleID, input)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writer.Header().Set("X-Sample-Status", string(sample.Status))
	writeJSON(writer, http.StatusOK, sample)
}
