package grpcserver

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/pkg/grpcapi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name      string
		checkResp error
		expected  codes.Code
	}{
		{
			name:      "Should return OK, if storage online",
			checkResp: nil,
			expected:  codes.OK,
		},
		{
			name:      "Should return not implemented, if storage doesn't support health check",
			checkResp: entities.ErrStorageUnpingable,
			expected:  codes.Unimplemented,
		},
		{
			name:      "Should return internal error, if storage offline",
			checkResp: entities.ErrUnexpected,
			expected:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(services.HealthCheckServiceMock)
			m.On("Ping", mock.Anything).Return(tt.checkResp)

			conn, closer := createTestServer(t, nil, m, nil)
			defer closer()

			client := grpcapi.NewHealthClient(conn)
			_, err := client.Ping(context.Background(), new(emptypb.Empty))

			rv, ok := status.FromError(err)

			require.True(t, ok)
			require.Equal(t, tt.expected, rv.Code())
		})
	}
}
