package storage

import (
	"errors"
)

type MemStorage struct {
	Data map[string]Record
}

func (strg *MemStorage) Push(record Record) error {
	id := RecordID(record.Name, record.Value.Kind())
	strg.Data[id] = record

	return nil
}

func (strg MemStorage) Get(id string) (Record, error) {
	record, ok := strg.Data[id]
	if !ok {
		return Record{}, errors.New("no value")
	}

	return record, nil
}

func (strg MemStorage) GetAll() ([]Record, error) {
	arr := make([]Record, len(strg.Data))
	i := 0

	for _, v := range strg.Data {
		arr[i] = v
		i++
	}

	return arr, nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Data: make(map[string]Record),
	}
}
