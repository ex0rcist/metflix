package httpserver

import (
	"context"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
)

type Server struct {
	server *http.Server
	notify chan error
}

func New(handler http.Handler, address entities.Address) *Server {
	httpServer := &http.Server{
		Handler: handler,
		Addr:    address.String(),
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
	}

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
