package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Ensure ServiceMock implements StorageService
var _ StorageService = (*ServiceMock)(nil)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) Get(ctx context.Context, name, kind string) (Record, error) {
	args := m.Called(name, kind)
	return args.Get(0).(Record), args.Error(1)
}

func (m *ServiceMock) Push(ctx context.Context, record Record) (Record, error) {
	args := m.Called(record)
	return args.Get(0).(Record), args.Error(1)
}

func (m *ServiceMock) PushList(ctx context.Context, records []Record) ([]Record, error) {
	args := m.Called(ctx, records)
	return args.Get(0).([]Record), args.Error(1)
}

func (m *ServiceMock) List(ctx context.Context) ([]Record, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Record), args.Error(1)
}
