package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/validators"
)

type Resource struct {
	storage storage.Storage
}

func (r Resource) Homepage(res http.ResponseWriter, _ *http.Request) {
	body := fmt.Sprintln("mainpage here.", "metrics list: ")

	// todo: errors
	records, _ := r.storage.GetAll()
	for _, record := range records {
		body += fmt.Sprintf("%s => %s: %v\r\n", record.Name, record.Value.Kind(), record.Value)
	}

	res.Write([]byte(body))
}

func (r Resource) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	metricName := req.PathValue("metricName")
	metricType := req.PathValue("metricType")
	metricValue := req.PathValue("metricValue")

	if err := validators.EnsureNamePresent(metricName); err != nil {
		res.WriteHeader(http.StatusOK)
		return
	}

	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if err := validators.EnsureNamePresent(metricName); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
	if err := validators.ValidateKind(metricType); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var record storage.Record

	switch metricType {
	case "counter":
		var value int64 = 0

		current, err := r.storage.Get(metricName)
		if err == nil {
			value = int64(current.Value.(metrics.Counter))
		}

		incr, _ := strconv.ParseInt(metricValue, 10, 64)
		value += incr

		record = storage.Record{Name: metricName, Value: metrics.Counter(value)}
	case "gauge":
		current, _ := strconv.ParseFloat(metricValue, 64)
		record = storage.Record{Name: metricName, Value: metrics.Gauge(current)}
	default:
		panic("why") // todo
	}

	r.storage.Push(record)

	res.WriteHeader(http.StatusOK)
}
