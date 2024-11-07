package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, router http.Handler, method, path string, payload []byte) (int, string, []byte) {
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

func createTestRouter() (http.Handler, *storage.ServiceMock, *services.PingerMock) {
	sm := &storage.ServiceMock{}
	pm := &services.PingerMock{}

	router := NewRouter(sm, pm, entities.Secret(""))

	return router, sm, pm
}

func TestHomepage(t *testing.T) {
	type result struct {
		code     int
		body     string
		contains []string
	}

	require := require.New(t)
	tests := []struct {
		name    string
		path    string
		metrics []storage.Record
		want    result
	}{
		{
			name: "default homepage",
			path: "/",
			want: result{code: http.StatusOK, body: "mainpage here.\n"},
		},
		{
			name: "has some metrics",
			path: "/",
			metrics: []storage.Record{
				{Name: "test1", Value: metrics.Counter(1)},
				{Name: "test2", Value: metrics.Gauge(2.3)},
			},
			want: result{code: http.StatusOK, contains: []string{"mainpage here.", "metrics list", "test1 => counter: 1", "test2 => gauge: 2.3"}},
		},
	}

	for _, tt := range tests {
		router, sm, _ := createTestRouter()
		sm.On("List").Return(tt.metrics, nil)

		t.Run(tt.name, func(t *testing.T) {
			code, _, body := testRequest(t, router, http.MethodGet, tt.path, nil)

			require.Equal(tt.want.code, code)

			if len(tt.want.body) > 0 {
				require.Equal(tt.want.body, string(body))
			}

			if len(tt.want.contains) > 0 {
				for _, v := range tt.want.contains {
					require.Contains(string(body), v)
				}
			}
		})
	}
}

func TestUpdateMetric(t *testing.T) {
	type result struct {
		code int
		body string
	}

	require := require.New(t)
	tests := []struct {
		name string
		path string
		mock func(m *storage.ServiceMock)
		want result
	}{
		{
			name: "Should push counter",
			path: "/update/counter/test/42",
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindCounter).Return(storage.Record{}, nil)
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42"},
		},
		{
			name: "Should push counter with existing value",
			path: "/update/counter/test/42",
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindCounter).Return(storage.Record{Name: "test", Value: metrics.Counter(21)}, nil)
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42"},
		},
		{
			name: "Should push gauge",
			path: "/update/gauge/test/42.42",
			mock: func(m *storage.ServiceMock) {
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42.42"},
		},
		{
			name: "Should fail on invalid kind",
			path: "/update/xxx/test/1",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "Should fail on empty metric name",
			path: "/update/counter//1",
			want: result{code: http.StatusNotFound},
		},
		{
			name: "Should fail on counter with invalid name",
			path: "/update/counter/inva!id/10",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "Should fail on counter with invalid value",
			path: "/update/counter/test/10.0",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "Should fail on gauge with invalid name",
			path: "/update/gauge/inval!d/42.42",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "Should fail on gauge with invalid value",
			path: "/update/gauge/test/42.42!",
			want: result{code: http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		router, sm, _ := createTestRouter()

		if tt.mock != nil {
			tt.mock(sm)
		}

		t.Run(tt.name, func(t *testing.T) {
			code, _, body := testRequest(t, router, http.MethodPost, tt.path, nil)

			require.Equal(tt.want.code, code)
			require.Equal(tt.want.body, string(body))
		})
	}
}

func TestUpdateJSONMetric(t *testing.T) {
	type result struct {
		code int
	}

	tests := []struct {
		name string
		mex  metrics.MetricExchange
		mock func(m *storage.ServiceMock)
		want result
	}{
		{
			name: "Should push counter",
			mex:  metrics.NewUpdateCounterMex("test", 42),
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindCounter).Return(storage.Record{}, nil)
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			want: result{code: http.StatusOK},
		},
		{
			name: "Should push gauge",
			mex:  metrics.NewUpdateGaugeMex("test", 42.42),
			mock: func(m *storage.ServiceMock) {
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)}, nil)
			},
			want: result{code: http.StatusOK},
		},
		{
			name: "Should fail on unknown metric kind",
			mex:  metrics.MetricExchange{ID: "42", MType: "test"},
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "Should fail on counter with invalid name",
			mex:  metrics.NewUpdateCounterMex("X)", 10),
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "Should fail on gauge with invalid name",
			mex:  metrics.NewUpdateGaugeMex("X;", 13.123),
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "Should fail on broken storage",
			mex:  metrics.NewUpdateCounterMex("fail", 13),
			mock: func(m *storage.ServiceMock) {
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{}, entities.ErrUnexpected)
			},
			want: result{code: http.StatusInternalServerError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			router, sm, _ := createTestRouter()

			if tt.mock != nil {
				tt.mock(sm)
			}

			payload, err := json.Marshal(tt.mex)
			require.NoError(err)

			code, contentType, body := testRequest(t, router, http.MethodPost, "/update", payload)

			assert.Equal(tt.want.code, code)

			if tt.want.code == http.StatusOK {
				assert.Equal("application/json", contentType)

				var resp metrics.MetricExchange
				err = json.Unmarshal(body, &resp)
				require.NoError(err)

				assert.Equal(tt.mex, resp)
			}
		})
	}
}

