package agent_test

import (
	"net/http"
	"testing"

	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// http://hassansin.github.io/Unit-Testing-http-client-in-Go
type RoundTripFunc func(req *http.Request) *http.Response

// todo: wtf...
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestNewApi(t *testing.T) {
	require := require.New(t)
	require.NotPanics(func() { agent.NewAPI("someaddr", nil) })
}

func TestApiClientReport(t *testing.T) {
	rtf := func(req *http.Request) *http.Response {
		assert.Equal(t, "http://localhost:8080/update/counter/Test/0", req.URL.String())
		assert.Equal(t, http.MethodPost, req.Method)

		return &http.Response{
			StatusCode: 200,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}
	}

	api := agent.NewAPI("http://localhost:8080", RoundTripFunc(rtf))
	api.Report("Test", metrics.Counter(0))
}
