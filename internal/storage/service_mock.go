package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ StorageService = (*ServiceMock)(nil)

// Service mock
type ServiceMock struct {
	mock.Mock
}

// Get record
func (m *ServiceMock) Get(ctx context.Context, name, kind string) (Record, error) {
	args := m.Called(name, kind)
	return args.Get(0).(Record), args.Error(1)
}

// Push record
func (m *ServiceMock) Push(ctx context.Context, record Record) (Record, error) {
	args := m.Called(record)
	return args.Get(0).(Record), args.Error(1)
}

// Push list of records
func (m *ServiceMock) PushList(ctx context.Context, records []Record) ([]Record, error) {
	args := m.Called(ctx, records)
	return args.Get(0).([]Record), args.Error(1)
}

// Get list of records
func (m *ServiceMock) List(ctx context.Context) ([]Record, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Record), args.Error(1)
}
