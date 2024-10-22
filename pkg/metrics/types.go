package metrics

import (
	"strconv"

	"github.com/ex0rcist/metflix/internal/entities"
)

// Available metric types.
const (
	KindCounter = "counter"
	KindGauge   = "gauge"
)

// Metric interface
type Metric interface {
	Kind() string
	String() string
}

// Gauge metric type - int64.
type Counter int64

// Return metric kind
func (c Counter) Kind() string {
	return KindCounter
}

// Stringer
func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

// Gauge metric type - float64.
type Gauge float64

// Return metric kind
func (g Gauge) Kind() string {
	return KindGauge
}

// Stringer
func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

// Convert string value to metrics.Counter
func ToCounter(value string) (Counter, error) {
	rawValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, entities.ErrMetricInvalidValue
	}

	return Counter(rawValue), nil
}

// Convert string value to metrics.Gauge
func ToGauge(value string) (Gauge, error) {
	rawValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, entities.ErrMetricInvalidValue
	}

	return Gauge(rawValue), nil
}

// Agent/Server exchange schema according to spec
type MetricExchange struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *Counter `json:"delta,omitempty"`
	Value *Gauge   `json:"value,omitempty"`
}
