package storage

import (
	"errors"
	"sort"
)

type MemStorage struct {
	Data map[RecordID]Record
}

func (strg *MemStorage) Push(record Record) error {
	recordID := CalculateRecordID(record.Name, record.Value.Kind())
	strg.Data[recordID] = record

	return nil
}

type By func(r1, r2 *Record) bool

func (strg MemStorage) Get(recordID RecordID) (Record, error) {
	record, ok := strg.Data[recordID]
	if !ok {
		return Record{}, errors.New("no value")
	}

	return record, nil
}

// Sorted by name
// Too complicated? mb easy way?
func (strg MemStorage) GetAll() ([]Record, error) {
	names := make([]string, 0, len(strg.Data))

	for k := range strg.Data {
		names = append(names, string(k))
	}

	arr := make([]Record, len(strg.Data))
	sort.Strings(names)

	for i, v := range names {
		arr[i] = strg.Data[RecordID(v)]
	}

	return arr, nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Data: make(map[RecordID]Record),
	}
}
