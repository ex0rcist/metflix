package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Pinger = (*PingerMock)(nil)

type PingerMock struct {
	mock.Mock
}

func (m *PingerMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
