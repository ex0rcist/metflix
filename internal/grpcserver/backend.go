package grpcserver

import (
	"net"

	"github.com/ex0rcist/metflix/internal/grpcserver/interceptors"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/services"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

type Backend struct {
	privateKey    security.PrivateKey
	trustedSubnet *net.IPNet

	server *grpc.Server

	healthService services.HealthChecker
	metricService services.MetricProvider
}

// Backend constructor
func NewBackend(opts ...Option) *grpc.Server {
	backend := &Backend{}

	for _, opt := range opts {
		opt(backend)
	}

	backend.setup()

	return backend.server
}

func (b *Backend) setup() {
	icep := b.prepareInterceptors()

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(icep...),
	}

	grpcServer := grpc.NewServer(grpcOpts...)

	b.server = grpcServer

	RegisterhHealthServer(grpcServer, b.healthService)
	RegisterMetricsServer(grpcServer, b.metricService, b.privateKey)
}

func (b *Backend) prepareInterceptors() []grpc.UnaryServerInterceptor {
	iceps := make([]grpc.UnaryServerInterceptor, 0, 2)
	iceps = append(iceps, interceptors.UnaryRequestsInterceptor)
	iceps = append(iceps, interceptors.UnaryRequestsFilter(b.trustedSubnet))

	return iceps
}

/* Options */

type Option func(*Backend)

func WithPrivateKey(privateKey security.PrivateKey) Option {
	return func(b *Backend) {
		b.privateKey = privateKey
	}
}

func WithTrustedSubnet(trustedSubnet *net.IPNet) Option {
	return func(b *Backend) {
		b.trustedSubnet = trustedSubnet
	}
}

func WithHealthService(healthService services.HealthChecker) Option {
	return func(b *Backend) {
		b.healthService = healthService
	}
}

func WithMetricService(metricService services.MetricProvider) Option {
	return func(b *Backend) {
		b.metricService = metricService
	}
}
