package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/rs/zerolog/log"
)

type API struct {
	address    *entities.Address
	httpClient *http.Client
}

func NewAPI(address *entities.Address, httpTransport http.RoundTripper) *API {
	if httpTransport == nil {
		httpTransport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: httpTransport,
	}

	return &API{
		address:    address,
		httpClient: client,
	}
}

func (c *API) Report(name string, metric metrics.Metric) *API {
	ctx := context.Background()

	// todo: another transport?
	url := "http://" + c.address.String() + "/update/"

	var mex metrics.MetricExchange

	// HELP: можно ли тут вместо приведения типов
	// использовать рефлексию через metric.(type) ?
	switch metric.Kind() {
	case "counter":
		mex = metrics.NewUpdateCounterMex(name, metric.(metrics.Counter))
	case "gauge":
		mex = metrics.NewUpdateGaugeMex(name, metric.(metrics.Gauge))
	default:
		log.Warn().Msg("unknown metric") // todo
	}

	body, err := json.Marshal(mex)
	if err != nil {
		log.Warn().Msg("unknown metric") // todo
	}

	log.Info().Str("target", url).Str("payload", string(body)).Msg("sending report")

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logging.LogError(ctx, err, "httpRequest error")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Warn().Msg("httpClient error") // todo
		return c
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body) // нужно прочитать ответ для keepalive?
	if err != nil {
		logging.LogError(ctx, entities.ErrMetricReport, "error reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		logging.LogError(ctx, entities.ErrMetricReport, string(respBody))
	}

	return c
}
