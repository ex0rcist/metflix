package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
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
	Address entities.Address `env:"ADDRESS"`
}

func New() (*Server, error) {
	config := &Config{
		Address: "0.0.0.0:8080",
	}

	err := parseFlags(config, os.Args[0], os.Args[1:])
	if err != nil {
		return nil, err
	}

	err = parseEnv(config)
	if err != nil {
		return nil, err
	}

	memStorage := storage.NewMemStorage()
	storageService := storage.NewService(memStorage)
	router := NewRouter(storageService)

	httpServer := &http.Server{
		Addr:    config.Address.String(),
		Handler: router,
	}

	return &Server{
		config:     config,
		httpServer: httpServer,
		Storage:    memStorage,
		Router:     router,
	}, nil
}

func parseFlags(config *Config, progname string, args []string) error {
	flags := pflag.NewFlagSet(progname, pflag.ContinueOnError)

	address := config.Address

	flags.VarP(&address, "address", "a", "address:port for HTTP API requests")
	err := flags.Parse(args)

	if err != nil {
		return err
	}

	// because VarP gets non-pointer value, set it manually
	flags.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			config.Address = address
		}
	})

	return nil
}

func parseEnv(config *Config) error {
	fmt.Println("====== " + os.Getenv("ADDRESS") + " ======")
	if err := env.Parse(config); err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() error {
	logging.LogInfo(s.config.String())
	logging.LogInfo("server ready")

	return s.httpServer.ListenAndServe()
}

func (c Config) String() string {
	return "server config: " + fmt.Sprintf("address=%s", c.Address)
}
