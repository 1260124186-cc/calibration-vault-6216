package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func TestTimelineForUnknownSampleReturnsNotFound(t *testing.T) {
	handler := New(service.New(store.NewMemory(), service.SystemClock{}))
	request := httptest.NewRequest(http.MethodGet, "/v1/intakes/missing-sample/timeline", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusNotFound, recorder.Body.String())
	}
}
