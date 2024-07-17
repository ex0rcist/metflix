package server

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/storage"
)

func NewRouter(storageService storage.StorageService) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(logging.RequestsLogger)

	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // no default body
	}))

	resource := NewMetricResource(storageService)

	router.Get("/", resource.Homepage) // TODO: resource?

	router.Post("/update/{metricKind}/{metricName}/{metricValue}", resource.UpdateMetric)
	router.Post("/update", resource.UpdateMetricJSON)

	router.Get("/value/{metricKind}/{metricName}", resource.GetMetric)
	router.Post("/value", resource.GetMetricJSON)

	return router
}
