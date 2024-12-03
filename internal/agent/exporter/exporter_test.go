package exporter

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/pkg/grpcapi"
	"github.com/ex0rcist/metflix/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestHTTPExporter(t *testing.T) {
	assert := assert.New(t)

	secSign := "mocked_sign"
	baseURL := entities.Address("localhost:8080")
	signer := mockSigner(secSign)

	wg := sync.WaitGroup{}
	wg.Add(1)

	server := newTestServer(t, baseURL.String(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal(r.Header.Get("HashSHA256"), secSign)

		w.WriteHeader(http.StatusOK)

		wg.Done()
	}))
	defer server.Close()

	exporter := NewHTTPExporter(context.Background(), &baseURL, signer, 1, nil)
	assert.NotNil(exporter)
	assert.Equal(baseURL, *exporter.baseURL)
	assert.Equal(signer, exporter.signer)

	// test sendd empty buffer
	err := exporter.Send()
	assert.Equal("cannot send empty buffer", err.Error())

	// test add()
	exporter.Add("test_counter", metrics.Counter(10))
	assert.Equal(1, len(exporter.buffer))
	assert.Equal("test_counter", exporter.buffer[0].ID)

	// test send()
	err = exporter.Send()
	assert.NoError(err)
	assert.NotNil(exporter.jobs)

	wg.Wait()

	assert.Equal(0, len(exporter.buffer))

	// test error()
	exporter.err = errors.New("test error")
	assert.Equal("metrics export failed: test error", exporter.Error().Error())

	// test reset()
	exporter.Reset()
	assert.Equal(0, len(exporter.buffer))
	assert.Nil(exporter.err)
}

func TestHTTPBatchExporter(t *testing.T) {
	assert := assert.New(t)

	secSign := "mocked_sign"
	baseURL := entities.Address("localhost:8080")
	signer := mockSigner(secSign)

	wg := sync.WaitGroup{}
	wg.Add(1)

	server := newTestServer(t, baseURL.String(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal(r.Header.Get("HashSHA256"), secSign)

		w.WriteHeader(http.StatusOK)

		wg.Done()
	}))
	defer server.Close()

	exporter := NewHTTPBatchExporter(context.Background(), &baseURL, signer, nil)
	assert.NotNil(exporter)
	assert.Equal(baseURL, *exporter.baseURL)
	assert.Equal(signer, exporter.signer)

	// test send empty buffer
	err := exporter.Send()
	assert.Equal("cannot send empty buffer", err.Error())

	// test add()
	exporter.Add("test_counter", metrics.Counter(10))
	assert.Equal(1, len(exporter.buffer))
	assert.Equal("test_counter", exporter.buffer[0].ID)

	// test send()
	err = exporter.Send()
	assert.NoError(err)

	wg.Wait()

	assert.Equal(0, len(exporter.buffer))

	// test error()
	exporter.err = errors.New("test error")
	assert.Equal("metrics export failed: test error", exporter.Error().Error())

	// test reset()
	exporter.Reset()
	assert.Equal(0, len(exporter.buffer))
	assert.Nil(exporter.err)
}

func TestGRPCExporter(t *testing.T) {
	assert := assert.New(t)

	baseURL := entities.Address("localhost:50051")

	wg := &sync.WaitGroup{}
	wg.Add(1)

	cancel := newGRPCTestServer(t, baseURL.String(), wg)
	defer cancel()

	exporter := NewGRPCExporter(&baseURL, nil)
	assert.NotNil(exporter)
	assert.Equal(baseURL, *exporter.baseURL)

	// test send empty buffer
	err := exporter.Send()
	assert.Equal("cannot send empty buffer", err.Error())

	// test add()
	exporter.Add("test_counter", metrics.Counter(10))
	assert.Equal(1, len(exporter.buffer))
	assert.Equal("test_counter", exporter.buffer[0].Id)

	// test send()
	err = exporter.Send()
	assert.NoError(err)

	assert.Equal(0, len(exporter.buffer))

	// test error()
	exporter.err = errors.New("test error")
	assert.Equal("metrics export failed: test error", exporter.Error().Error())

	// test reset()
	exporter.Reset()
	assert.Equal(0, len(exporter.buffer))
	assert.Nil(exporter.err)
}

func mockSigner(signature string) security.Signer {
	signer := new(security.MockSigner)
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

type TestMetricsServer struct {
	grpcapi.UnimplementedMetricsServer

	wg *sync.WaitGroup
}

func (s *TestMetricsServer) BatchUpdate(ctx context.Context, req *grpcapi.BatchUpdateRequest) (*grpcapi.BatchUpdateResponse, error) {
	s.wg.Done()

	return &grpcapi.BatchUpdateResponse{}, nil
}

func newGRPCTestServer(t *testing.T, bind string, wg *sync.WaitGroup) func() {
	server := grpc.NewServer()
	grpcapi.RegisterMetricsServer(server, &TestMetricsServer{wg: wg})

	lis, err := net.Listen("tcp", bind)
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		err := server.Serve(lis)
		if err != nil {
			t.Logf("Failed to serve: %v", err)
		}

		wg.Wait()
	}()

	return func() {
		server.Stop()
		_ = lis.Close()
	}
}
