// storage_test.go
package storage

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/pkg/metrics"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgresStorage_Push(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := PostgresStorage{Pool: mockPool}

	ctx := context.Background()
	record := Record{Name: "testName", Value: metrics.Counter(123)}
	key := record.CalculateRecordID()

	txMock := new(PGXTxMock)
	mockPool.On("Begin", mock.Anything).Return(txMock, nil)
	txMock.
		On("Exec", mock.Anything, mock.Anything, key, record.Name, record.Value.Kind(), record.Value.String()).
		Return(pgconn.CommandTag{}, nil)

	txMock.On("Commit", mock.Anything).Return(nil)

	err := storage.Push(ctx, key, record)
	if err != nil {
		t.Fatalf("expected no error on Push, got: %v", err)
	}
}

func TestPostgresStorage_PushList(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := PostgresStorage{Pool: mockPool}

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

func TestPostgresStorage_Get(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := PostgresStorage{Pool: mockPool}

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

func TestPostgresStorage_List(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := PostgresStorage{Pool: mockPool}

	ctx := context.Background()
	expectedRecords := []Record{
		{Name: "name1", Value: metrics.Counter(123)},
		{Name: "name2", Value: metrics.Gauge(456)},
	}

	mockRows := new(PGXRowsMock)
	mockPool.On("Query", ctx, mock.AnythingOfType("string"), []interface{}(nil)).Return(mockRows, nil)
	mockRows.On("Next").Return(true).Twice()
	mockRows.On("Next").Return(false)

	counter := 0
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		rec := expectedRecords[counter]
		*args.Get(0).(*string) = rec.Name
		*args.Get(1).(*string) = rec.Value.Kind()

		switch expectedRecords[counter].Value.Kind() {
		case metrics.KindCounter:
			value, _ := rec.Value.(metrics.Counter)
			*args.Get(2).(*float64) = float64(value)
		case metrics.KindGauge:
			value, _ := rec.Value.(metrics.Gauge)
			*args.Get(2).(*float64) = float64(value)
		}

		counter++
	}).Twice().Return(nil)
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("CommandTag").Return(pgconn.NewCommandTag("select"))

	records, err := storage.List(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecords, records)

	mockPool.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestPostgresStorage_Ping(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := PostgresStorage{Pool: mockPool}

	ctx := context.Background()
	mockPool.On("Ping", ctx).Return(nil)

	err := storage.Ping(ctx)
	assert.NoError(t, err)

	mockPool.AssertExpectations(t)
}

func TestPostgresStorage_Close(t *testing.T) {
	mockPool := NewPGXPoolMock()
	storage := PostgresStorage{Pool: mockPool}

	ctx := context.Background()

	mockPool.On("Close").Return(nil)

	storage.Close(ctx)
	mockPool.AssertExpectations(t)
}
