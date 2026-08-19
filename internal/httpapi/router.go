package httpapi

import (
	"net/http"
	"strings"

	"example.com/calibration-vault/internal/service"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(service *service.Service) *Server {
	server := &Server{service: service, mux: http.NewServeMux()}
	server.register()
	return server
}

func (s *Server) register() {
	s.mux.HandleFunc("/healthz", s.handleHealth)
	s.mux.HandleFunc("/v1/intakes", s.handleIntakes)
	s.mux.HandleFunc("/v1/intakes/", s.handleIntakeRoute)
	s.mux.HandleFunc("/v1/summary", s.handleSummary)
	s.mux.HandleFunc("/v1/reports/operations", s.handleOperationsReport)
	s.mux.HandleFunc("/v1/metadata", s.handleMetadata)
	s.mux.HandleFunc("/v1/search/tags", s.handleTagSearch)
	s.mux.HandleFunc("/v1/inspect/", s.handleInspection)
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 保留原始请求 context，使客户端断连或取消的信号能沿调用链传递到业务层
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) handleIntakeRoute(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(request.URL.Path, "/v1/intakes/")
	if strings.Trim(path, "/") == "batch" {
		s.handleBatchIntake(writer, request)
		return
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" {
		writeProblem(writer, http.StatusNotFound, "not_found", "route not found")
		return
	}
	sampleID, action := parts[0], parts[1]
	switch {
	case request.Method == http.MethodPost && action == "review":
		s.handleReview(writer, request, sampleID)
	case request.Method == http.MethodPost && action == "release":
		s.handleRelease(writer, request, sampleID)
	case request.Method == http.MethodGet && action == "timeline":
		s.handleTimeline(writer, request, sampleID)
	default:
		writeProblem(writer, http.StatusNotFound, "not_found", "route not found")
	}
}
