package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ex0rcist/metflix/internal/grpcserver"
	"github.com/ex0rcist/metflix/internal/httpserver"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/storage"
)

const shutdownTimeout = 60 * time.Second

// Backend heart
type Server struct {
	config         *Config
	httpServer     *HTTPServer
	grpcServer     *GRPCServer
	profilerServer *ProfilerServer
	storage        storage.MetricsStorage
	privateKey     security.PrivateKey
}

// Server constructor
func New() (*Server, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	dataStorage, err := setupStorage(config)
	if err != nil {
		return nil, err
	}

	privateKey, err := preparePrivateKey(config)
	if err != nil {
		return nil, err
	}

	metricService := services.NewMetricService(dataStorage)
	healthService := services.NewHealthCheckService(dataStorage)

	httpServer := setupHTTPServer(config, metricService, healthService, privateKey)
	grpcServer := setupGRPCServer(config, metricService, healthService, privateKey)
	profilerServer := setupProfilerServer(config)

	return &Server{
		config:         config,
		httpServer:     httpServer,
		grpcServer:     grpcServer,
		profilerServer: profilerServer,
		storage:        dataStorage,
		privateKey:     privateKey,
	}, nil
}

// Start all subservices
func (s *Server) Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s.httpServer.Start()
	s.grpcServer.Start()
	s.profilerServer.Start()

	logging.LogInfo(s.String())
	logging.LogInfo("server ready")

	select {
	case s := <-interrupt:
		logging.LogInfo("interrupt: signal " + s.String())
	case err := <-s.httpServer.Notify():
		logging.LogError(err, "Server -> Start() -> s.httpServer.Notify")
	case err := <-s.grpcServer.Notify():
		logging.LogError(err, "Server -> Start() -> s.grpcServer.Notify")
	case err := <-s.profilerServer.Notify():
		logging.LogError(err, "Server -> Start() - s.profilerServer.Notify")
	}

	logging.LogInfo("shutting down...")

	stopped := make(chan struct{})
	stopCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	go func() {
		s.shutdown(stopCtx)
		close(stopped)
	}()

	select {
	case <-stopped:
		logging.LogInfo("server shutdown successful")

	case <-stopCtx.Done():
		logging.LogInfo("shutdown timeout exceeded")
	}
}

// Stringer for logging
func (s *Server) String() string {
	str := []string{
		fmt.Sprintf("address=%s", s.config.Address),
	}

	if stringer, ok := s.storage.(fmt.Stringer); ok {
		str = append(str, stringer.String())
	}

	if len(s.config.Secret) > 0 {
		str = append(str, fmt.Sprintf("secret=%s", s.config.Secret))
	}

	if len(s.config.PrivateKeyPath) > 0 {
		str = append(str, fmt.Sprintf("private-key=%v", s.config.PrivateKeyPath))
	}

	if s.config.TrustedSubnet != nil {
		str = append(str, fmt.Sprintf("trusted-subnet=%v", s.config.TrustedSubnet.String()))
	}

	return "server config: " + strings.Join(str, "; ")
}

func (s *Server) shutdown(ctx context.Context) {
	logging.LogInfo("shutting down HTTP API")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logging.LogError(err)
	}

	logging.LogInfo("shutting down gRPC API")
	s.grpcServer.Shutdown()

	logging.LogInfo("shutting down storage")
	if err := s.storage.Close(ctx); err != nil {
		logging.LogError(err)
	}

	logging.LogInfo("shutting down profiler")
	if err := s.profilerServer.Shutdown(ctx); err != nil {
		logging.LogError(err)
	}
}

func setupHTTPServer(
	config *Config,
	metricService services.MetricProvider,
	healthService services.HealthChecker,
	privateKey security.PrivateKey,
) *HTTPServer {
	healthResource := httpserver.NewHealthResource(healthService)
	metricResource := httpserver.NewMetricResource(metricService)

	handler := httpserver.NewBackend(
		httpserver.WithTrustedSubnet(config.TrustedSubnet),
		httpserver.WithSignSecret(config.Secret),
		httpserver.WithPrivateKey(privateKey),

		httpserver.WithHealthResource(healthResource),
		httpserver.WithMetricResource(metricResource),
	)

	return NewHTTPServer(handler, config.Address)
}

func setupGRPCServer(
	config *Config,
	metricService services.MetricProvider,
	healthService services.HealthChecker,
	privateKey security.PrivateKey,
) *GRPCServer {
	srv := grpcserver.NewBackend(
		grpcserver.WithTrustedSubnet(config.TrustedSubnet),
		grpcserver.WithPrivateKey(privateKey),

		grpcserver.WithHealthService(healthService),
		grpcserver.WithMetricService(metricService),
	)

	return NewGRPCServer(srv, config.GRPCAddress)
}

func setupProfilerServer(config *Config) *ProfilerServer {
	return NewProfilerServer(config)
}

func setupStorage(config *Config) (storage.MetricsStorage, error) {
	return storage.NewStorage(
		config.DatabaseDSN,
		config.StorePath,
		config.StoreInterval,
		config.RestoreOnStart,
	)
}

func preparePrivateKey(config *Config) (security.PrivateKey, error) {
	var (
		privateKey security.PrivateKey
		err        error
	)

	if len(config.PrivateKeyPath) != 0 {
		privateKey, err = security.NewPrivateKey(config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
	}

	return privateKey, err
}
