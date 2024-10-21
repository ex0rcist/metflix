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

type Metric interface {
	Kind() string
	String() string
}

// Gauge metric type - int64.
type Counter int64

func (c Counter) Kind() string {
	return KindCounter
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

// Gauge metric type - float64.
type Gauge float64

func (g Gauge) Kind() string {
	return KindGauge
}

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func ToCounter(value string) (Counter, error) {
	rawValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, entities.ErrMetricInvalidValue
	}

	return Counter(rawValue), nil
}

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
