package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSignerService struct {
	mock.Mock
}

func (m *MockSignerService) CalculateSignature(data []byte) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockSignerService) VerifySignature(data []byte, signature string) (bool, error) {
	args := m.Called(data, signature)
	return args.Bool(0), args.Error(1)
}

func TestSignResponseMiddleware(t *testing.T) {
	secret := entities.Secret("my-secret-key")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("test response"))
		require.NoError(t, err)
	})

	signedHandler := SignResponse(handler, secret)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr := httptest.NewRecorder()

	signedHandler.ServeHTTP(rr, req)

	result := rr.Result()
	defer func() {
		err := result.Body.Close()
		if err != nil {
			logging.LogError(err)
		}
	}()

	hash := result.Header.Get("HashSHA256")
	assert.NotEmpty(t, hash)

	body, _ := io.ReadAll(result.Body)
	assert.Equal(t, "test response", string(body))
}

func TestSignResponseMiddlewareWithoutSecret(t *testing.T) {
	secret := entities.Secret("")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("test response"))
		require.NoError(t, err)
	})

	signedHandler := SignResponse(handler, secret)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr := httptest.NewRecorder()

	signedHandler.ServeHTTP(rr, req)

	result := rr.Result()
	defer func() {
		err := result.Body.Close()
		if err != nil {
			logging.LogError(err)
		}
	}()

	hash := result.Header.Get("HashSHA256")
	assert.Empty(t, hash)

	body, _ := io.ReadAll(result.Body)
	assert.Equal(t, "test response", string(body))
}

func TestCheckSignedRequestMiddleware(t *testing.T) {
	secret := entities.Secret("my-secret-key")
	signer := services.NewSignerService(secret)
	body := []byte("test request")

	signature, err := signer.CalculateSignature(body)
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		require.NoError(t, err)
	})

	checkSignedHandler := CheckSignedRequest(handler, secret)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(body))
	req.Header.Set("HashSHA256", signature)

	rr := httptest.NewRecorder()
	checkSignedHandler.ServeHTTP(rr, req)

	result := rr.Result()
	defer func() {
		err := result.Body.Close()
		if err != nil {
			logging.LogError(err)
		}
	}()

	assert.Equal(t, http.StatusOK, result.StatusCode)

	respBody, _ := io.ReadAll(result.Body)
	assert.Equal(t, "ok", string(respBody))
}

func TestCheckSignedRequestMiddlewareInvalidSignature(t *testing.T) {
	secret := entities.Secret("my-secret-key")
	body := []byte("test request")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		require.NoError(t, err)
	})

	checkSignedHandler := CheckSignedRequest(handler, secret)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(body))
	req.Header.Set("HashSHA256", "invalid-signature")

	rr := httptest.NewRecorder()
	checkSignedHandler.ServeHTTP(rr, req)

	result := rr.Result()
	defer func() {
		err := result.Body.Close()
		if err != nil {
			logging.LogError(err)
		}
	}()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)

	respBody, _ := io.ReadAll(result.Body)
	assert.Equal(t, "Failed to verify signature\n", string(respBody))
}

func TestCheckSignedRequestMiddlewareWithoutSecret(t *testing.T) {
	secret := entities.Secret("")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		require.NoError(t, err)
	})

	checkSignedHandler := CheckSignedRequest(handler, secret)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", nil)

	rr := httptest.NewRecorder()
	checkSignedHandler.ServeHTTP(rr, req)

	result := rr.Result()
	defer func() {
		err := result.Body.Close()
		if err != nil {
			logging.LogError(err)
		}
	}()

	assert.Equal(t, http.StatusOK, result.StatusCode)

	respBody, _ := io.ReadAll(result.Body)
	assert.Equal(t, "ok", string(respBody))
}
