package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/spf13/pflag"
)

type Server struct {
	config     *Config
	httpServer *http.Server
	Storage    storage.MetricsStorage
	Router     http.Handler
}

type Config struct {
	Address        entities.Address `env:"ADDRESS"`
	StoreInterval  int              `env:"STORE_INTERVAL"`
	StorePath      string           `env:"FILE_STORAGE_PATH"`
	RestoreOnStart bool             `env:"RESTORE"`
	DatabaseDSN    string           `env:"DATABASE_DSN"`
}

func New() (*Server, error) {
	config := &Config{
		Address:        "0.0.0.0:8080",
		StoreInterval:  300,
		RestoreOnStart: true,
	}

	err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	storageKind := detectStorageKind(config)
	dataStorage, err := newDataStorage(storageKind, config)
	if err != nil {
		return nil, err
	}

	storageService := storage.NewService(dataStorage)
	pingerService := services.NewPingerService(dataStorage)
	router := NewRouter(storageService, pingerService)

	httpServer := &http.Server{
		Addr:    config.Address.String(),
		Handler: router,
	}

	return &Server{
		config:     config,
		httpServer: httpServer,
		Storage:    dataStorage,
		Router:     router,
	}, nil
}

func (s *Server) Run() error {
	logging.LogInfo(s.String())
	logging.LogInfo("server ready")

	return s.httpServer.ListenAndServe()
}

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

	return "server config: " + strings.Join(str, "; ")
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

	// define flags
	storeInterval := flags.IntP("store-interval", "i", config.StoreInterval, "interval (s) for dumping metrics to the disk, zero value means saving after each request")
	storePath := flags.StringP("store-file", "f", config.StorePath, "path to file to store metrics")
	restoreOnStart := flags.BoolP("restore", "r", config.RestoreOnStart, "whether to restore state on startup")
	databaseDSN := flags.StringP("database", "d", config.DatabaseDSN, "PostgreSQL database DSN")

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	// fill values
	flags.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			config.Address = address
		case "store-interval":
			config.StoreInterval = *storeInterval
		case "store-file":
			config.StorePath = *storePath
		case "restore":
			config.RestoreOnStart = *restoreOnStart
		case "database":
			config.DatabaseDSN = *databaseDSN
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

func newDataStorage(kind string, config *Config) (storage.MetricsStorage, error) {
	switch kind {
	case storage.KindMemory:
		return storage.NewMemStorage(), nil
	case storage.KindFile:
		return storage.NewFileStorage(config.StorePath, config.StoreInterval, config.RestoreOnStart)
	case storage.KindDatabase:
		return storage.NewDatabaseStorage(config.DatabaseDSN)
	default:
		return nil, fmt.Errorf("unknown storage type")
	}
}
