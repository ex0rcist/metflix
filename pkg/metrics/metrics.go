// Package metrics provides client REST API for metrics collector (server).
package metrics

// Create new MetricExchange struct with metrics.Counter value to be used for updating metrics.
func NewUpdateCounterMex(name string, value Counter) MetricExchange {
	return MetricExchange{ID: name, MType: value.Kind(), Delta: &value}
}

// Create new MetricExchange struct with metrics.Gauge value to be used for updating metrics.
func NewUpdateGaugeMex(name string, value Gauge) MetricExchange {
	return MetricExchange{ID: name, MType: value.Kind(), Value: &value}
}

// Create new MetMetricExchange struct to be used for retrieving of counter metric.
func NewGetCounterMex(name string) MetricExchange {
	return MetricExchange{ID: name, MType: KindCounter}
}

// Create new MetMetricExchange struct to be used for retrieving of gauge metric.
func NewGetGaugeMex(name string) MetricExchange {
	return MetricExchange{ID: name, MType: KindGauge}
}
