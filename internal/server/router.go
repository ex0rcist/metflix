package server

import (
	"net/http"

	chimdlw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/middleware"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/storage"
)

func NewRouter(
	storageService storage.StorageService,
	pingerService services.Pinger,
	secret entities.Secret,
) http.Handler {
	router := chi.NewRouter()

	router.Use(chimdlw.RealIP)
	router.Use(chimdlw.StripSlashes)

	router.Use(middleware.RequestsLogger)

	router.Use(func(next http.Handler) http.Handler {
		return middleware.CheckSignedRequest(next, secret)
	})

	router.Use(middleware.DecompressRequest)
	router.Use(middleware.CompressResponse)

	router.Use(func(next http.Handler) http.Handler {
		return middleware.SignResponse(next, secret)
	})

	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // no default body
	}))

	registerMetricsEndpoints(storageService, router)
	registerPingerEndpoint(pingerService, router)

	return router
}

func registerMetricsEndpoints(storageService storage.StorageService, router *chi.Mux) {
	resource := NewMetricResource(storageService)

	router.Get("/", resource.Homepage)

	router.Post("/update/{metricKind}/{metricName}/{metricValue}", resource.UpdateMetric)
	router.Post("/update", resource.UpdateMetricJSON)
	router.Post("/updates", resource.BatchUpdateMetricsJSON)

	router.Get("/value/{metricKind}/{metricName}", resource.GetMetric)
	router.Post("/value", resource.GetMetricJSON)
}

func registerPingerEndpoint(pingerService services.Pinger, router *chi.Mux) {
	resource := NewPingerResource(pingerService)

	router.Get("/ping", resource.Ping)
}
