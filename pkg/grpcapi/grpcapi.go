package grpcapi

import "github.com/ex0rcist/metflix/pkg/metrics"

func NewUpdateCounterMex(name string, value metrics.Counter) *MetricExchange {
	return &MetricExchange{Id: name, Mtype: value.Kind(), Delta: int64(value)}
}

func NewUpdateGaugeMex(name string, value metrics.Gauge) *MetricExchange {
	return &MetricExchange{Id: name, Mtype: value.Kind(), Value: float64(value)}
}
