package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindOrCreateRequestID_ExistingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "existing-id")

	requestID := findOrCreateRequestID(req)
	if requestID != "existing-id" {
		t.Fatalf("expected 'existing-id', got %s", requestID)
	}
}

func TestFindOrCreateRequestID_NewID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	requestID := findOrCreateRequestID(req)
	if requestID == "" {
		t.Fatalf("expected non-empty request ID")
	}
}
