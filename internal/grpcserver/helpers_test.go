package grpcserver

import (
	"context"
	"net"
	"testing"

	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func createTestServer(
	t *testing.T,
	metrics *services.MetricServiceMock,
	healthcheck *services.HealthCheckServiceMock,
	prvKey security.PrivateKey,
) (*grpc.ClientConn, func()) {
	t.Helper()
	require := require.New(t)

	if metrics == nil {
		metrics = &services.MetricServiceMock{}
	}

	if healthcheck == nil {
		healthcheck = &services.HealthCheckServiceMock{}
	}

	lis := bufconn.Listen(1024 * 1024)
	srv := NewBackend(
		WithHealthService(healthcheck),
		WithMetricService(metrics),
		WithPrivateKey(prvKey),
	)

	go func() {
		require.NoError(srv.Serve(lis))
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("0.0.0.0", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(err)

	closer := func() {
		require.NoError(conn.Close())
		srv.Stop()
		require.NoError(lis.Close())
	}

	return conn, closer
}
