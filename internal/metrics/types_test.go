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
