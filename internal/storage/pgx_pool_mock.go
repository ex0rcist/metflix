package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

var _ PGXPool = (*PGXPoolMock)(nil)

// PGX.pool mock
type PGXPoolMock struct {
	mock.Mock
}

// Constructor
func NewPGXPoolMock() *PGXPoolMock {
	return new(PGXPoolMock)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) Begin(ctx context.Context) (pgx.Tx, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(pgx.Tx), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) Acquire(ctx context.Context) (c *pgxpool.Conn, err error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(*pgxpool.Conn), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults) {
	mArgs := m.Called(ctx, b)
	return mArgs.Get(0).(pgx.BatchResults)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Rows), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Row)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXPoolMock) Close() {
	_ = m.Called()
}

// ************** PGXBatchResultsMock ************** //

// PGXbatch result mock
type PGXBatchResultsMock struct {
	mock.Mock
}

// Stub comment, required for linter. See original for comment.
func (m *PGXBatchResultsMock) Exec() (pgconn.CommandTag, error) {
	mArgs := m.Called()
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXBatchResultsMock) Query() (pgx.Rows, error) {
	mArgs := m.Called()
	return mArgs.Get(0).(pgx.Rows), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXBatchResultsMock) QueryRow() pgx.Row {
	mArgs := m.Called()
	return mArgs.Get(0).(pgx.Row)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXBatchResultsMock) Close() error {
	mArgs := m.Called()
	return mArgs.Error(0)
}

// ************** PGXRowMock ************** //

// PGX.row mock
type PGXRowMock struct {
	mock.Mock
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowMock) Scan(args ...any) error {
	mArgs := m.Called(args...)
	return mArgs.Error(0)
}

// ************** PGXRowsMock ************** //

// PGX.rows mock
type PGXRowsMock struct {
	mock.Mock
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) Close() {
	_ = m.Called()
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) Scan(args ...any) error {
	mArgs := m.Called(args...)
	return mArgs.Error(0)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) Err() error {
	mArgs := m.Called()
	return mArgs.Error(0)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) CommandTag() pgconn.CommandTag {
	mArgs := m.Called()
	return mArgs.Get(0).(pgconn.CommandTag)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) FieldDescriptions() []pgconn.FieldDescription {
	mArgs := m.Called()
	return mArgs.Get(0).([]pgconn.FieldDescription)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) Values() ([]any, error) {
	mArgs := m.Called()
	return mArgs.Get(0).([]any), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) RawValues() [][]byte {
	mArgs := m.Called()
	return mArgs.Get(0).([][]byte)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXRowsMock) Conn() *pgx.Conn {
	mArgs := m.Called()
	return mArgs.Get(0).(*pgx.Conn)
}

// ************** PGXTxMock ************** //

// PGX.tx mock
type PGXTxMock struct {
	mock.Mock
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) Begin(ctx context.Context) (pgx.Tx, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(pgx.Tx), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) Commit(ctx context.Context) error {
	mArgs := m.Called(ctx)
	return mArgs.Error(0)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) Rollback(ctx context.Context) error {
	mArgs := m.Called(ctx)
	return mArgs.Error(0)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	mArgs := m.Called(ctx, tableName)
	return mArgs.Get(0).(int64), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	mArgs := m.Called(ctx, b)
	return mArgs.Get(0).(pgx.BatchResults)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) LargeObjects() pgx.LargeObjects {
	mArgs := m.Called()
	return mArgs.Get(0).(pgx.LargeObjects)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	mArgs := m.Called(ctx, name, sql)
	return mArgs.Get(0).(*pgconn.StatementDescription), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error) {
	varargs := append([]any{ctx, sql}, args...)
	mArgs := m.Called(varargs...)
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Rows), mArgs.Error(1)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Row)
}

// Stub comment, required for linter. See original for comment.
func (m *PGXTxMock) Conn() *pgx.Conn {
	mArgs := m.Called()
	return mArgs.Get(0).(*pgx.Conn)
}
