package grpcserver

import (
	"context"
	"errors"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/pkg/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HealthServer verifies current health status of the service.
type HealthServer struct {
	grpcapi.UnimplementedHealthServer
	healthService services.HealthChecker
}

// RegisterhHealthServer creates new instance of gRPC serving Health API and attaches it to the server.
func RegisterhHealthServer(server *grpc.Server, healthService services.HealthChecker) {
	s := &HealthServer{healthService: healthService}

	grpcapi.RegisterHealthServer(server, s)
}

// Ping verifies connection to the database.
func (s HealthServer) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.healthService.Ping(ctx)
	if err == nil {
		return new(emptypb.Empty), nil
	}

	if errors.Is(err, entities.ErrStorageUnpingable) {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return nil, status.Error(codes.Internal, err.Error())
}
