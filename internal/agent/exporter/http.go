package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ex0rcist/metflix/internal/compression"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/retrier"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/ex0rcist/metflix/pkg/metrics"
)

var _ Exporter = (*HTTPExporter)(nil)

// An exporter to send metrics one-by-one in parallel.
type HTTPExporter struct {
	baseURL   *entities.Address
	client    *http.Client
	signer    security.Signer
	publicKey security.PublicKey
	context   context.Context

	buffer []metrics.MetricExchange
	jobs   chan metrics.MetricExchange
	err    error
}

// Constructor.
func NewHTTPExporter(
	ctx context.Context,
	baseURL *entities.Address,
	signer security.Signer,
	numWorkers int,
	publicKey security.PublicKey,
) *HTTPExporter {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	exporter := &HTTPExporter{
		baseURL:   baseURL,
		context:   ctx,
		client:    client,
		signer:    signer,
		publicKey: publicKey,
		jobs:      make(chan metrics.MetricExchange, 30),
	}

	exporter.spawnWorkers(numWorkers)

	return exporter
}

// Add metric to buffer.
func (e *HTTPExporter) Add(name string, value metrics.Metric) Exporter {
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
		logging.LogError(entities.ErrMetricReport, "unknown metric")

		e.err = entities.ErrMetricUnknown
		return e
	}

	e.buffer = append(e.buffer, mex)

	return e
}

// Send buffer out.
func (e *HTTPExporter) Send() error {
	if e.err != nil {
		return e.err
	}

	if len(e.buffer) == 0 {
		return fmt.Errorf("cannot send empty buffer")
	}

	logging.LogDebugF("sending %d jobs to channel", len(e.buffer))

	for _, mex := range e.buffer {
		e.jobs <- mex
	}

	e.Reset()

	return nil
}

// Reset buffer.
func (e *HTTPExporter) Reset() {
	e.buffer = make([]metrics.MetricExchange, 0)
	e.err = nil
}

// Get error if any.
func (e *HTTPExporter) Error() error {
	if e.err == nil {
		return nil
	}

	return fmt.Errorf("metrics export failed: %w", e.err)
}

func (e *HTTPExporter) spawnWorkers(numWorkers int) {
	for w := 1; w <= numWorkers; w++ {
		go e.worker(w)
	}
}

func (e *HTTPExporter) worker(id int) {
	delays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for mex := range e.jobs {
		logging.LogDebugF("worker #%d started job", id)

		err := retrier.New(
			func() error { return e.doSend(mex) },
			func(err error) bool {
				_, ok := err.(entities.RetriableError)
				return ok
			},
			retrier.WithDelays(delays),
		).Run(e.context)

		if err != nil {
			logging.LogError(err, "error during async working")
		}

		logging.LogDebugF("worker #%d ended job", id)
	}
}

func (e *HTTPExporter) doSend(mex metrics.MetricExchange) error {
	requestID := utils.GenerateRequestID()
	ctx := setupLoggerCtx(requestID)

	body, err := json.Marshal(mex)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during marshaling", err.Error())
		return err
	}

	payload, err := compression.Pack(body)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during compression", err.Error())
		return err
	}

	if e.publicKey != nil {
		payload, err = security.Encrypt(io.Reader(payload), e.publicKey)
		if err != nil {
			return err
		}
	}

	url := "http://" + e.baseURL.String() + "/update"

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "httpRequest error", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("X-Request-Id", requestID)

	clientIP, err := utils.GetOutboundIP()
	if err != nil {
		return err
	}
	req.Header.Set("X-Real-IP", clientIP.String())

	if e.signer != nil {
		signature, signErr := e.signer.CalculateSignature(payload.Bytes())
		if signErr != nil {
			logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error during signing", signErr.Error())
			return signErr
		}

		req.Header.Set("HashSHA256", signature)
	}

	logHTTPRequest(ctx, url, req.Header, body)

	resp, err := e.client.Do(req)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error making http request", err.Error())
		return entities.RetriableError{Err: err, RetryAfter: 10 * time.Second}
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logging.LogError(closeErr)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reading response body", err.Error())
		return err
	}

	logHTTPResponse(ctx, resp, respBody)

	if resp.StatusCode != http.StatusOK {
		formatedBody := strings.ReplaceAll(string(respBody), "\n", "")
		logging.LogErrorCtx(ctx, entities.ErrMetricReport, "error reporting stats", resp.Status, formatedBody)
		return err
	}

	return nil
}
