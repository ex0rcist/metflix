package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/server"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/stretchr/testify/require"
)

// https://github.com/go-chi/chi/blob/cca4135d8dddff765463feaf1118047a9e506b4a/middleware/get_head_test.go
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (int, string, []byte) {
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
		storage := storage.NewMemStorage()
		router := server.NewRouter(storage)
		testServer := httptest.NewServer(router)
		defer testServer.Close()

		if len(tt.metrics) > 0 {
			for _, r := range tt.metrics {
				err := storage.Push(r)
				require.NoError(err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			code, _, body := testRequest(t, testServer, http.MethodGet, tt.path, nil)

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
		want result
	}{
		{
			name: "push counter",
			path: "/update/counter/test/42",
			want: result{code: http.StatusOK, body: ""},
		},
		{
			name: "push gauge",
			path: "/update/gauge/test/42.42",
			want: result{code: http.StatusOK, body: ""},
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
		t.Run(tt.name, func(t *testing.T) {
			router := server.NewRouter(storage.NewMemStorage())
			testServer := httptest.NewServer(router)
			defer testServer.Close()

			code, _, body := testRequest(t, testServer, http.MethodPost, tt.path, nil)

			require.Equal(tt.want.code, code)
			require.Equal(tt.want.body, string(body))
		})
	}
}
