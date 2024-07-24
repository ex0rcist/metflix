package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/storage"
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

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, contentType, respBody
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
		m := storage.ServiceMock{}
		m.On("List").Return(tt.metrics, nil)

		router := NewRouter(&m)

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
			name: "push counter",
			path: "/update/counter/test/42",
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", "counter").Return(storage.Record{}, nil)
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42"},
		},
		{
			name: "push counter with existing value",
			path: "/update/counter/test/42",
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", "counter").Return(storage.Record{Name: "test", Value: metrics.Counter(21)}, nil)
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42"},
		},
		{
			name: "push gauge",
			path: "/update/gauge/test/42.42",
			mock: func(m *storage.ServiceMock) {
				m.On("Push", mock.AnythingOfType("Record")).Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42.42"},
		},
		{
			name: "fail on invalid kind",
			path: "/update/xxx/test/1",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "fail on empty metric name",
			path: "/update/counter//1",
			want: result{code: http.StatusNotFound},
		},
		{
			name: "fail on counter with invalid name",
			path: "/update/counter/inva!id/10",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "fail on counter with invalid value",
			path: "/update/counter/test/10.0",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "fail on gauge with invalid name",
			path: "/update/gauge/inval!d/42.42",
			want: result{code: http.StatusBadRequest},
		},
		{
			name: "fail on gauge with invalid value",
			path: "/update/gauge/test/42.42!",
			want: result{code: http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		m := new(storage.ServiceMock)

		if tt.mock != nil {
			tt.mock(m)
		}

		router := NewRouter(m)

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
				m.On("Get", "test", "counter").Return(storage.Record{}, nil)
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

			m := storage.ServiceMock{}

			if tt.mock != nil {
				tt.mock(&m)
			}

			router := NewRouter(&m)

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
				m.On("Get", "test", "counter").Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			want: result{code: http.StatusOK, body: "42"},
		},
		{
			name: "get gauge",
			path: "/value/gauge/test",
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", "gauge").Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)}, nil)
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
			m := storage.ServiceMock{}
			router := NewRouter(&m)

			if tt.mock != nil {
				tt.mock(&m)
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
				m.On("Get", "test", "counter").Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
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
				m.On("Get", "test", "gauge").Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)}, nil)
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
				m.On("Get", "test", "counter").Return(storage.Record{}, entities.ErrRecordNotFound)
			},
			expected: result{
				code: http.StatusNotFound,
			},
		},
		{
			name: "Should fail on unknown gauge",
			mex:  metrics.NewGetGaugeMex("test"),
			mock: func(m *storage.ServiceMock) {
				m.On("Get", "test", "gauge").Return(storage.Record{}, entities.ErrRecordNotFound)
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
				m.On("Get", "test", "gauge").Return(storage.Record{}, entities.ErrUnexpected)
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

			m := storage.ServiceMock{}

			if tt.mock != nil {
				tt.mock(&m)
			}

			router := NewRouter(&m)

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
