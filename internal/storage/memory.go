package storage

import (
	"github.com/ex0rcist/metflix/internal/entities"
)

// check that MemStorage implements MetricsStorage
var _ MetricsStorage = (*MemStorage)(nil)

type MemStorage struct {
	Data map[string]Record
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Data: make(map[string]Record),
	}
}

func (s *MemStorage) Push(id string, record Record) error {
	s.Data[id] = record

	return nil
}

func (s *MemStorage) Get(id string) (Record, error) {
	record, ok := s.Data[id]
	if !ok {
		return Record{}, entities.ErrMetricNotFound
	}

	return record, nil
}

func (s *MemStorage) List() ([]Record, error) {
	arr := make([]Record, len(s.Data))

	i := 0
	for _, record := range s.Data {
		arr[i] = record
		i++
	}

	return arr, nil
}
