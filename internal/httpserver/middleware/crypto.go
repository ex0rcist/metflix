package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
)

// Decrypt request body using RSA.
func DecryptRequest(next http.Handler, key security.PrivateKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if key == nil { // skip middleware entirely
			next.ServeHTTP(w, r)
			return
		}

		msg, err := security.Decrypt(r.Body, key)
		if err != nil {
			logging.LogError(err, "error decoding request")
			http.Error(w, "decrypt failed", http.StatusBadRequest)

			return
		}

		r.Body = io.NopCloser(bytes.NewReader(msg.Bytes()))
		next.ServeHTTP(w, r)
	})
}
