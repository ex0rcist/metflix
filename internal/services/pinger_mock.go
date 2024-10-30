package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Pinger = (*PingerMock)(nil)

// Pinger mock
type PingerMock struct {
	mock.Mock
}

// Ping service
func (m *PingerMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
