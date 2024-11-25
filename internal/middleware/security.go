package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/go-chi/chi/middleware"
)

// CustomResponseWriter is a wrapper around http.ResponseWriter that captures the response body
type CustomResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

// Write body
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// Sign response middleware
func SignResponse(next http.Handler, secret entities.Secret) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if len(secret) == 0 { // skip middleware entirely
			next.ServeHTTP(w, r)
			return
		}

		// wrap the ResponseWriter with chi's middleware.WrapResponseWriter
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// create a buffer to capture the response body
		bodyBuffer := &bytes.Buffer{}

		// create a custom ResponseWriter to capture the response body
		crw := &CustomResponseWriter{ResponseWriter: ww, body: bodyBuffer}

		// pass the custom ResponseWriter to the next handler
		next.ServeHTTP(crw, r)

		signer := security.NewSignerService(secret)
		signature, _ := signer.CalculateSignature(bodyBuffer.Bytes())

		w.Header().Set("HashSHA256", signature)

		// write the captured body to the original ResponseWriter
		_, err := w.Write(bodyBuffer.Bytes())
		if err != nil {
			logging.LogErrorCtx(ctx, fmt.Errorf("got empty signature for request"))
			return
		}
	})
}

// Ensure incoming request satisfies it's signature.
func CheckSignedRequest(next http.Handler, secret entities.Secret) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if len(secret) == 0 { // skip middleware entirely
			next.ServeHTTP(w, r)
			return
		}

		protected := map[string]struct{}{"POST": {}, "PUT": {}, "PATCH": {}}
		if _, ok := protected[r.Method]; !ok {
			logging.LogDebugCtx(ctx, "no need to check sign for that method")
			next.ServeHTTP(w, r)
			return
		}

		hash := r.Header.Get("HashSHA256")
		if len(hash) == 0 {
			// just pass it through for backward compatibility
			next.ServeHTTP(w, r)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		defer func() {
			if closeErr := r.Body.Close(); closeErr != nil {
				logging.LogError(closeErr)
			}
		}()

		if err != nil {
			logging.LogErrorCtx(ctx, fmt.Errorf("failed to read request body"))
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}

		signer := security.NewSignerService(secret)
		ok, _ := signer.VerifySignature(bodyBytes, hash)
		if !ok {
			logging.LogErrorCtx(ctx, fmt.Errorf("failed to verify request signature"))
			http.Error(w, "Failed to verify signature", http.StatusBadRequest)
			return
		}

		logging.LogDebugCtx(ctx, "got correct signature")

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		next.ServeHTTP(w, r)
	})
}

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

// Ensure incoming request is from a trsuted subnet
func FilterUntrustedRequest(next http.Handler, trustedSubnet *net.IPNet) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if trustedSubnet == nil { // skip middleware entirely
			next.ServeHTTP(w, r)
			return
		}

		clientIP := net.ParseIP(r.Header.Get("X-Real-IP"))

		if !trustedSubnet.Contains(clientIP) {
			logging.LogError(entities.UntrustedSubnetError(clientIP))
			http.Error(w, "", http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}
