package server

import (
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/httpserver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ProfilerServer struct {
	*httpserver.Server
}

func NewProfilerServer(address entities.Address) *ProfilerServer {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Mount("/debug", middleware.Profiler())

	server := httpserver.New(r, address)

	prf := &ProfilerServer{server}

	return prf
}
