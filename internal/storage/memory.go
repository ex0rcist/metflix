package storage

import (
	"github.com/ex0rcist/metflix/internal/entities"
)

// check that MemStorage implements MetricsStorage
var _ MetricsStorage = (*MemStorage)(nil)

type MemStorage struct {
	data map[string]Record
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]Record),
	}
}

func (s *MemStorage) Push(id string, record Record) error {
	s.data[id] = record

	return nil
}

func (s *MemStorage) Get(id string) (Record, error) {
	record, ok := s.data[id]
	if !ok {
		return Record{}, entities.ErrRecordNotFound
	}

	return record, nil
}

func (s *MemStorage) List() ([]Record, error) {
	arr := make([]Record, len(s.data))

	i := 0
	for _, record := range s.data {
		arr[i] = record
		i++
	}

	return arr, nil
}
