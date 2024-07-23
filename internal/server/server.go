package server

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/spf13/pflag"
)

type Server struct {
	config  *Config
	Storage storage.MetricsStorage
	Router  http.Handler
}

type Config struct {
	Address entities.Address `env:"ADDRESS"`
}

func New() (*Server, error) {
	config := &Config{
		Address: "0.0.0.0:8080",
	}

	memStorage := storage.NewMemStorage()
	storageService := storage.NewService(memStorage)

	router := NewRouter(storageService)

	return &Server{
		config:  config,
		Storage: memStorage,
		Router:  router,
	}, nil
}

func (s *Server) ParseFlags() error {
	address := s.config.Address

	pflag.VarP(&address, "address", "a", "address:port for HTTP API requests")
	pflag.Parse()

	// because VarP gets non-pointer value, set it manually
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			s.config.Address = address
		}
	})

	if err := env.Parse(s.config); err != nil {
		return entities.NewStackError(err)
	}

	return nil
}

func (s *Server) Run() error {
	logging.LogInfo(s.config.String())
	logging.LogInfo("server ready")

	return http.ListenAndServe(s.config.Address.String(), s.Router)
}

func (c Config) String() string {
	return "server config: " + fmt.Sprintf("address=%s\t", c.Address)
}
