package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
	httpServer     *httpserver.Server
	profilerServer *ProfilerServer
	storage        storage.MetricsStorage
	router         http.Handler
	privateKey     security.PrivateKey
}

// Server constructor
func New() (*Server, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	dataStorage, err := storage.NewStorage(
		config.DatabaseDSN,
		config.StorePath,
		config.StoreInterval,
		config.RestoreOnStart,
	)
	if err != nil {
		return nil, err
	}

	var privateKey security.PrivateKey
	if len(config.PrivateKeyPath) != 0 {
		privateKey, err = security.NewPrivateKey(config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
	}

	storageService := storage.NewService(dataStorage)
	pingerService := services.NewPingerService(dataStorage)
	router := httpserver.NewRouter(
		storageService,
		pingerService,
		config.Secret,
		privateKey,
		config.TrustedSubnet,
	)

	httpServer := httpserver.New(router, config.Address)
	profilerServer := NewProfilerServer(config.ProfilerAddress)

	return &Server{
		config:         config,
		httpServer:     httpServer,
		profilerServer: profilerServer,
		storage:        dataStorage,
		router:         router,
		privateKey:     privateKey,
	}, nil
}

// Start all subservices
func (s *Server) Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s.httpServer.Start()
	s.profilerServer.Start()

	logging.LogInfo(s.String())
	logging.LogInfo("server ready")

	select {
	case s := <-interrupt:
		logging.LogInfo("interrupt: signal " + s.String())
	case err := <-s.httpServer.Notify():
		logging.LogError(err, "Server -> Start() -> s.httpServer.Notify")
	case err := <-s.profilerServer.Notify():
		logging.LogError(err, "Server -> Start() - s.profiler.Notify")
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

	logging.LogInfo("shutting down storage")
	if err := s.storage.Close(ctx); err != nil {
		logging.LogError(err)
	}

	logging.LogInfo("shutting down profiler")
	if err := s.profilerServer.Shutdown(ctx); err != nil {
		logging.LogError(err)
	}
}
