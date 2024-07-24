package storage

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
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
