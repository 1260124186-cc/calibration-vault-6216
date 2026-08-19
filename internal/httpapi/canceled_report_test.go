package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func TestCanceledOperationsReportReturnsRequestCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	handler := New(service.New(store.NewMemory(), service.SystemClock{}))
	request := httptest.NewRequest(http.MethodGet, "/v1/reports/operations", nil).WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusRequestTimeout {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusRequestTimeout, recorder.Body.String())
	}
}
