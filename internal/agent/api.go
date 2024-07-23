package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ex0rcist/metflix/internal/compression"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/utils"
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
	url := "http://" + c.address.String() + "/update"

	var requestID = utils.GenerateRequestID()
	var ctx = setupLoggerCtx(requestID)
	var mex metrics.MetricExchange

	// HELP: можно ли тут вместо приведения типов
	// использовать рефлексию через metric.(type) ?
	switch metric.Kind() {
	case "counter":
		mex = metrics.NewUpdateCounterMex(name, metric.(metrics.Counter))
	case "gauge":
		mex = metrics.NewUpdateGaugeMex(name, metric.(metrics.Gauge))
	default:
		logging.LogError(entities.ErrMetricReport, "unknown metric")
		return c
	}

	body, err := json.Marshal(mex)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during marshaling", err.Error())
		return c
	}

	payload, err := compression.Pack(body)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during compression", err.Error())
		return c
	}

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "httpRequest error", err.Error())
		return c
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("X-Request-Id", requestID)

	logRequest(ctx, url, req.Header, body)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error making http request", err.Error())
		return c
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reading response body", err.Error())
		return c
	}

	logResponse(ctx, resp, respBody)

	if resp.StatusCode != http.StatusOK {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reporting stat", resp.Status, string(respBody))
	}

	return c
}

func setupLoggerCtx(requestID string) context.Context {
	// empty context for now
	ctx := context.Background()

	// setup logger with rid attached
	logger := log.Logger.With().Ctx(ctx).Str("rid", requestID).Logger()

	// return context for logging
	return logger.WithContext(ctx)
}

func logRequest(ctx context.Context, url string, headers http.Header, body []byte) {
	logging.LogInfoCtx(ctx, "sending request to: "+url)
	logging.LogDebugCtx(ctx, fmt.Sprintf("request: headers=%s; body=%s", utils.HeadersToStr(headers), string(body)))
}

func logResponse(ctx context.Context, resp *http.Response, respBody []byte) {
	logging.LogDebugCtx(ctx, fmt.Sprintf("response: %v; headers=%s; body=%s", resp.Status, utils.HeadersToStr(resp.Header), respBody))
}
