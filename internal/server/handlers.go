package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/validators"
)

type MetricResource struct {
	storageService storage.StorageService
}

func NewMetricResource(storageService storage.StorageService) MetricResource {
	return MetricResource{
		storageService: storageService,
	}
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, code int, err error) {
	logging.LogErrorCtx(ctx, err)

	w.WriteHeader(code) // only header for now
}

func (r MetricResource) Homepage(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	body := fmt.Sprintln("mainpage here.")

	records, err := r.storageService.List()
	if err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}
	if len(records) > 0 {
		body += fmt.Sprintln("metrics list:")

		for _, record := range records {
			body += fmt.Sprintf("%s => %s: %v\n", record.Name, record.Value.Kind(), record.Value)
		}
	}

	rw.Header().Set("Content-Type", "text/html")

	_, err = rw.Write([]byte(body))
	if err != nil {
		logging.LogErrorCtx(ctx, err)
	}
}

func (r MetricResource) UpdateMetric(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	mex := metrics.MetricExchange{
		ID:    req.PathValue("metricName"),
		MType: req.PathValue("metricKind"),
	}

	rawValue := req.PathValue("metricValue")

	switch mex.MType {
	case metrics.KindCounter:
		delta, err := metrics.ToCounter(rawValue)
		if err != nil {
			writeErrorResponse(ctx, rw, errToStatus(err), err)
			return
		}

		mex.Delta = &delta
	case metrics.KindGauge:
		value, err := metrics.ToGauge(rawValue)
		if err != nil {
			writeErrorResponse(ctx, rw, errToStatus(err), err)
			return
		}

		mex.Value = &value
	}

	record, err := toRecord(&mex)
	if err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}

	newRecord, err := r.storageService.Push(record)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	rw.WriteHeader(http.StatusOK)

	if _, err = io.WriteString(rw, newRecord.Value.String()); err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}
}

func (r MetricResource) UpdateMetricJSON(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	mex := new(metrics.MetricExchange)
	if err := json.NewDecoder(req.Body).Decode(mex); err != nil {
		if err == io.EOF {
			err = errors.New("no json provided")
		}

		writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
		return
	}

	record, err := toRecord(mex)
	if err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}

	newRecord, err := r.storageService.Push(record)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	mex, err = toMetricExchange(newRecord)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(rw).Encode(mex); err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}
}

func (r MetricResource) GetMetric(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	metricName := req.PathValue("metricName")
	metricKind := req.PathValue("metricKind")

	if err := validators.ValidateMetric(metricName, metricKind); err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}

	var record storage.Record
	record, err := r.storageService.Get(metricName, metricKind)
	if err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}

	body := record.Value.String()

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write([]byte(body))
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}
}

func (r MetricResource) GetMetricJSON(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	mex := new(metrics.MetricExchange)
	if err := json.NewDecoder(req.Body).Decode(mex); err != nil {
		writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
		return
	}

	if err := validators.ValidateMetric(mex.ID, mex.MType); err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}

	record, err := r.storageService.Get(mex.ID, mex.MType)
	if err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}

	mex, err = toMetricExchange(record)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(rw).Encode(mex); err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}
}

func errToStatus(err error) int {
	switch err {
	case entities.ErrRecordNotFound, entities.ErrMetricMissingName:
		return http.StatusNotFound
	case
		entities.ErrMetricUnknown, entities.ErrMetricInvalidValue,
		entities.ErrMetricInvalidName, entities.ErrMetricLongName,
		entities.ErrMetricMissingValue:

		return http.StatusBadRequest
	case entities.ErrUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
