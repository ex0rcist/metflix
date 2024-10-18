package httpserver

import (
	"fmt"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/validators"
)

// from JSON-struct to storage.Record
func toRecord(mex *metrics.MetricExchange) (storage.Record, error) {
	var record storage.Record

	if err := validators.ValidateMetric(mex.ID, mex.MType); err != nil {
		return record, err
	}

	switch mex.MType {
	case metrics.KindCounter:
		if mex.Delta == nil {
			return record, entities.ErrMetricMissingValue
		}

		record = storage.Record{Name: mex.ID, Value: *mex.Delta}
	case metrics.KindGauge:
		if mex.Value == nil {
			return record, entities.ErrMetricMissingValue
		}

		record = storage.Record{Name: mex.ID, Value: *mex.Value}
	default:
		return record, entities.ErrMetricUnknown
	}

	return record, nil
}

func toMetricExchange(record storage.Record) (*metrics.MetricExchange, error) {
	req := &metrics.MetricExchange{ID: record.Name, MType: record.Value.Kind()}

	switch record.Value.Kind() {
	case metrics.KindCounter:
		delta, _ := record.Value.(metrics.Counter)
		req.Delta = &delta

	case metrics.KindGauge:
		value, _ := record.Value.(metrics.Gauge)
		req.Value = &value
	}

	return req, nil
}

func toMetricExchangeList(records []storage.Record) ([]*metrics.MetricExchange, error) {
	result := make([]*metrics.MetricExchange, len(records))

	for i, record := range records {
		req, err := toMetricExchange(record)
		if err != nil {
			return nil, fmt.Errorf("failed to convert record to MetricExchange: %w", err)
		}

		result[i] = req
	}

	return result, nil
}
