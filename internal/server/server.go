package server

import (
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/spf13/pflag"
)

type Server struct {
	Config  *Config
	Storage storage.Storage
	Router  http.Handler
}

type Config struct {
	Address entities.Address
}

func New() (*Server, error) {
	config := &Config{
		Address: "0.0.0.0:8080",
	}

	storage := storage.NewMemStorage()
	router := NewRouter(storage)

	return &Server{
		Config:  config,
		Storage: storage,
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

	return nil
}

func (s *Server) Run() error {
	err := http.ListenAndServe(s.Config.Address.String(), s.Router) // HELP: почему тип, реализующий String() не приводится к строке автоматически?
	return err
}
