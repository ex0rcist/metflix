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

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/httpserver"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/spf13/pflag"
)

const shutdownTimeout = 60 * time.Second

// Backend heart
type Server struct {
	config     *Config
	httpServer *httpserver.Server
	profiler   *ProfilerServer
	storage    storage.MetricsStorage
	router     http.Handler
	privateKey security.PrivateKey
}

// Backend config
type Config struct {
	Address         entities.Address  `env:"ADDRESS"`
	StoreInterval   int               `env:"STORE_INTERVAL"`
	StorePath       string            `env:"FILE_STORAGE_PATH"`
	RestoreOnStart  bool              `env:"RESTORE"`
	DatabaseDSN     string            `env:"DATABASE_DSN"`
	Secret          entities.Secret   `env:"KEY"`
	ProfilerAddress entities.Address  `env:"PROFILER_ADDRESS"`
	PrivateKeyPath  entities.FilePath `env:"CRYPTO_KEY"`
}

// Server constructor
func New() (*Server, error) {
	config := &Config{
		Address:         "0.0.0.0:8080",
		StoreInterval:   300,
		RestoreOnStart:  true,
		ProfilerAddress: "0.0.0.0:8081",
	}

	err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	dataStorage, err := newDataStorage(config)
	if err != nil {
		return nil, err
	}

	var privateKey security.PrivateKey
	if len(config.PrivateKeyPath) != 0 {
		privateKey, err = security.NewPrivateKey(config.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("server - New - security.NewPrivateKey: %w", err)
		}
	}

	storageService := storage.NewService(dataStorage)
	pingerService := services.NewPingerService(dataStorage)
	router := httpserver.NewRouter(storageService, pingerService, config.Secret, privateKey)

	httpServer := httpserver.New(router, config.Address)
	pprofiler := NewProfilerServer(config.ProfilerAddress)

	return &Server{
		config:     config,
		httpServer: httpServer,
		storage:    dataStorage,
		router:     router,
		profiler:   pprofiler,
		privateKey: privateKey,
	}, nil
}

// Start all subservices
func (s *Server) Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s.httpServer.Start()
	s.profiler.Start()

	logging.LogInfo(s.String())
	logging.LogInfo("server ready")

	select {
	case s := <-interrupt:
		logging.LogInfo("interrupt: signal " + s.String())
	case err := <-s.httpServer.Notify():
		logging.LogError(err, "Server -> Start() -> s.httpServer.Notify")
	case err := <-s.profiler.Notify():
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
	kind := detectStorageKind(s.config)

	str := []string{
		fmt.Sprintf("address=%s", s.config.Address),
		fmt.Sprintf("storage=%s", kind),
	}

	if kind == storage.KindFile {
		str = append(str, fmt.Sprintf("store-interval=%d", s.config.StoreInterval))
		str = append(str, fmt.Sprintf("store-path=%s", s.config.StorePath))
		str = append(str, fmt.Sprintf("restore=%t", s.config.RestoreOnStart))
	}

	if kind == storage.KindDatabase {
		str = append(str, fmt.Sprintf("database=%s", s.config.DatabaseDSN))
	}

	if len(s.config.Secret) > 0 {
		str = append(str, fmt.Sprintf("secret=%s", s.config.Secret))
	}

	if len(s.config.PrivateKeyPath) > 0 {
		str = append(str, fmt.Sprintf("private-key=%v", s.config.PrivateKeyPath))
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
	if err := s.profiler.Shutdown(ctx); err != nil {
		logging.LogError(err)
	}
}

func parseConfig(config *Config) error {
	err := parseFlags(config, os.Args[0], os.Args[1:])
	if err != nil {
		return err
	}

	err = parseEnv(config)
	if err != nil {
		return err
	}

	return nil
}

func parseFlags(config *Config, progname string, args []string) error {
	flags := pflag.NewFlagSet(progname, pflag.ContinueOnError)

	address := config.Address
	flags.VarP(&address, "address", "a", "address:port for HTTP API requests")

	secret := config.Secret
	flags.VarP(&secret, "secret", "k", "a key to sign outgoing data")

	privateKeyPath := config.PrivateKeyPath
	flags.VarP(&privateKeyPath, "crypto-key", "", "path to public key to encrypt agent -> server communications")

	// define flags
	flags.IntVarP(&config.StoreInterval, "store-interval", "i", config.StoreInterval, "interval (s) for dumping metrics to the disk, zero value means saving after each request")
	flags.StringVarP(&config.StorePath, "store-file", "f", config.StorePath, "path to file to store metrics")
	flags.BoolVarP(&config.RestoreOnStart, "restore", "r", config.RestoreOnStart, "whether to restore state on startup")
	flags.StringVarP(&config.DatabaseDSN, "database", "d", config.DatabaseDSN, "PostgreSQL database DSN")

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	// fill values
	flags.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			config.Address = address
		case "secret":
			config.Secret = secret
		case "crypto-key":
			config.PrivateKeyPath = privateKeyPath
		}
	})

	return nil
}

func parseEnv(config *Config) error {
	if err := env.Parse(config); err != nil {
		return err
	}

	return nil
}

func detectStorageKind(c *Config) string {
	var sk string

	switch {
	case c.DatabaseDSN != "":
		sk = storage.KindDatabase
	case c.StorePath != "":
		sk = storage.KindFile
	default:
		sk = storage.KindMemory
	}

	return sk
}

func newDataStorage(config *Config) (storage.MetricsStorage, error) {
	storageKind := detectStorageKind(config)

	switch storageKind {
	case storage.KindMemory:
		return storage.NewMemStorage(), nil
	case storage.KindFile:
		return storage.NewFileStorage(config.StorePath, config.StoreInterval, config.RestoreOnStart)
	case storage.KindDatabase:
		return storage.NewPostgresStorage(config.DatabaseDSN)
	default:
		return nil, fmt.Errorf("unknown storage type")
	}
}
