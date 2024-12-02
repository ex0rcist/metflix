package httpserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	type result struct {
		code int
	}

	tests := []struct {
		name         string
		pingResponse error
		expected     result
	}{
		{
			name:         "should return ok if storage is ok",
			pingResponse: nil,
			expected:     result{code: http.StatusOK},
		},
		{
			name:         "should return not implemented if storage doesn't support ping",
			pingResponse: entities.ErrStorageUnpingable,
			expected:     result{code: http.StatusNotImplemented},
		},
		{
			name:         "should return internal server error if storage offline",
			pingResponse: entities.ErrUnexpected,
			expected:     result{code: http.StatusInternalServerError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _, healthMock := createHealthTestBackend()
			healthMock.On("Ping", mock.Anything).Return(tt.pingResponse)

			code, _, _ := testHealthRequest(t, router, http.MethodGet, "/ping", nil)

			if tt.expected.code != code {
				t.Fatalf("expected response to be %d, got: %d", tt.expected.code, code)
			}
		})
	}
}

func testHealthRequest(t *testing.T, router http.Handler, method, path string, payload []byte) (int, string, []byte) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	body := bytes.NewReader(payload)

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logging.LogError(closeErr)
		}
	}()

	contentType := resp.Header.Get("Content-Type")

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, contentType, respBody
}

func createHealthTestBackend() (http.Handler, *services.MetricServiceMock, *services.HealthCheckServiceMock) {
	healthServiceMock := &services.HealthCheckServiceMock{}
	healthMock := NewHealthResource(healthServiceMock)

	metricServiceMock := &services.MetricServiceMock{}
	metricMock := NewMetricResource(metricServiceMock)

	handler := NewBackend(
		WithHealthResource(healthMock),
		WithMetricResource(metricMock),
	)

	return handler, metricServiceMock, healthServiceMock
}
