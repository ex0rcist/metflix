package middleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilterUntrustedRequest(t *testing.T) {
	tests := []struct {
		trustedSubnet    string
		requestIP        string
		expectedNextCall bool
		expectedStatus   int
	}{
		{"", "192.168.1.1", true, http.StatusOK},
		{"192.168.1.0/24", "192.168.1.50", true, http.StatusOK},
		{"192.168.1.0/24", "10.0.0.1", false, http.StatusForbidden},
		{"192.168.1.0/24", "invalid-ip", false, http.StatusForbidden},
	}

	for _, tt := range tests {

		nextHandlerCalled := false // initial
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextHandlerCalled = true // set true if middleware pass
		})

		var trustedSubnet *net.IPNet
		if len(tt.trustedSubnet) > 0 {
			_, trustedSubnet, _ = net.ParseCIDR(tt.trustedSubnet)
		}

		middleware := FilterUntrustedRequest(nextHandler, trustedSubnet)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Real-IP", tt.requestIP)

		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)

		if tt.expectedNextCall && !nextHandlerCalled {
			t.Fatal("Expected next handler to be called")
		}

		if !tt.expectedNextCall && nextHandlerCalled {
			t.Fatal("Expected next handler to not be called")
		}

		if rr.Code != tt.expectedStatus {
			t.Fatalf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
		}
	}
}
