// Package httpserver implements REST API for metrics collector server.
package httpserver

// @Title Metrics collector API
// @Description Service for storing metrics data.
// @Version 1.0

// @Contact.name  Evgeniy Shuvalov
// @Contact.email evshuvalov@yandex.ru

// @Tag.name Metrics
// @Tag.description "Metrics API"

// @Tag.name Healthcheck
// @Tag.description "API to inspect service health state"

import (
	"net"
	"net/http"

	_ "github.com/ex0rcist/metflix/docs/api"
	httpSwagger "github.com/swaggo/http-swagger"

	chimdlw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/httpserver/middleware"
	"github.com/ex0rcist/metflix/internal/security"
)

type Backend struct {
	router        *chi.Mux
	signSecret    entities.Secret
	privateKey    security.PrivateKey
	trustedSubnet *net.IPNet

	healthResource *HealthResource
	metricResource *MetricResource
}

// Backend constructor
func NewBackend(opts ...Option) http.Handler {
	backend := &Backend{
		router: chi.NewRouter(),
	}

	for _, opt := range opts {
		opt(backend)
	}

	backend.registerMiddlewares()
	backend.registerEndpoints()

	return backend.router
}

func (b *Backend) registerMiddlewares() {
	middlewares := []func(http.Handler) http.Handler{
		chimdlw.RealIP,
		chimdlw.StripSlashes,
		middleware.RequestsLogger,

		func(next http.Handler) http.Handler {
			return middleware.CheckSignedRequest(next, b.signSecret)
		},

		func(next http.Handler) http.Handler {
			return middleware.DecryptRequest(next, b.privateKey)
		},

		func(next http.Handler) http.Handler {
			return middleware.FilterUntrustedRequest(next, b.trustedSubnet)
		},

		middleware.DecompressRequest,
		middleware.CompressResponse,

		func(next http.Handler) http.Handler {
			return middleware.SignResponse(next, b.signSecret)
		},
	}

	b.router.Use(middlewares...)
}

func (b *Backend) registerEndpoints() {
	b.registerMetricsEndpoints()
	b.registerHealthEndpoint()

	// setup default 404
	b.router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // no default body
	}))

	// setup documentation
	b.router.Get("/swagger/*", httpSwagger.WrapHandler)
}

func (b *Backend) registerMetricsEndpoints() {
	if b.metricResource == nil {
		return
	}

	b.router.Get("/", b.metricResource.Homepage)

	b.router.Post("/update/{metricKind}/{metricName}/{metricValue}", b.metricResource.UpdateMetric)
	b.router.Post("/update", b.metricResource.UpdateMetricJSON)
	b.router.Post("/updates", b.metricResource.UpdateMetricsBatch)

	b.router.Get("/value/{metricKind}/{metricName}", b.metricResource.GetMetric)
	b.router.Post("/value", b.metricResource.GetMetricJSON)
}

func (b *Backend) registerHealthEndpoint() {
	if b.healthResource == nil {
		return
	}

	b.router.Get("/ping", b.healthResource.Ping)
}

/* Options */

type Option func(*Backend)

func WithSignSecret(secret entities.Secret) Option {
	return func(b *Backend) {
		b.signSecret = secret
	}
}

func WithPrivateKey(privateKey security.PrivateKey) Option {
	return func(b *Backend) {
		b.privateKey = privateKey
	}
}

func WithTrustedSubnet(trustedSubnet *net.IPNet) Option {
	return func(b *Backend) {
		b.trustedSubnet = trustedSubnet
	}
}

func WithHealthResource(healthResource *HealthResource) Option {
	return func(b *Backend) {
		b.healthResource = healthResource
	}
}

func WithMetricResource(metricResource *MetricResource) Option {
	return func(b *Backend) {
		b.metricResource = metricResource
	}
}
