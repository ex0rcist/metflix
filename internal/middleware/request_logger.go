package middleware

import (
	"net/http"
	"time"

	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

// Log requests middleware
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

func findOrCreateRequestID(r *http.Request) string {
	requestID := r.Header.Get("X-Request-Id")

	if requestID == "" {
		requestID = utils.GenerateRequestID()
	}

	return requestID
}
