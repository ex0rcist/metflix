package server

import (
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
	case "counter":
		if mex.Delta == nil {
			return record, entities.ErrMetricMissingValue
		}

		record = storage.Record{Name: mex.ID, Value: *mex.Delta}
	case "gauge":
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
	case "counter":
		delta, _ := record.Value.(metrics.Counter)
		req.Delta = &delta

	case "gauge":
		value, _ := record.Value.(metrics.Gauge)
		req.Value = &value
	}

	return req, nil
}