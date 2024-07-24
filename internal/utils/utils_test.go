package utils

import (
	"net/http"
	"regexp"
	"testing"
)

func TestHeadersToStr(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		expected string
	}{
		{
			name: "Single header with single value",
			headers: http.Header{
				"Content-Type": {"application/json"},
			},
			expected: "Content-Type:application/json",
		},
		{
			name: "Single header with multiple values",
			headers: http.Header{
				"Accept": {"text/plain", "text/html"},
			},
			expected: "Accept:text/html, Accept:text/plain",
		},
		{
			name: "Multiple headers with single values",
			headers: http.Header{
				"Content-Type": {"application/json"},
				"User-Agent":   {"Go-http-client/1.1"},
			},
			expected: "Content-Type:application/json, User-Agent:Go-http-client/1.1",
		},
		{
			name:     "Empty headers",
			headers:  http.Header{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HeadersToStr(tt.headers)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	regex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

	requestID := GenerateRequestID()
	if !regex.MatchString(requestID) {
		t.Errorf("GenerateRequestID() returned invalid UUIDv4: %s", requestID)
	}
}
