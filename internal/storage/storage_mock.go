package storage

import (
	"github.com/stretchr/testify/mock"
)

// Ensure StorageMock implements MetricsStorage
var _ MetricsStorage = (*StorageMock)(nil)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) Get(id string) (Record, error) {
	args := m.Called(id)
	return args.Get(0).(Record), args.Error(1)
}

func (m *StorageMock) Push(id string, record Record) error {
	args := m.Called(id, record)
	return args.Error(0)
}

func (m *StorageMock) List() ([]Record, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Record), args.Error(1)
}

func (m *StorageMock) Kind() string {
	return KindMock
}
