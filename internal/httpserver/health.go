package httpserver

import (
	"errors"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/services"
)

// Resource to handle ping requests
type HealthResource struct {
	healthService services.HealthChecker
}

// Constructor
func NewHealthResource(healthService services.HealthChecker) *HealthResource {
	return &HealthResource{
		healthService: healthService,
	}
}

// Ping godoc
// @Tags Healthcheck
// @Router /ping [get]
// @Summary Verify server up and running
// @ID health_info
// @Success 200
// @Failure 500 {string} string http.StatusInternalServerError
// @Failure 501 {string} string http.StatusNotImplemented
func (res HealthResource) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := res.healthService.Ping(ctx)
	if err == nil {
		return
	}

	if errors.Is(err, entities.ErrStorageUnpingable) {
		writeErrorResponse(ctx, w, http.StatusNotImplemented, err)
		return
	}

	writeErrorResponse(ctx, w, http.StatusInternalServerError, err)
}
