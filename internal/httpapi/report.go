package httpapi

import (
	"net/http"
	"strings"
)

func (s *Server) handleOperationsReport(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	report, err := s.service.OperationsReport(request.Context())
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writer.Header().Set("Cache-Control", "no-store")
	writeJSON(writer, http.StatusOK, report)
}

func (s *Server) handleMetadata(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	metadata, err := s.service.Metadata(request.Context())
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	metadata["service"] = "calibration-vault"
	writeJSON(writer, http.StatusOK, metadata)
}

func (s *Server) handleTagSearch(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	tag := request.URL.Query().Get("tag")
	items, err := s.service.SearchTag(request.Context(), tag)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"tag":   tag,
		"items": items,
		"count": len(items),
	})
}

func (s *Server) handleInspection(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	sampleID := strings.TrimPrefix(request.URL.Path, "/v1/inspect/")
	if sampleID == "" || strings.Contains(sampleID, "/") {
		writeProblem(writer, http.StatusBadRequest, "invalid_sample_id", "sample id is required")
		return
	}
	inspection, err := s.service.Inspect(request.Context(), sampleID)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, inspection)
}
