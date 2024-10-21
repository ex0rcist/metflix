package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/profiler"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/validators"
	"github.com/ex0rcist/metflix/pkg/metrics"
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

// Homepage godoc
// @Tags Metrics
// @Router / [get]
// @Summary Yet another homepage
// @ID homepage
// @Produce text/html
// @Success 200 {string} string
// @Failure 500 {string} string http.StatusInternalServerError
func (r MetricResource) Homepage(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	body := fmt.Sprintln("mainpage here.")

	records, err := r.storageService.List(ctx)
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

// UpdateMetric godoc
// @Tags Metrics
// @Router /update/{type}/{name}/{value} [post]
// @Summary Push metric data.
// @ID metrics_update
// @Produce plain
// @Param type path string true "Metrics type (e.g. `counter`, `gauge`)."
// @Param name path string true "Metrics name."
// @Param value path string true "Metrics value, must be convertable to `int64` or `float64`."
// @Success 200 {string} string
// @Failure 400 {string} string http.StatusBadRequest
// @Failure 500 {string} string http.StatusInternalServerError
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

	newRecord, err := r.storageService.Push(ctx, record)
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

// UpdateMetricJSON godoc
// @Tags Metrics
// @Router /update [post]
// @Summary Push metric data as JSON
// @ID metrics_json_update
// @Accept  json
// @Param request body metrics.MetricExchange true "Request parameters."
// @Success 200 {object} metrics.MetricExchange
// @Failure 400 {string} string http.StatusBadRequest
// @Failure 500 {string} string http.StatusInternalServerError
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

	newRecord, err := r.storageService.Push(ctx, record)
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

// BatchUpdateMetricsJSON godoc
// @Tags Metrics
// @Router /updates [post]
// @Summary Push list of metrics data as JSON
// @ID metrics_json_update_list
// @Accept  json
// @Param request body []metrics.MetricExchange true "List of metrics to update."
// @Success 200 {object} []metrics.MetricExchange
// @Failure 400 {string} string http.StatusBadRequest
// @Failure 500 {string} string http.StatusInternalServerError
func (r MetricResource) BatchUpdateMetricsJSON(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	records, err := parseJSONMetricsList(req)
	if err != nil {
		if err == io.EOF {
			err = errors.New("no json provided")
		}

		writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
		return
	}

	recorded, err := r.storageService.PushList(ctx, records)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	resp, err := toMetricExchangeList(recorded)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	profiler.GetProfiler().SaveMemoryProfile()
}

// GetMetric godoc
// @Tags Metrics
// @Router /value/{type}/{name} [get]
// @Summary Get metric's value as string
// @ID metrics_info
// @Produce plain
// @Param type path string true "Metrics type (e.g. `counter`, `gauge`)."
// @Param name path string true "Metrics name."
// @Success 200 {string} string
// @Failure 400 {string} string http.StatusBadRequest
// @Failure 404 {string} string http.StatusNotFound
// @Failure 500 {string} string http.StatusInternalServerError
func (r MetricResource) GetMetric(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	metricName := req.PathValue("metricName")
	metricKind := req.PathValue("metricKind")

	if err := validators.ValidateMetric(metricName, metricKind); err != nil {
		writeErrorResponse(ctx, rw, errToStatus(err), err)
		return
	}

	var record storage.Record
	record, err := r.storageService.Get(ctx, metricName, metricKind)
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

// GetMetricJSON godoc
// @Tags Metrics
// @Router /value [post]
// @Summary Get metrics value as JSON
// @ID metrics_json_info
// @Accept  json
// @Produce json
// @Param request body metrics.MetricExchange true "Request parameters: `id` and `type` are required."
// @Success 200 {object} metrics.MetricExchange
// @Failure 400 {string} string http.StatusBadRequest
// @Failure 404 {string} string http.StatusNotFound
// @Failure 500 {string} string http.StatusInternalServerError
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

	record, err := r.storageService.Get(ctx, mex.ID, mex.MType)
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

func parseJSONMetricsList(r *http.Request) ([]storage.Record, error) {
	req := make([]metrics.MetricExchange, 0)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	records := make([]storage.Record, len(req))

	for i := range req {
		record, err := toRecord(&req[i])
		if err != nil {
			return nil, err
		}

		records[i] = record
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no metrics provided")
	}

	return records, nil
}

type PingerResource struct {
	pinger services.Pinger
}

func NewPingerResource(pinger services.Pinger) PingerResource {
	return PingerResource{
		pinger: pinger,
	}
}

// Ping godoc
// @Tags Healthcheck
// @Router /ping [get]
// @Summary Verify server up and running
// @ID health_info
// @Success 200
// @Failure 500 {string} string http.StatusInternalServerError
// @Failure 501 {string} string http.StatusNotImplemented
func (pr PingerResource) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := pr.pinger.Ping(ctx)
	if err == nil {
		return
	}

	if errors.Is(err, entities.ErrStorageUnpingable) {
		writeErrorResponse(ctx, w, http.StatusNotImplemented, err)
		return
	}

	writeErrorResponse(ctx, w, http.StatusInternalServerError, err)
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