func TestBatchUpdateMetricsJSON(t *testing.T) {
	type result struct {
		code int
	}

	batchRequest := []metrics.MetricExchange{
		metrics.NewUpdateCounterMex("PollCount", 42),
		metrics.NewUpdateGaugeMex("Alloc", 42.42),
	}

	batchResponse := []storage.Record{
		{Name: "PollCount", Value: metrics.Counter(42)},
		{Name: "Alloc", Value: metrics.Gauge(42.42)},
	}

	tests := []struct {
		name        string
		mex         []metrics.MetricExchange
		mock        func(m *storage.ServiceMock)
		recorderRv  []storage.Record
		recorderErr error
		expected    result
	}{
		{
			name: "should push different metrics",
			mex:  batchRequest,
			mock: func(m *storage.ServiceMock) {
				m.On("PushList", mock.Anything, mock.Anything).Return(batchResponse, nil)
			},
			expected: result{code: http.StatusOK},
		},
		{
			name:     "should fail on empty list",
			mex:      make([]metrics.MetricExchange, 0),
			expected: result{code: http.StatusBadRequest},
		},
		{
			name:     "should fail on no counter value",
			mex:      []metrics.MetricExchange{{ID: "xxx", MType: "counter"}},
			expected: result{code: http.StatusBadRequest},
		},
		{
			name:     "should fail on no gauge value",
			mex:      []metrics.MetricExchange{{ID: "xxx", MType: "gauge"}},
			expected: result{code: http.StatusBadRequest},
		},
		{
			name:     "should fail on unknown metric",
			mex:      []metrics.MetricExchange{{ID: "xxx", MType: "unknown"}},
			expected: result{code: http.StatusBadRequest},
		},
		{
			name: "should fail if storage offline",
			mex:  batchRequest,
			mock: func(m *storage.ServiceMock) {
				m.On("PushList", mock.Anything, mock.Anything).Return([]storage.Record{}, entities.ErrUnexpected)
			},
			expected: result{code: http.StatusInternalServerError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			router, sm, _ := createTestRouter()
			if tt.mock != nil {
				tt.mock(sm)
			}

			payload, err := json.Marshal(tt.mex)
			require.NoError(err)

			code, _, _ := testRequest(t, router, http.MethodPost, "/updates", payload)
			require.Equal(tt.expected.code, code)
		})
	}
}

