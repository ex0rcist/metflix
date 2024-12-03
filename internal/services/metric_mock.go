package services

import (
	"context"

	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/stretchr/testify/mock"
)

var _ MetricProvider = (*MetricServiceMock)(nil)

// Service mock
type MetricServiceMock struct {
	mock.Mock
}

// Get record
func (m *MetricServiceMock) Get(ctx context.Context, name, kind string) (storage.Record, error) {
	args := m.Called(name, kind)
	return args.Get(0).(storage.Record), args.Error(1)
}

// Push record
func (m *MetricServiceMock) Push(ctx context.Context, record storage.Record) (storage.Record, error) {
	args := m.Called(record)
	return args.Get(0).(storage.Record), args.Error(1)
}

// Push list of records
func (m *MetricServiceMock) PushList(ctx context.Context, records []storage.Record) ([]storage.Record, error) {
	args := m.Called(ctx, records)
	return args.Get(0).([]storage.Record), args.Error(1)
}

// Get list of records
func (m *MetricServiceMock) List(ctx context.Context) ([]storage.Record, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]storage.Record), args.Error(1)
}
