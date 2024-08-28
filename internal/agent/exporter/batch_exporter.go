package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ex0rcist/metflix/internal/compression"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/utils"
)

var _ Exporter = (*BatchExporter)(nil)

type BatchExporter struct {
	baseURL *entities.Address
	client  *http.Client
	signer  services.Signer

	buffer []metrics.MetricExchange
	err    error
}

func NewBatchExporter(baseURL *entities.Address, signer services.Signer) *BatchExporter {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	return &BatchExporter{
		baseURL: baseURL,
		client:  client,
		signer:  signer,
	}
}

func (e *BatchExporter) Add(name string, value metrics.Metric) Exporter {
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

func (e *BatchExporter) Send() error {
	if e.err != nil {
		return e.err
	}

	if len(e.buffer) == 0 {
		return fmt.Errorf("cannot send empty buffer")
	}

	err := utils.NewRetrier(
		func() error { return e.doSend() },
		func(err error) bool {
			_, ok := err.(entities.RetriableError)
			return ok
		},
		[]time.Duration{
			1 * time.Second,
			3 * time.Second,
			5 * time.Second,
		},
	).Run()

	e.Reset()

	return err
}

func (e *BatchExporter) Reset() {
	e.buffer = make([]metrics.MetricExchange, 0)
	e.err = nil
}

func (e *BatchExporter) Error() error {
	if e.err == nil {
		return nil
	}

	return fmt.Errorf("metrics export failed: %w", e.err)
}

func (e *BatchExporter) doSend() error {
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("X-Request-Id", requestID)

	if e.signer != nil {
		signature, err := e.signer.CalculateSignature(payload.Bytes())
		if err != nil {
			logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during signing", err.Error())
			return err
		}

		req.Header.Set("HashSHA256", signature)
	}

	logRequest(ctx, url, req.Header, body)

	resp, err := e.client.Do(req)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error making http request", err.Error())
		return entities.RetriableError{Err: err, RetryAfter: 10 * time.Second}
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reading response body", err.Error())
		return err
	}

	logResponse(ctx, resp, respBody)

	if resp.StatusCode != http.StatusOK {
		formatedBody := strings.ReplaceAll(string(respBody), "\n", "")
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reporting stats", resp.Status, formatedBody)
		return err
	}

	return nil
}
