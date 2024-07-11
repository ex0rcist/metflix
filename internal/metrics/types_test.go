package metrics_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
)

func TestKind(t *testing.T) {
	tests := []struct {
		name   string
		metric metrics.Metric
		want   string
	}{
		{
			name:   "kind = counter",
			metric: metrics.Counter(1),
			want:   "counter",
		},
		{
			name:   "kind = gauge",
			metric: metrics.Gauge(1.01),
			want:   "gauge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metric.Kind(); got != tt.want {
				t.Errorf("Metric.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name   string
		metric metrics.Metric
		want   string
	}{
		{
			name:   "String() for counter",
			metric: metrics.Counter(42),
			want:   "42",
		},
		{
			name:   "String() for gauge",
			metric: metrics.Gauge(42.01),
			want:   "42.01",
		},
		{
			name:   "String() for small gauge",
			metric: metrics.Gauge(0.42),
			want:   "0.42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metric.String(); got != tt.want {
				t.Errorf("Metric.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
