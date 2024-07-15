package logging

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func RequestsLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := findOrCreateRequestID(r)
		start := time.Now()

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

		// execute
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ctx := logger.WithContext(r.Context())
		next.ServeHTTP(ww, r.WithContext(ctx))

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
		requestID = uuid.NewV4().String()
	}

	return requestID
}
