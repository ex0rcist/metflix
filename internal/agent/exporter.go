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

type Exporter interface {
	Add(name string, value metrics.Metric) Exporter
	Send() Exporter
	Error() error
	Reset()
}

type MetricsExporter struct {
	baseURL *entities.Address
	client  *http.Client

	buffer []metrics.MetricExchange
	err    error
}

func NewMetricsExporter(baseURL *entities.Address, httpTransport http.RoundTripper) *MetricsExporter {
	if httpTransport == nil {
		httpTransport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: httpTransport,
	}

	return &MetricsExporter{
		baseURL: baseURL,
		client:  client,
	}
}

func (e *MetricsExporter) Add(name string, value metrics.Metric) Exporter {
	if e.err != nil {
		return e
	}

	var mex metrics.MetricExchange
	switch value.Kind() {
	case metrics.KindCounter:
		mex = metrics.NewUpdateCounterMex(name, value.(metrics.Counter))

	case metrics.KindGauge:
		mex = metrics.NewUpdateGaugeMex(name, value.(metrics.Gauge))

	default:
		e.err = entities.ErrMetricUnknown
		return e
	}

	e.buffer = append(e.buffer, mex)

	return e
}

// Send metrics stored in internal buffer to metrics collector in single batch request.
func (e *MetricsExporter) Send() Exporter {
	if e.err != nil {
		return e
	}

	if len(e.buffer) == 0 {
		e.err = fmt.Errorf("cannot send empty buffer")
		return e
	}

	e.err = e.doSend()

	return e
}

func (e *MetricsExporter) doSend() error {
	requestID := utils.GenerateRequestID()
	ctx := setupLoggerCtx(requestID)

	body, err := json.Marshal(e.buffer)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during marshaling", err.Error())
		return err
	}

	payload, err := compression.Pack(body)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during compression", err.Error())
		return err
	}

	url := "http://" + e.baseURL.String() + "/updates"
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "httpRequest error", err.Error())
		return err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("X-Request-Id", requestID)

	logRequest(ctx, url, req.Header, body)

	resp, err := e.client.Do(req)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error making http request", err.Error())
		return err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reading response body", err.Error())
		return err
	}

	logResponse(ctx, resp, respBody)

	if resp.StatusCode != http.StatusOK {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reporting stats", resp.Status, string(respBody))
		return err
	}

	return nil
}

func (e *MetricsExporter) Reset() {
	e.buffer = make([]metrics.MetricExchange, 0)
	e.err = nil
}

func (e *MetricsExporter) Error() error {
	if e.err == nil {
		return nil
	}

	return fmt.Errorf("metrics export failed: %w", e.err)
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
