package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/storage"
)

const defaultTimeout = 5 * time.Second

var _ HealthChecker = HealthCheckService{}

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type HealthCheckService struct {
	storage storage.MetricsStorage
}

// Interface to check if storage supports healthcheck
type PingableStorage interface {
	Ping(ctx context.Context) error
}

// Pinger constructor.
func NewHealthCheckService(storage storage.MetricsStorage) *HealthCheckService {
	return &HealthCheckService{storage: storage}
}

// Ping-pong.
func (s HealthCheckService) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	strg, ok := s.storage.(PingableStorage)
	if !ok {
		return fmt.Errorf("storage ping failed: %w", entities.ErrStorageUnpingable)
	}

	if err := strg.Ping(ctx); err != nil {
		return fmt.Errorf("storage check failed: %w", err)
	}

	return nil
}
