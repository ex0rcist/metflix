package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/utils"
)

func TestDecryptRequest(t *testing.T) {
	prvKey, pubKey := mockKeys()

	tt := []struct {
		name           string
		pubKey         security.PublicKey
		prvKey         security.PrivateKey
		message        []byte
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "no keys configured",
			message:        []byte("test message"),
			expectedStatus: http.StatusOK,
			expectedBody:   "test message",
		},
		{
			name:           "prv key configured, message encoded",
			pubKey:         pubKey,
			prvKey:         prvKey,
			message:        encrypt([]byte("test message 2"), pubKey),
			expectedStatus: http.StatusOK,
			expectedBody:   "test message 2",
		},
		{
			name:           "prv key configured, message not (correctly) encoded",
			pubKey:         pubKey,
			prvKey:         prvKey,
			message:        []byte("test message 3"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		handler := createDecryptMiddleware(tc.prvKey)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tc.message))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		result := rec.Result()
		defer func() {
			_ = result.Body.Close()
		}()

		if tc.expectedStatus != result.StatusCode {
			t.Errorf("expected response status to be %d, got %d", tc.expectedStatus, result.StatusCode)
		}

		if tc.expectedStatus == http.StatusOK {
			if rec.Body.String() != tc.expectedBody {
				t.Errorf("expected response body to be %s, got '%s'", tc.expectedBody, rec.Body.String())
			}
		}
	}
}

func mockKeys() (security.PrivateKey, security.PublicKey) {
	root, _ := utils.GetProjectRoot()

	prvKey, _ := security.NewPrivateKey(entities.FilePath(filepath.Join(root, "example_key.pem")))
	pubKey, _ := security.NewPublicKey(entities.FilePath(filepath.Join(root, "example_key.pub.pem")))

	return prvKey, pubKey
}

func createDecryptMiddleware(key security.PrivateKey) http.Handler {
	return DecryptRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := new(bytes.Buffer)
		_, rErr := body.ReadFrom(r.Body)
		if rErr != nil {
			panic(rErr)
		}

		_, wErr := w.Write(body.Bytes())
		if wErr != nil {
			panic(wErr)
		}
	}), key)
}

func encrypt(message []byte, key security.PublicKey) []byte {
	buff, _ := security.Encrypt(bytes.NewReader(message), key)
	return buff.Bytes()
}
