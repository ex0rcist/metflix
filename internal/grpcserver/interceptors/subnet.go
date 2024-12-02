package interceptors

import (
	"context"
	"net"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Interceptor to match if trusted subnet used
func UnaryRequestsFilter(trustedSubnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if trustedSubnet == nil {
			return handler(ctx, req)
		}

		_, rawIP := extractMetaData(ctx)
		clientIP := net.ParseIP(rawIP)

		if !trustedSubnet.Contains(clientIP) {
			err := entities.ErrUntrustedSubnet
			logging.LogErrorCtx(ctx, err)

			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		return handler(ctx, req)
	}
}
