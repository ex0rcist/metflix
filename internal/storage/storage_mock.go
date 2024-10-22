package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ MetricsStorage = (*StorageMock)(nil)

// Storage mock
type StorageMock struct {
	mock.Mock
}

// Get record
func (m *StorageMock) Get(ctx context.Context, id string) (Record, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Record), args.Error(1)
}

// Push record
func (m *StorageMock) Push(ctx context.Context, id string, record Record) error {
	args := m.Called(ctx, id, record)
	return args.Error(0)
}

// Push list of records
func (m *StorageMock) PushList(ctx context.Context, data map[string]Record) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

// List records
func (m *StorageMock) List(ctx context.Context) ([]Record, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Record), args.Error(1)
}

// Close storage
func (m *StorageMock) Close(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}
