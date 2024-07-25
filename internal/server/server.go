package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/utils"
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
	router := NewRouter(storageService)

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

	// restore storage if possible
	if s.storageNeedsRestore() {
		if err := s.restoreStorage(); err != nil {
			return err
		}
	}

	logging.LogInfo("server ready")

	// start dumping if neede
	if s.storageNeedsDumping() {
		go s.startStorageDumping()
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) String() string {
	str := []string{
		fmt.Sprintf("address=%s", s.config.Address),
		fmt.Sprintf("storage=%s", s.Storage.Kind()),
	}

	if s.Storage.Kind() == storage.KindFile {
		str = append(str, fmt.Sprintf("store-interval=%d", s.config.StoreInterval))
		str = append(str, fmt.Sprintf("store-path=%s", s.config.StorePath))
		str = append(str, fmt.Sprintf("restore=%t", s.config.RestoreOnStart))
	}

	return "server config: " + strings.Join(str, "; ")
}

func (s *Server) storageNeedsRestore() bool {
	return s.Storage.Kind() == storage.KindFile && s.config.RestoreOnStart
}

func (s *Server) restoreStorage() error {
	// HELP: не уверен что тут корректное решение... Но иначе нужно добавлять Restore() в интерфейс, а его логически нет у MemStorage
	return s.Storage.(*storage.FileStorage).Restore()
}

func (s *Server) storageNeedsDumping() bool {
	return s.Storage.Kind() == storage.KindFile && s.config.StoreInterval > 0
}

func (s *Server) startStorageDumping() {
	ticker := time.NewTicker(utils.IntToDuration(s.config.StoreInterval))
	defer ticker.Stop()

	for {
		_, ok := <-ticker.C
		if !ok {
			break
		}

		// HELP: не уверен что тут корректное решение... Но иначе нужно добавлять Dump() в интерфейс, а его логически нет у MemStorage
		if err := s.Storage.(*storage.FileStorage).Dump(); err != nil {
			logging.LogError(fmt.Errorf("error during FileStorage Dump(): %s", err.Error()))
		}
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

	// define flags
	storeInterval := flags.IntP("store-interval", "i", config.StoreInterval, "interval (s) for dumping metrics to the disk, zero value means saving after each request")
	storePath := flags.StringP("store-file", "f", config.StorePath, "path to file to store metrics")
	restoreOnStart := flags.BoolP("restore", "r", config.RestoreOnStart, "whether to restore state on startup")

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
		return storage.NewFileStorage(config.StorePath, config.StoreInterval), nil
	default:
		return nil, fmt.Errorf("unknown storage type")
	}
}
