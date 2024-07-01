package metrics

type Metric interface {
	Kind() string
}

type Counter int64

func (c Counter) Kind() string {
	return "counter"
}

type Gauge float64

func (g Gauge) Kind() string {
	return "gauge"
}
