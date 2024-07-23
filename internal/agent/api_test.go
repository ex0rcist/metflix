package agent

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ex0rcist/metflix/internal/compression"
	"github.com/ex0rcist/metflix/internal/entities"
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

	require.NotPanics(func() {
		address := entities.Address("localhost")
		NewAPI(&address, nil)
	})
}

func TestApiClientReport(t *testing.T) {
	rtf := func(req *http.Request) *http.Response {
		assert.Equal(t, "http://localhost:8080/update", req.URL.String())
		assert.Equal(t, http.MethodPost, req.Method)

		payload, err := json.Marshal(metrics.NewUpdateCounterMex("test", 42))
		require.NoError(t, err)

		expectedPayload, err := compression.Pack(payload)
		require.NoError(t, err)

		actualPayload, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		assert.Equal(t, expectedPayload.Bytes(), actualPayload)

		return &http.Response{
			StatusCode: 200,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}
	}

	address := entities.Address("localhost:8080")

	api := NewAPI(&address, RoundTripFunc(rtf))
	api.Report("test", metrics.Counter(42))
}
