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
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestHomepage(t *testing.T) {
	type result struct {
		code int
		body string
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
			want: result{code: http.StatusOK, body: "mainpage here.\nmetrics list:\ntest1 => counter: 1\ntest2 => gauge: 2.3\n"},
		},
	}

	for _, tt := range tests {
		storage := storage.NewMemStorage()
		router := server.NewRouter(storage)
		testServer := httptest.NewServer(router)
		defer testServer.Close()

		if len(tt.metrics) > 0 {
			for _, r := range tt.metrics {
				storage.Push(r)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			req, body := testRequest(t, testServer, http.MethodGet, tt.path, nil)

			require.Equal(tt.want.code, req.StatusCode)
			require.Equal(tt.want.body, string(body))
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

			req, body := testRequest(t, testServer, http.MethodPost, tt.path, nil)

			require.Equal(tt.want.code, req.StatusCode)
			require.Equal(tt.want.body, string(body))
		})
	}
}
