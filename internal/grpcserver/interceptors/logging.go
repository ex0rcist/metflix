package interceptors

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryRequestsInterceptor is grpc unary interceptor which logs incoming requests and responses.
func UnaryRequestsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	// setup child logger for middleware
	logger := log.Logger.With().
		Str("transport", "gRPC").
		Logger()

	ctx = logger.WithContext(ctx)
	requestID, clientIP := extractMetaData(ctx)

	// log started
	logger.Info().
		Str("rid", requestID).
		Str("method", info.FullMethod).
		Str("remote-addr", clientIP).
		Msg("Started")

	// execute
	resp, err := handler(ctx, req)
	status, _ := status.FromError(err)

	// log completed
	logger.Info().
		Float64("elapsed", time.Since(start).Seconds()).
		Str("status", status.Code().String()).
		Msg("Completed")

	return resp, err
}
