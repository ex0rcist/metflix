package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Ensure StorageMock implements MetricsStorage
var _ MetricsStorage = (*StorageMock)(nil)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) Get(ctx context.Context, id string) (Record, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Record), args.Error(1)
}

func (m *StorageMock) Push(ctx context.Context, id string, record Record) error {
	args := m.Called(ctx, id, record)
	return args.Error(0)
}

func (m *StorageMock) PushList(ctx context.Context, data map[string]Record) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *StorageMock) List(ctx context.Context) ([]Record, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Record), args.Error(1)
}

func (m *StorageMock) Close(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}
