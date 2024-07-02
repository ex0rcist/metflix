package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ex0rcist/metflix/internal/storage"
)

func NewRouter(storage storage.Storage) http.Handler {
	router := chi.NewRouter()
	resource := Resource{storage: storage}

	router.Get("/", resource.Homepage) // TODO: resource?
	router.Post("/update/{metricType}/{metricName}/{metricValue}", resource.UpdateMetric)

	return router
}
