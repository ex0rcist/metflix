package storage_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestRecordId(t *testing.T) {
	require := require.New(t)

	record1 := storage.Record{Name: "test", Value: metrics.Counter(42)}
	record2 := storage.Record{Name: "test", Value: metrics.Gauge(42.42)}

	require.Equal("test_counter", record1.CalculateRecordID())
	require.Equal("test_gauge", record2.CalculateRecordID())
}

func TestRecordIdWithEmptyName(t *testing.T) {
	require := require.New(t)
	require.Equal("", storage.CalculateRecordID("", "counter"))
}

func TestRecordIdWithEmptyKind(t *testing.T) {
	require := require.New(t)
	require.Equal("", storage.CalculateRecordID("test", ""))
}
