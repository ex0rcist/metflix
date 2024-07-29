package storage

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/stretchr/testify/require"
)

func TestCalculateRecordID(t *testing.T) {
	tests := []struct {
		test     string
		name     string
		kind     string
		expected string
	}{
		{test: "valid inputs", name: "metricName", kind: "metricKind", expected: "metricName_metricKind"},
		{test: "empty name", name: "", kind: "metricKind", expected: ""},
		{test: "empty kind", name: "metricName", kind: "", expected: ""},
		{test: "both empty", name: "", kind: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateRecordID(tt.name, tt.kind)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRecord_CalculateRecordID(t *testing.T) {
	tests := []struct {
		name     string
		record   Record
		expected string
	}{
		{name: "valid record with counter", record: Record{Name: "metricName", Value: metrics.Counter(100)}, expected: "metricName_counter"},
		{name: "valid record with gauge", record: Record{Name: "metricName", Value: metrics.Gauge(100.0)}, expected: "metricName_gauge"},
		{name: "empty name", record: Record{Name: "", Value: metrics.Counter(100)}, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.record.CalculateRecordID()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRecordMarshalingNormalData(t *testing.T) {
	tests := []struct {
		name   string
		source Record
	}{
		{name: "Should convert counter", source: Record{Name: "PollCount", Value: metrics.Counter(10)}},
		{name: "Should convert gauge", source: Record{Name: "Alloc", Value: metrics.Gauge(42.0)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json, err := tt.source.MarshalJSON()
			if err != nil {
				t.Fatalf("expected no error marshaling json, got: %v", err)
			}

			target := new(Record)
			err = target.UnmarshalJSON(json)
			if err != nil {
				t.Fatalf("expected no error unmarshaling json, got: %v", err)
			}

			if tt.source != *target {
				t.Fatal("expected records to be equal:", tt.source, target)
			}
		})
	}
}

func TestRecordUnmarshalingCorruptedData(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected error
	}{
		{name: "Should fail on broken json", data: `{"name": "xxx",`, expected: &json.SyntaxError{}},
		{name: "Should fail on invalid counter", data: `{"name": "xxx", "kind": "counter", "value": "12.345"}`, expected: strconv.ErrSyntax},
		{name: "Should fail on invalid gauge", data: `{"name": "xxx", "kind": "gauge", "value": "12.)"}`, expected: strconv.ErrSyntax},
		{name: "Should fail on unknown kind", data: `{"name": "xxx", "kind": "unknown", "value": "12"}`, expected: entities.ErrMetricUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(Record)

			err := r.UnmarshalJSON([]byte(tt.data))
			require.ErrorAs(t, err, &tt.expected)
		})
	}
}
