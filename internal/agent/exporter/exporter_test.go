package exporter

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSigner struct {
	mock.Mock
}

func (m *MockSigner) CalculateSignature(data []byte) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockSigner) VerifySignature(data []byte, hash string) (bool, error) {
	args := m.Called(data, hash)
	return args.Bool(0), args.Error(1)
}

func mockSigner(signature string) security.Signer {
	signer := new(MockSigner)
	signer.On("CalculateSignature", mock.Anything).Return(signature, nil)

	return signer
}

func newTestServer(t *testing.T, bind string, handler http.HandlerFunc) *httptest.Server {
	listener, err := net.Listen("tcp", bind)
	require.NoError(t, err)

	server := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: handler},
	}

	server.Start()
	return server
}

func TestNewLimitedExporter(t *testing.T) {
	baseURL := entities.Address("localhost:8080")

	signer := mockSigner("mocked_sign")

	exporter := NewLimitedExporter(&baseURL, signer, 3, nil)

	assert.NotNil(t, exporter)
	assert.Equal(t, baseURL, *exporter.baseURL)
	assert.Equal(t, signer, exporter.signer)
	assert.NotNil(t, exporter.jobs)
}

func TestAdd(t *testing.T) {
	baseURL := entities.Address("localhost:8080")

	signer := mockSigner("mocked_sign")

	exporter := NewLimitedExporter(&baseURL, signer, 3, nil)

	counter := metrics.Counter(10)
	exporter.Add("test_counter", counter)

	assert.Equal(t, 1, len(exporter.buffer))
	assert.Equal(t, "test_counter", exporter.buffer[0].ID)
}

func TestLimitedSend(t *testing.T) {
	signer := mockSigner("mocked_sign")
	baseURL := entities.Address("127.0.0.1:8080")

	wg := sync.WaitGroup{}
	wg.Add(1)

	server := newTestServer(t, "127.0.0.1:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, r.Header.Get("HashSHA256"), "mocked_sign")

		w.WriteHeader(http.StatusOK)

		wg.Done()
	}))
	defer server.Close()

	exporter := NewLimitedExporter(&baseURL, signer, 1, nil)
	exporter.Add("test_counter", metrics.Counter(10))

	err := exporter.Send()
	assert.NoError(t, err)

	wg.Wait()
}

func TestBatchSend(t *testing.T) {
	signer := mockSigner("mocked_sign")
	baseURL := entities.Address("127.0.0.1:8080")

	wg := sync.WaitGroup{}
	wg.Add(1)

	server := newTestServer(t, "127.0.0.1:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, r.Header.Get("HashSHA256"), "mocked_sign")

		w.WriteHeader(http.StatusOK)

		wg.Done()
	}))
	defer server.Close()

	exporter := NewBatchExporter(&baseURL, signer, nil)
	exporter.Add("test_counter", metrics.Counter(10))

	err := exporter.Send()
	assert.NoError(t, err)

	wg.Wait()

	assert.Equal(t, len(exporter.buffer), 0)
}

func TestSendEmptyBuffer(t *testing.T) {
	signer := mockSigner("mocked_sign")
	baseURL := entities.Address("localhost:8080")

	exporter := NewLimitedExporter(&baseURL, signer, 1, nil)

	err := exporter.Send()

	assert.Error(t, err)
	assert.Equal(t, "cannot send empty buffer", err.Error())
}

func TestError(t *testing.T) {
	signer := mockSigner("mocked_sign")
	baseURL := entities.Address("localhost:8080")

	exporter := NewLimitedExporter(&baseURL, signer, 1, nil)

	exporter.err = errors.New("test error")
	assert.Equal(t, "metrics export failed: test error", exporter.Error().Error())
}

func TestReset(t *testing.T) {
	signer := mockSigner("mocked_sign")
	baseURL := entities.Address("localhost:8080")

	exporter := NewLimitedExporter(&baseURL, signer, 1, nil)

	exporter.Add("test_counter", metrics.Counter(10))
	exporter.Reset()

	assert.Equal(t, 0, len(exporter.buffer))
	assert.Nil(t, exporter.err)
}

func TestWorker(t *testing.T) {
	signer := mockSigner("mocked_sign")
	baseURL := entities.Address("localhost:8080")

	exporter := NewLimitedExporter(&baseURL, signer, 1, nil)

	exporter.Add("test_counter", metrics.Counter(10))
	err := exporter.Send()

	assert.NoError(t, err)

	time.Sleep(1 * time.Second) // wait for worker to process the job

	assert.Equal(t, 0, len(exporter.buffer))
}
