package server

import (
	"net/http"

	"github.com/ex0rcist/metflix/internal/storage"
)

func NewRouter(storage storage.Storage) http.Handler {
	mux := http.NewServeMux()
	resource := Resource{storage: storage}

	mux.HandleFunc("GET /update/{metricType}/{metricName}/{metricValue}", resource.UpdateMetric)
	mux.HandleFunc(`/`, resource.Homepage)

	return mux
}
