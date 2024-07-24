package storage

import (
	"errors"
	"sort"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/metrics"
)

// Common interface for service layer
// Storage service prepares data before calling storage
type StorageService interface {
	List() ([]Record, error)
	Push(record Record) (Record, error)
	Get(name, kind string) (Record, error)
}

// ensure Service implements StorageService
var _ StorageService = Service{}

type Service struct {
	storage MetricsStorage
}

func NewService(storage MetricsStorage) Service {
	return Service{storage: storage}
}

func (s Service) Get(name, kind string) (Record, error) {
	id := CalculateRecordID(name, kind)

	record, err := s.storage.Get(id)
	if err != nil {
		return Record{}, err
	}

	return record, nil
}

func (s Service) Push(record Record) (Record, error) {
	newValue, err := s.calculateNewValue(record)
	if err != nil {
		return Record{}, err
	}

	record.Value = newValue
	err = s.storage.Push(record.CalculateRecordID(), record)

	if err != nil {
		return Record{}, err
	}

	return record, nil
}

func (s Service) List() ([]Record, error) {
	records, err := s.storage.List()
	if err != nil {
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})

	return records, nil
}

func (s Service) calculateNewValue(record Record) (metrics.Metric, error) {
	if record.Value.Kind() != "counter" {
		return record.Value, nil
	}

	id := record.CalculateRecordID()
	if id == "" {
		return record.Value, entities.ErrMetricMissingName
	}

	storedRecord, err := s.storage.Get(id)
	if errors.Is(err, entities.ErrRecordNotFound) {
		return record.Value, nil
	} else if err != nil {
		return nil, err
	}

	return storedRecord.Value.(metrics.Counter) + record.Value.(metrics.Counter), nil
}
