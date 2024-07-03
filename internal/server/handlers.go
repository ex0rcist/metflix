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

func NewResource(s storage.Storage) Resource {
	return Resource{
		storage: s,
	}
}

func (r Resource) Homepage(res http.ResponseWriter, _ *http.Request) {
	body := fmt.Sprintln("mainpage here.")

	// todo: errors
	records, _ := r.storage.GetAll()
	if len(records) > 0 {
		body += fmt.Sprintln("metrics list:")

		for _, record := range records {
			body += fmt.Sprintf("%s => %s: %v\n", record.Name, record.Value.Kind(), record.Value)
		}
	}

	_, err := res.Write([]byte(body))
	if err != nil {
		panic("unable to write") // todo
	}
}

func (r Resource) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	metricName := req.PathValue("metricName")

	metricType := req.PathValue("metricType")
	metricValue := req.PathValue("metricValue")

	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if err := validators.EnsureNamePresent(metricName); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	//При попытке передать запрос с некорректным именем метрики возвращать http.StatusBadRequest.
	if err := validators.ValidateName(metricName); err != nil {
		res.WriteHeader(http.StatusBadRequest)
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

		incr, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		value += incr

		record = storage.Record{Name: metricName, Value: metrics.Counter(value)}
	case "gauge":
		current, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		record = storage.Record{Name: metricName, Value: metrics.Gauge(current)}
	default:
		panic("why") // todo
	}

	err := r.storage.Push(record)
	if err != nil {
		panic("cannot push") // todo
	}

	res.WriteHeader(http.StatusOK)
}
