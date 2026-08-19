package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func TestBatchValidationReturnsBadRequest(t *testing.T) {
	handler := New(service.New(store.NewMemory(), service.SystemClock{}))
	request := httptest.NewRequest(http.MethodPost, "/v1/intakes/batch", strings.NewReader(`{"batch_reference":"daily-1","items":[]}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
}
