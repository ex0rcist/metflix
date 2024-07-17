package storage_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/storage"

	"github.com/stretchr/testify/require"
)

func TestPushCounter(t *testing.T) {
	require := require.New(t)

	strg := storage.NewMemStorage()

	name := "test"
	id := storage.CalculateRecordID(name, "counter")

	value := metrics.Counter(42)
	record := storage.Record{Name: name, Value: value}
	err := strg.Push(id, record)

	require.NoError(err)
	require.Equal(value, strg.Data[record.CalculateRecordID()].Value)
}

func TestPushGauge(t *testing.T) {
	require := require.New(t)

	strg := storage.NewMemStorage()

	value := metrics.Gauge(42.42)
	record := storage.Record{Name: "test", Value: value}
	err := strg.Push(record.CalculateRecordID(), record)

	require.NoError(err)
	require.Equal(value, strg.Data[record.CalculateRecordID()].Value)
}

func TestPushWithSameName(t *testing.T) {
	require := require.New(t)

	strg := storage.NewMemStorage()

	counterValue := metrics.Counter(42)
	gaugeValue := metrics.Gauge(42.42)

	record1 := storage.Record{Name: "test", Value: counterValue}
	id1 := record1.CalculateRecordID()
	err1 := strg.Push(id1, record1)
	require.NoError(err1)

	record2 := storage.Record{Name: "test", Value: gaugeValue}
	id2 := record2.CalculateRecordID()
	err2 := strg.Push(id2, record2)
	require.NoError(err2)

	require.Equal(counterValue, strg.Data[record1.CalculateRecordID()].Value)
	require.Equal(gaugeValue, strg.Data[record2.CalculateRecordID()].Value)
}

func TestGet(t *testing.T) {
	require := require.New(t)

	strg := storage.NewMemStorage()

	value := metrics.Counter(6)
	record := storage.Record{Name: "test", Value: value}
	id := record.CalculateRecordID()
	err := strg.Push(id, record)
	require.NoError(err)

	gotRecord, err := strg.Get(record.CalculateRecordID())
	require.NoError(err)
	require.Equal(value, gotRecord.Value)
}

func TestGetNonExistantKey(t *testing.T) {
	require := require.New(t)

	strg := storage.NewMemStorage()

	_, err := strg.Get("none")
	require.Error(err)
}

func TestList(t *testing.T) {
	require := require.New(t)

	strg := storage.NewMemStorage()

	records := []storage.Record{
		{Name: "test1", Value: metrics.Counter(1)},
		{Name: "test2", Value: metrics.Counter(2)},
		{Name: "test3", Value: metrics.Gauge(3.4)},
	}

	for _, r := range records {
		err := strg.Push(r.CalculateRecordID(), r)
		require.NoError(err)
	}

	allRecords, err := strg.List()

	require.NoError(err)
	require.ElementsMatch(records, allRecords)
}