func TestGetMetric(t *testing.T) {
	type result struct {
		code int
		body string
	}

	tests := []struct {
		name string
		path string
		mock func(m *storage.ServiceMock)
		want result
	}{
		{
			name: "get counter",
			path: "/value/counter/test",
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindCounter).Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42"},
		},
		{
			name: "get gauge",
			path: "/value/gauge/test",
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindGauge).Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42.42"},
		},
		{
			name: "fail on invalid kind",
			path: "/value/xxx/test",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "fail on empty metric name",
			path: "/value/counter",
			want: result{code: http.StatusNotFound},
		},
		{
			name: "fail on counter with invalid name",
			path: "/value/counter/inva!id",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "fail on gauge with invalid name",
			path: "/value/gauge/inval!d",
			want: result{code: http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, sm, _ := createTestRouter()

			if tt.mock != nil {
				tt.mock(sm)
			}

			code, _, body := testRequest(t, router, http.MethodGet, tt.path, nil)

			assert.Equal(t, tt.want.code, code)
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}

func TestGetMetricJSON(t *testing.T) {
	type result struct {
		code int
		body metrics.MetricExchange
	}

	tests := []struct {
		name     string
		mex      metrics.MetricExchange
		mock     func(m *storage.ServiceMock)
		expected result
	}{
		{
			name: "Should get counter",
			mex:  metrics.NewGetCounterMex("test"),
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindCounter).Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			expected: result{
				code: http.StatusOK,
				body: metrics.NewUpdateCounterMex("test", 42),
			},
		},
		{
			name: "Should get gauge",
			mex:  metrics.NewGetGaugeMex("test"),
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindGauge).Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)}, nil)
			},
			expected: result{
				code: http.StatusOK,
				body: metrics.NewUpdateGaugeMex("test", 42.42),
			},
		},
		{
			name: "Should fail on unknown metric kind",
			mex:  metrics.MetricExchange{ID: "test", MType: "unknown"},
			expected: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Should fail on unknown counter",
			mex:  metrics.NewGetCounterMex("test"),
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindCounter).Return(storage.Record{}, entities.ErrRecordNotFound)
			},
			expected: result{
				code: http.StatusNotFound,
			},
		},
		{
			name: "Should fail on unknown gauge",
			mex:  metrics.NewGetGaugeMex("test"),
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindGauge).Return(storage.Record{}, entities.ErrRecordNotFound)
			},
			expected: result{
				code: http.StatusNotFound,
			},
		},
		{
			name: "Should fail on counter with invalid name",
			mex:  metrics.NewGetCounterMex("X)"),
			expected: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Should fail on gauge with invalid name",
			mex:  metrics.NewGetGaugeMex("X;"),
			expected: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Should fail on broken service",
			mex:  metrics.NewGetGaugeMex("test"),
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", metrics.KindGauge).Return(storage.Record{}, entities.ErrUnexpected)
			},
			expected: result{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			router, sm, _ := createTestRouter()

			if tt.mock != nil {
				tt.mock(sm)
			}

			payload, err := json.Marshal(tt.mex)
			require.NoError(err)

			code, contentType, body := testRequest(t, router, http.MethodPost, "/value", payload)
			assert.Equal(tt.expected.code, code)

			if tt.expected.code == http.StatusOK {
				assert.Equal("application/json", contentType)

				var resp metrics.MetricExchange
				err = json.Unmarshal(body, &resp)
				require.NoError(err)

				assert.Equal(tt.expected.body, resp)
			}
		})
	}
}

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
			expected: result{
				code: http.StatusOK,
			},
		},
		{
			name:         "should return not implemented if storage doesn't support ping",
			pingResponse: entities.ErrStorageUnpingable,
			expected: result{
				code: http.StatusNotImplemented,
			},
		},
		{
			name:         "should return internal server error if storage offline",
			pingResponse: entities.ErrUnexpected,
			expected: result{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _, pm := createTestRouter()
			pm.On("Ping", mock.Anything).Return(tt.pingResponse)

			code, _, _ := testRequest(t, router, http.MethodGet, "/ping", nil)

			if tt.expected.code != code {
				t.Fatalf("expected response to be %d, got: %d", tt.expected.code, code)
			}
		})
	}
}
