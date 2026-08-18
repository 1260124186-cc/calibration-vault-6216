package httpapi

import (
	"net/http"
	"time"
)

func (s *Server) handleHealth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeProblem(writer, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().UTC(),
	})
}
