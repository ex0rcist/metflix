package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ex0rcist/metflix/internal/storage"
)

func NewRouter(storage storage.Storage) http.Handler {
	router := chi.NewRouter()
	resource := Resource{storage: storage}

	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no default body
		w.WriteHeader(http.StatusNotFound)
	}))

	router.Get("/", resource.Homepage) // TODO: resource?
	router.Post("/update/{metricKind}/{metricName}/{metricValue}", resource.UpdateMetric)
	router.Get("/value/{metricKind}/{metricName}", resource.ShowMetric)

	return router
}
