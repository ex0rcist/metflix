package server

import (
	"net/http"

	chimdlw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/ex0rcist/metflix/internal/middleware"
	"github.com/ex0rcist/metflix/internal/storage"
)

func NewRouter(storageService storage.StorageService) http.Handler {
	router := chi.NewRouter()

	router.Use(chimdlw.RealIP)
	router.Use(chimdlw.StripSlashes)

	router.Use(middleware.RequestsLogger)
	router.Use(middleware.DecompressRequest)
	router.Use(middleware.CompressResponse)

	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // no default body
	}))

	resource := NewMetricResource(storageService)

	router.Get("/", resource.Homepage)

	router.Post("/update/{metricKind}/{metricName}/{metricValue}", resource.UpdateMetric)
	router.Post("/update", resource.UpdateMetricJSON)

	router.Get("/value/{metricKind}/{metricName}", resource.GetMetric)
	router.Post("/value", resource.GetMetricJSON)

	return router
}
