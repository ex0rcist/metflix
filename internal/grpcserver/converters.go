package grpcserver

import (
	"fmt"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/validators"
	"github.com/ex0rcist/metflix/pkg/grpcapi"
	"github.com/ex0rcist/metflix/pkg/metrics"
)

func toRecord(req *grpcapi.MetricExchange) (storage.Record, error) {
	var record storage.Record

	if err := validators.ValidateMetric(req.Id, req.Mtype); err != nil {
		return record, err
	}

	switch req.Mtype {
	case metrics.KindCounter:
		record = storage.Record{Name: req.Id, Value: metrics.Counter(req.Delta)}
	case metrics.KindGauge:
		record = storage.Record{Name: req.Id, Value: metrics.Gauge(req.Value)}
	default:
		return record, entities.ErrMetricUnknown
	}

	return record, nil
}

func toMetricExchange(record storage.Record) (*grpcapi.MetricExchange, error) {
	req := &grpcapi.MetricExchange{
		Id:    record.Name,
		Mtype: record.Value.Kind(),
	}

	switch record.Value.Kind() {
	case metrics.KindCounter:
		delta, _ := record.Value.(metrics.Counter)
		req.Delta = int64(delta)

	case metrics.KindGauge:
		value, _ := record.Value.(metrics.Gauge)
		req.Value = float64(value)
	}

	return req, nil
}

func toRecordsList(req *grpcapi.BatchUpdateRequest) ([]storage.Record, error) {
	rv := make([]storage.Record, len(req.Data))

	for i := range req.Data {
		record, err := toRecord(req.Data[i])
		if err != nil {
			return nil, err
		}

		rv[i] = record
	}

	if len(rv) == 0 {
		return nil, entities.ErrMetricBatchIncomplete
	}

	return rv, nil
}

func toMetricExchangeList(records []storage.Record) ([]*grpcapi.MetricExchange, error) {
	rv := make([]*grpcapi.MetricExchange, len(records))

	for i, record := range records {
		req, err := toMetricExchange(record)
		if err != nil {
			return nil, fmt.Errorf("toMetricExchangeList - toMetricExchange: %w", err)
		}

		rv[i] = req
	}

	return rv, nil
}
