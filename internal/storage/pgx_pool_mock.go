package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

var _ PGXPool = (*PGXPoolMock)(nil)

type PGXPoolMock struct {
	mock.Mock
}

func NewPGXPoolMock() *PGXPoolMock {
	return new(PGXPoolMock)
}

func (m *PGXPoolMock) Begin(ctx context.Context) (pgx.Tx, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(pgx.Tx), mArgs.Error(1)
}

func (m *PGXPoolMock) Acquire(ctx context.Context) (c *pgxpool.Conn, err error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(*pgxpool.Conn), mArgs.Error(1)
}

func (m *PGXPoolMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}

func (m *PGXPoolMock) SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults) {
	mArgs := m.Called(ctx, b)
	return mArgs.Get(0).(pgx.BatchResults)
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

// ************** PGXBatchResultsMock ************** //

type PGXBatchResultsMock struct {
	mock.Mock
}

func (m *PGXBatchResultsMock) Exec() (pgconn.CommandTag, error) {
	mArgs := m.Called()
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}

func (m *PGXBatchResultsMock) Query() (pgx.Rows, error) {
	mArgs := m.Called()
	return mArgs.Get(0).(pgx.Rows), mArgs.Error(1)
}

func (m *PGXBatchResultsMock) QueryRow() pgx.Row {
	mArgs := m.Called()
	return mArgs.Get(0).(pgx.Row)
}

func (m *PGXBatchResultsMock) Close() error {
	mArgs := m.Called()
	return mArgs.Error(0)
}

// ************** PGXRowMock ************** //

type PGXRowMock struct {
	mock.Mock
}

func (m *PGXRowMock) Scan(args ...any) error {
	mArgs := m.Called(args...)
	return mArgs.Error(0)
}

// ************** PGXRowsMock ************** //

type PGXRowsMock struct {
	mock.Mock
}

func (m *PGXRowsMock) Close() {
	_ = m.Called()
}

func (m *PGXRowsMock) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *PGXRowsMock) Scan(args ...any) error {
	mArgs := m.Called(args...)
	return mArgs.Error(0)
}

func (m *PGXRowsMock) Err() error {
	mArgs := m.Called()
	return mArgs.Error(0)
}

func (m *PGXRowsMock) CommandTag() pgconn.CommandTag {
	mArgs := m.Called()
	return mArgs.Get(0).(pgconn.CommandTag)
}

func (m *PGXRowsMock) FieldDescriptions() []pgconn.FieldDescription {
	mArgs := m.Called()
	return mArgs.Get(0).([]pgconn.FieldDescription)
}

func (m *PGXRowsMock) Values() ([]any, error) {
	mArgs := m.Called()
	return mArgs.Get(0).([]any), mArgs.Error(1)
}

func (m *PGXRowsMock) RawValues() [][]byte {
	mArgs := m.Called()
	return mArgs.Get(0).([][]byte)
}

func (m *PGXRowsMock) Conn() *pgx.Conn {
	mArgs := m.Called()
	return mArgs.Get(0).(*pgx.Conn)
}
