package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Profiler server to serve pprof
type ProfilerServer struct {
	server *http.Server
	notify chan error
}

// Profiler server constructor
func NewProfilerServer(config *Config) *ProfilerServer {
	handler := chi.NewRouter()

	handler.Use(middleware.Logger)
	handler.Mount("/debug", middleware.Profiler())

	httpServer := &http.Server{
		Handler: handler,
		Addr:    config.ProfilerAddress.String(),
	}

	return &ProfilerServer{server: httpServer, notify: make(chan error, 1)}
}

// Run server in a goroutine.
func (s *ProfilerServer) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Return channel to handle errors.
func (s *ProfilerServer) Notify() <-chan error {
	return s.notify
}

// Shutdown server.
func (s *ProfilerServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
