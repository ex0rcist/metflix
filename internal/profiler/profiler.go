package profiler

import (
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/httpserver"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "net/http/pprof"
)

type Profiler struct {
	*httpserver.Server
}

func New(address entities.Address) *Profiler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Mount("/debug", middleware.Profiler())

	// tweak memory profiling rate to spot more allocations.
	// runtime.MemProfileRate = 2048

	server := httpserver.New(r, address)

	return &Profiler{server}
}
