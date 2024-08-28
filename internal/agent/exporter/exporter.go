package exporter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/utils"

	"github.com/rs/zerolog/log"
)

const (
	KindBatch   = "batch"
	KindLimited = "limited"
)

type Exporter interface {
	Add(name string, value metrics.Metric) Exporter
	Send() error
	Error() error
	Reset()
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
