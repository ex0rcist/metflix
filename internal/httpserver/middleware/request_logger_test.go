package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindOrCreateRequestID(t *testing.T) {
	tt := []struct {
		reqID string
	}{
		{reqID: "existing-id"},
		{},
	}

	for _, tc := range tt {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		if len(tc.reqID) > 0 {
			req.Header.Set("X-Request-Id", tc.reqID)

			requestID := findOrCreateRequestID(req)
			if requestID != tc.reqID {
				t.Fatalf("expected '%s', got %s", tc.reqID, requestID)
			}
		} else {
			requestID := findOrCreateRequestID(req)
			if requestID == "" {
				t.Fatalf("expected non-empty reqID")
			}
		}
	}
}
