package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/ex0rcist/metflix/internal/compression"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

func RequestsLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := findOrCreateRequestID(r)

		// setup child logger for middleware
		logger := log.Logger.With().
			Str("rid", requestID).
			Logger()

		// log started
		logger.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote-addr", r.RemoteAddr). // middleware.RealIP
			Msg("Started")

		logger.Debug().
			Msgf("request: %s", utils.HeadersToStr(r.Header))

		// TODO: context?

		// execute
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ctx := logger.WithContext(r.Context())
		next.ServeHTTP(ww, r.WithContext(ctx))

		logger.Debug().
			Msgf("response: %s", utils.HeadersToStr(ww.Header()))

		// log completed
		logger.Info().
			Float64("elapsed", time.Since(start).Seconds()).
			Int("status", ww.Status()).
			Int("size", ww.BytesWritten()).
			Msg("Completed")
	})
}

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

func findOrCreateRequestID(r *http.Request) string {
	requestID := r.Header.Get("X-Request-Id")

	if requestID == "" {
		requestID = utils.GenerateRequestID()
	}

	return requestID
}

func needGzipEncoding(r *http.Request) bool {
	if len(r.Header.Get("Accept-Encoding")) == 0 {
		return false
	}

	for _, encoding := range r.Header.Values("Accept-Encoding") {
		if encoding == "gzip" {
			return true
		}
	}

	return false
}
