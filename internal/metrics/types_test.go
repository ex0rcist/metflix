package metrics

import (
	"testing"
)

func TestKind(t *testing.T) {
	tests := []struct {
		name   string
		metric Metric
		want   string
	}{
		{
			name:   "kind = counter",
			metric: Counter(1),
			want:   KindCounter,
		},
		{
			name:   "kind = gauge",
			metric: Gauge(1.01),
			want:   KindGauge,
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
		metric Metric
		want   string
	}{
		{
			name:   "String() for counter",
			metric: Counter(42),
			want:   "42",
		},
		{
			name:   "String() for gauge",
			metric: Gauge(42.01),
			want:   "42.01",
		},
		{
			name:   "String() for small gauge",
			metric: Gauge(0.42),
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
