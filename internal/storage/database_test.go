// storage_test.go
package storage

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pashagolub/pgxmock/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDatabaseStorage_Push(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	storage := DatabaseStorage{Pool: mockPool}

	ctx := context.Background()
	record := Record{Name: "testName", Value: metrics.Counter(123)}
	key := record.CalculateRecordID()

	mockPool.ExpectBegin()
	mockPool.ExpectExec("INSERT INTO metrics").WithArgs(key, record.Name, record.Value.Kind(), record.Value.String()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mockPool.ExpectCommit()

	err = storage.Push(ctx, key, record)
	require.NoError(t, err)

	if err := mockPool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDatabaseStorage_PushList(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := DatabaseStorage{Pool: mockPool}

	ctx := context.Background()
	data := map[string]Record{
		"name1_counter": {Name: "name1", Value: metrics.Counter(123)},
		"name2_gauge":   {Name: "name2", Value: metrics.Gauge(456)},
	}

	mockBatchResults := new(PGXBatchResultsMock)
	mockPool.On("SendBatch", ctx, mock.Anything).Return(mockBatchResults)
	mockBatchResults.On("Exec").Return(pgconn.CommandTag{}, nil).Twice()
	mockBatchResults.On("Close").Return(nil)

	err := storage.PushList(ctx, data)
	assert.NoError(t, err)

	mockPool.AssertExpectations(t)
	mockBatchResults.AssertExpectations(t)
}

func TestDatabaseStorage_Get(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := DatabaseStorage{Pool: mockPool}

	ctx := context.Background()
	expectedRecord := Record{Name: "testName", Value: metrics.Counter(123)}
	key := expectedRecord.CalculateRecordID()

	mockRow := new(PGXRowMock)
	mockPool.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything).Run(func(mArgs mock.Arguments) {
		*mArgs.Get(0).(*string) = expectedRecord.Name
		*mArgs.Get(1).(*string) = expectedRecord.Value.Kind()
		*mArgs.Get(2).(*float64) = 123
	}).Return(nil)

	record, err := storage.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord, record)

	mockPool.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestDatabaseStorage_Ping(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := DatabaseStorage{Pool: mockPool}

	ctx := context.Background()
	mockPool.On("Ping", ctx).Return(nil)

	err := storage.Ping(ctx)
	assert.NoError(t, err)

	mockPool.AssertExpectations(t)
}

func TestDatabaseStorage_Close(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := DatabaseStorage{Pool: mockPool}

	ctx := context.Background()

	mockPool.On("Close").Return(nil)

	storage.Close(ctx)
	mockPool.AssertExpectations(t)
}
