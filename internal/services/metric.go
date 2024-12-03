package services

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/pkg/metrics"
)

type MetricProvider interface {
	List(ctx context.Context) ([]storage.Record, error)
	Push(ctx context.Context, record storage.Record) (storage.Record, error)
	PushList(ctx context.Context, records []storage.Record) ([]storage.Record, error)
	Get(ctx context.Context, name, kind string) (storage.Record, error)
}

var _ MetricProvider = MetricService{}

// Service struct, containing storage
type MetricService struct {
	storage storage.MetricsStorage
}

// Service constructor
func NewMetricService(storage storage.MetricsStorage) MetricService {
	return MetricService{storage: storage}
}

// Get record from bound storage
func (s MetricService) Get(ctx context.Context, name, kind string) (storage.Record, error) {
	id := storage.CalculateRecordID(name, kind)

	record, err := s.storage.Get(ctx, id)
	if err != nil {
		return storage.Record{}, err
	}

	return record, nil
}

// Push record to bound storage
func (s MetricService) Push(ctx context.Context, record storage.Record) (storage.Record, error) {
	newValue, err := s.calculateNewValue(ctx, record)
	if err != nil {
		return storage.Record{}, err
	}

	record.Value = newValue
	err = s.storage.Push(ctx, record.CalculateRecordID(), record)

	if err != nil {
		return storage.Record{}, err
	}

	return record, nil
}

// Push list of records to bound storage
func (s MetricService) PushList(ctx context.Context, records []storage.Record) ([]storage.Record, error) {
	data := make(map[string]storage.Record)

	for _, record := range records {
		id := record.CalculateRecordID()

		if prev, ok := data[id]; ok {
			if record.Value.Kind() == metrics.KindCounter {
				record.Value = prev.Value.(metrics.Counter) + record.Value.(metrics.Counter)
			}

			data[id] = record

			continue
		}

		newValue, err := s.calculateNewValue(ctx, record)
		if err != nil {
			return nil, fmt.Errorf("unable to calculate new value: %w", err)
		}

		record.Value = newValue
		data[id] = record
	}

	if err := s.storage.PushList(ctx, data); err != nil {
		return nil, fmt.Errorf("unable to PushList(): %w", err)
	}

	result := make([]storage.Record, 0, len(data))
	for _, v := range data {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// List records from bound storage
func (s MetricService) List(ctx context.Context) ([]storage.Record, error) {
	records, err := s.storage.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})

	return records, nil
}

func (s MetricService) calculateNewValue(ctx context.Context, record storage.Record) (metrics.Metric, error) {
	if record.Value.Kind() != metrics.KindCounter {
		return record.Value, nil
	}

	id := record.CalculateRecordID()
	if id == "" {
		return record.Value, entities.ErrMetricMissingName
	}

	storedRecord, err := s.storage.Get(ctx, id)
	if errors.Is(err, entities.ErrRecordNotFound) {
		return record.Value, nil
	} else if err != nil {
		return nil, err
	}

	return storedRecord.Value.(metrics.Counter) + record.Value.(metrics.Counter), nil
}
