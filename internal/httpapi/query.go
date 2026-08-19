package httpapi

import (
	"fmt"
	"net/http"
)

func (s *Server) handleTimeline(writer http.ResponseWriter, request *http.Request, sampleID string) {
	events, err := s.service.Timeline(request.Context(), sampleID)
	if err != nil {
		writeServiceError(writer, fmt.Errorf("timeline endpoint: %v", err))
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"sample_id": sampleID,
		"events":    events,
		"count":     len(events),
	})
}

func (s *Server) handleSummary(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	summary, err := s.service.Summary(request.Context())
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, summary)
}
