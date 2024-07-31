package storage

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_Push(t *testing.T) {
	ctx := context.Background()
	strg := NewMemStorage()

	records := []Record{
		{Name: metrics.KindCounter, Value: metrics.Counter(42)},
		{Name: metrics.KindGauge, Value: metrics.Gauge(42.42)},
	}

	for _, r := range records {
		id := r.CalculateRecordID()

		err := strg.Push(ctx, id, r)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if s, _ := strg.Get(ctx, id); r != s {
			t.Fatalf("expected record %v, got %v", r, s)
		}
	}
}

func TestMemStorage_Push_WithSameName(t *testing.T) {
	ctx := context.Background()
	strg := NewMemStorage()

	counterValue := metrics.Counter(42)
	gaugeValue := metrics.Gauge(42.42)

	records := []Record{
		{Name: "test", Value: counterValue},
		{Name: "test", Value: gaugeValue},
	}

	for _, r := range records {
		id := r.CalculateRecordID()

		if err := strg.Push(ctx, id, r); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	storedCounter, err := strg.Get(ctx, records[0].CalculateRecordID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storedGauge, err := strg.Get(ctx, records[1].CalculateRecordID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if storedCounter != records[0] {
		t.Fatalf("expected stored %v, got %v", records[0], storedCounter)
	}

	if storedGauge != records[1] {
		t.Fatalf("expected stored %v, got %v", records[1], storedGauge)
	}
}

func TestMemStorage_Get(t *testing.T) {
	ctx := context.Background()
	strg := NewMemStorage()
	record := Record{Name: "1", Value: metrics.Counter(42)}
	err := strg.Push(ctx, record.CalculateRecordID(), record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		name      string
		id        string
		want      Record
		wantError bool
	}{
		{name: "existing record", id: record.CalculateRecordID(), want: record, wantError: false},
		{name: "non-existing record", id: "test", want: Record{}, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strg.Get(ctx, tt.id)
			if (err != nil) != tt.wantError {
				t.Fatalf("expected error: %v, got %v", tt.wantError, err)
			}
			if got != tt.want {
				t.Fatalf("expected record %v, got %v", tt.want, got)
			}
		})
	}
}

func TestMemStorage_List(t *testing.T) {
	ctx := context.Background()
	storage := NewMemStorage()

	records := []Record{
		{Name: metrics.KindGauge, Value: metrics.Gauge(42.42)},
		{Name: metrics.KindCounter, Value: metrics.Counter(42)},
	}

	err := storage.Push(ctx, records[0].CalculateRecordID(), records[0])
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = storage.Push(ctx, records[1].CalculateRecordID(), records[1])
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := storage.List(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(got) != len(records) {
		t.Fatalf("expected %d records, got %d", len(records), len(got))
	}

	require.ElementsMatch(t, records, got)
}
