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
	id := CalculateRecordID(record.Name, record.Value.Kind())

	newValue, err := s.calculateNewValue(id, record)
	if err != nil {
		return Record{}, err
	}

	record.Value = newValue
	err = s.storage.Push(id, record)

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

func (s Service) calculateNewValue(id string, newRecord Record) (metrics.Metric, error) {
	if newRecord.Value.Kind() != "counter" {
		return newRecord.Value, nil
	}

	storedRecord, err := s.storage.Get(id)
	if errors.Is(err, entities.ErrMetricNotFound) {
		return newRecord.Value, nil
	} else if err != nil {
		return nil, err
	}

	return storedRecord.Value.(metrics.Counter) + newRecord.Value.(metrics.Counter), nil
}
