package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/go-chi/chi/middleware"
)

// CustomResponseWriter is a wrapper around http.ResponseWriter that captures the response body
type CustomResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func SignResponse(next http.Handler, secret entities.Secret) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if len(secret) == 0 { // skip middleware entirely
			next.ServeHTTP(w, r)
			return
		}

		// Wrap the ResponseWriter with chi's middleware.WrapResponseWriter
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Create a buffer to capture the response body
		bodyBuffer := &bytes.Buffer{}

		// Create a custom ResponseWriter to capture the response body
		crw := &CustomResponseWriter{ResponseWriter: ww, body: bodyBuffer}

		// Pass the custom ResponseWriter to the next handler
		next.ServeHTTP(crw, r)

		signer := services.NewSignerService(secret)
		signature, _ := signer.CalculateSignature(bodyBuffer.Bytes())

		w.Header().Set("HashSHA256", signature)

		// Write the captured body to the original ResponseWriter
		_, err := w.Write(bodyBuffer.Bytes())
		if err != nil {
			logging.LogErrorCtx(ctx, fmt.Errorf("got empty signature for request"))
			return
		}
	})
}

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
			logging.LogErrorCtx(ctx, fmt.Errorf("got empty signature for request"))
			http.Error(w, "failed to verify signature", http.StatusBadRequest)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		r.Body.Close() //  must close
		if err != nil {
			logging.LogErrorCtx(ctx, fmt.Errorf("failed to read request body"))
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}

		signer := services.NewSignerService(secret)
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
