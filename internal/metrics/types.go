package metrics

import (
	"strconv"

	"github.com/ex0rcist/metflix/internal/entities"
)

const (
	KindCounter = "counter"
	KindGauge   = "gauge"
)

type Metric interface {
	Kind() string
	String() string
}

type Counter int64

func (c Counter) Kind() string {
	return KindCounter
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

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

func NewUpdateCounterMex(name string, value Counter) MetricExchange {
	return MetricExchange{ID: name, MType: value.Kind(), Delta: &value}
}

func NewUpdateGaugeMex(name string, value Gauge) MetricExchange {
	return MetricExchange{ID: name, MType: value.Kind(), Value: &value}
}

func NewGetCounterMex(name string) MetricExchange {
	return MetricExchange{ID: name, MType: KindCounter}
}

func NewGetGaugeMex(name string) MetricExchange {
	return MetricExchange{ID: name, MType: KindGauge}
}
