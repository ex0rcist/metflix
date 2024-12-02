package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ex0rcist/metflix/internal/compression"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
)

// Decompress request if possible
func DecompressRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		encoding := r.Header.Get("Content-Encoding")
		if len(encoding) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		decompressor := compression.NewDecompressor(r, ctx)
		defer decompressor.Close()

		err := decompressor.Decompress()
		if err != nil {
			switch {
			case errors.Is(err, entities.ErrEncodingUnsupported):
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, entities.ErrEncodingInternal):
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Compress response if needed
func CompressResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !needGzipEncoding(r) {
			logging.LogDebugCtx(ctx, "compression not requested or not supported by client")

			next.ServeHTTP(w, r)
			return
		}

		compressor := compression.NewCompressor(w, ctx)
		defer compressor.Close()

		next.ServeHTTP(compressor, r)
	})
}

func needGzipEncoding(r *http.Request) bool {
	if len(r.Header.Get("Accept-Encoding")) == 0 {
		return false
	}

	for _, encoding := range r.Header.Values("Accept-Encoding") {
		if strings.Contains(encoding, "gzip") {
			return true
		}
	}

	return false
}
