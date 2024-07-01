package server

import (
	"net/http"

	"github.com/ex0rcist/metflix/internal/storage"
)

type Server struct {
	Storage storage.Storage
	Router  http.Handler
}

func New() (*Server, error) {
	storage := storage.NewMemStorage()
	router := NewRouter(storage)

	return &Server{
		Storage: storage,
		Router:  router,
	}, nil
}

func (s *Server) Run() error {
	err := http.ListenAndServe(`:8080`, s.Router)
	return err
}
