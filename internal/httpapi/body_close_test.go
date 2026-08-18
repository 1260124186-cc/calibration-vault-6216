package httpapi_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/calibration-vault/internal/httpapi"
	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error {
	b.closed = true
	return nil
}

func TestIntakeClosesRequestBody(t *testing.T) {
	application := service.New(store.NewMemory(), service.SystemClock{})
	handler := httpapi.New(application)
	body := &trackedBody{Reader: bytes.NewBufferString(`{"sample_id":"body-close","source":"line-a","priority":"routine"}`)}
	request := httptest.NewRequest(http.MethodPost, "/v1/intakes", body)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusCreated)
	}
	if !body.closed {
		t.Fatal("request body was not closed")
	}
}
