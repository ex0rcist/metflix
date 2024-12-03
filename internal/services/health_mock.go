package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ HealthChecker = (*HealthCheckServiceMock)(nil)

type HealthCheckServiceMock struct {
	mock.Mock
}

func (m *HealthCheckServiceMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
