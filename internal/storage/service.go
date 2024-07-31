package storage

import (
	"context"
	"errors"
	"sort"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/metrics"
)

// Common interface for service layer
// Storage service prepares data before calling storage
type StorageService interface {
	List(ctx context.Context) ([]Record, error)
	Push(ctx context.Context, record Record) (Record, error)
	Get(ctx context.Context, name, kind string) (Record, error)
}

// ensure Service implements StorageService
var _ StorageService = Service{}

type Service struct {
	Storage MetricsStorage
}

func NewService(storage MetricsStorage) Service {
	return Service{Storage: storage}
}

func (s Service) Get(ctx context.Context, name, kind string) (Record, error) {
	id := CalculateRecordID(name, kind)

	record, err := s.Storage.Get(ctx, id)
	if err != nil {
		return Record{}, err
	}

	return record, nil
}

func (s Service) Push(ctx context.Context, record Record) (Record, error) {
	newValue, err := s.calculateNewValue(ctx, record)
	if err != nil {
		return Record{}, err
	}

	record.Value = newValue
	err = s.Storage.Push(ctx, record.CalculateRecordID(), record)

	if err != nil {
		return Record{}, err
	}

	return record, nil
}

func (s Service) List(ctx context.Context) ([]Record, error) {
	records, err := s.Storage.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})

	return records, nil
}

func (s Service) calculateNewValue(ctx context.Context, record Record) (metrics.Metric, error) {
	if record.Value.Kind() != metrics.KindCounter {
		return record.Value, nil
	}

	id := record.CalculateRecordID()
	if id == "" {
		return record.Value, entities.ErrMetricMissingName
	}

	storedRecord, err := s.Storage.Get(ctx, id)
	if errors.Is(err, entities.ErrRecordNotFound) {
		return record.Value, nil
	} else if err != nil {
		return nil, err
	}

	return storedRecord.Value.(metrics.Counter) + record.Value.(metrics.Counter), nil
}
