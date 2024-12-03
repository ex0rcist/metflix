package server

import (
	"context"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
)

// HTTP-server wrapper.
type HTTPServer struct {
	server *http.Server
	notify chan error
}

// Constructor.
func NewHTTPServer(handler http.Handler, address entities.Address) *HTTPServer {
	httpServer := &http.Server{
		Handler: handler,
		Addr:    address.String(),
	}

	return &HTTPServer{
		server: httpServer,
		notify: make(chan error, 1),
	}
}

// Run server in a goroutine.
func (s *HTTPServer) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Return channel to handle errors.
func (s *HTTPServer) Notify() <-chan error {
	return s.notify
}

// Shutdown server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
