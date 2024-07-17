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
	Config  *Config
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
		Config:  config,
		Storage: memStorage,
		Router:  router,
	}, nil
}

func (s *Server) ParseFlags() error {
	address := s.Config.Address

	pflag.VarP(&address, "address", "a", "address:port for HTTP API requests") // HELP: "&"" because Set() has pointer receiver?
	pflag.Parse()

	// because VarP gets non-pointer value, set it manually
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			s.Config.Address = address
		}
	})

	if err := env.Parse(s.Config); err != nil {
		return logging.NewError(err)
	}

	return nil
}

func (s *Server) Run() error {
	// HELP: почему тип, реализующий String() не приводится к строке автоматически?
	err := http.ListenAndServe(s.Config.Address.String(), s.Router)
	return err
}

func (c Config) String() string {
	out := "server config: "

	out += fmt.Sprintf("address=%s\t", c.Address)
	return out
}
