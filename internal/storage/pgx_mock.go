package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

var _ PGXPool = (*PGXPoolMock)(nil)

type PGXPoolMock struct {
	mock.Mock
}

func NewPGXPoolMock() *PGXPoolMock {
	return new(PGXPoolMock)
}

func (m *PGXPoolMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}

func (m *PGXPoolMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *PGXPoolMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Rows), mArgs.Error(1)
}

func (m *PGXPoolMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Row)
}

func (m *PGXPoolMock) Close() {
	_ = m.Called()
}
