package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/validators"
)

type Resource struct {
	storage storage.Storage
}

func NewResource(s storage.Storage) Resource {
	return Resource{
		storage: s,
	}
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, code int, err error) {
	logging.LogError(ctx, err)

	w.WriteHeader(code) // only header for now

	// resp := fmt.Sprintf("%d %v", code, err)
	// http.Error(w, resp, code)
}

func (r Resource) Homepage(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	body := fmt.Sprintln("mainpage here.")

	// todo: errors
	records, _ := r.storage.GetAll()
	if len(records) > 0 {
		body += fmt.Sprintln("metrics list:")

		for _, record := range records {
			body += fmt.Sprintf("%s => %s: %v\n", record.Name, record.Value.Kind(), record.Value)
		}
	}

	_, err := rw.Write([]byte(body))
	if err != nil {
		logging.LogError(ctx, err)
	}
}

func (r Resource) UpdateMetric(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	// logger := zerolog.Ctx(req.Context())

	metricName := req.PathValue("metricName")
	metricKind := req.PathValue("metricKind")
	metricValue := req.PathValue("metricValue")

	//logger.Info().Msg(fmt.Sprintf("updating1 metric... %v", metricName))
	logging.LogInfo(ctx, fmt.Sprintf("updating metric... %v", metricName))

	if err := validators.EnsureNamePresent(metricName); err != nil {
		writeErrorResponse(ctx, rw, http.StatusNotFound, err)
		return
	}

	if err := validators.ValidateName(metricName); err != nil {
		writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
		return
	}

	if err := validators.ValidateKind(metricKind); err != nil {
		writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
		return
	}

	var rec storage.Record

	switch metricKind {
	case "counter":
		var currentValue int64

		recordID := storage.CalculateRecordID(metricName, metricKind)
		current, err := r.storage.Get(recordID)

		if err != nil && errors.Is(err, entities.ErrMetricNotFound) {
			currentValue = 0
		} else if err != nil {
			writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
			return
		} else {
			currentValue = int64(current.Value.(metrics.Counter))
		}

		incr, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
			return
		}
		currentValue += incr

		rec = storage.Record{Name: metricName, Value: metrics.Counter(currentValue)}
	case "gauge":
		current, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
			return
		}

		rec = storage.Record{Name: metricName, Value: metrics.Gauge(current)}
	default:
		writeErrorResponse(ctx, rw, http.StatusBadRequest, entities.ErrMetricUnknown)
		return
	}

	err := r.storage.Push(rec)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusInternalServerError, err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (r Resource) ShowMetric(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	metricName := req.PathValue("metricName")
	metricKind := req.PathValue("metricKind")

	if err := validators.EnsureNamePresent(metricName); err != nil {
		writeErrorResponse(ctx, rw, http.StatusNotFound, err)
		return
	}

	if err := validators.ValidateName(metricName); err != nil {
		writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
		return
	}

	if err := validators.ValidateKind(metricKind); err != nil {
		writeErrorResponse(ctx, rw, http.StatusBadRequest, err)
		return
	}

	var record storage.Record

	recordID := storage.CalculateRecordID(metricName, metricKind)
	record, err := r.storage.Get(recordID)
	if err != nil {
		writeErrorResponse(ctx, rw, http.StatusNotFound, err)
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
