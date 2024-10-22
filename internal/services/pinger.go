package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/storage"
)

var _ Pinger = PingerService{}

// Pinger interface
type Pinger interface {
	Ping(ctx context.Context) error
}

// Pinger service to make healthcheck
type PingerService struct {
	storage storage.MetricsStorage
}

// Pinger constructor.
func NewPingerService(storage storage.MetricsStorage) PingerService {
	return PingerService{storage: storage}
}

// Ping-pong.
func (s PingerService) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	strg, ok := s.storage.(Pinger)
	if !ok {
		return fmt.Errorf("storage ping failed: %w", entities.ErrStorageUnpingable)
	}

	if err := strg.Ping(ctx); err != nil {
		return fmt.Errorf("storage check failed: %w", err)
	}

	return nil
}
