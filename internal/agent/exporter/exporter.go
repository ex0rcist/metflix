package exporter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/ex0rcist/metflix/pkg/metrics"

	"github.com/rs/zerolog/log"
)

// Exporter interface.
type Exporter interface {
	Add(name string, value metrics.Metric) Exporter
	Send() error
	Error() error
	Reset()
}

// Create new instance of Exporter for specified transport.
func New(
	ctx context.Context,
	transport string,
	address *entities.Address,
	rateLimit int,
	signer security.Signer,
	publicKey security.PublicKey,
) (Exporter, error) {
	var exp Exporter

	switch transport {
	case entities.TransportHTTP:
		if rateLimit > 0 {
			exp = NewHTTPExporter(ctx, address, signer, rateLimit, publicKey)
		} else {
			exp = NewHTTPBatchExporter(ctx, address, signer, publicKey)
		}
	case entities.TransportGRPC:
		exp = NewGRPCExporter(address, publicKey)
	default:
		return exp, entities.ErrUnknownTransport(transport)
	}

	return exp, nil
}

func setupLoggerCtx(requestID string) context.Context {
	// empty context for now
	ctx := context.Background()

	// setup logger with rid attached
	logger := log.Logger.With().Ctx(ctx).Str("rid", requestID).Logger()

	// return context for logging
	return logger.WithContext(ctx)
}

func logHTTPRequest(ctx context.Context, url string, headers http.Header, body []byte) {
	logging.LogInfoCtx(ctx, "sending request to: "+url)
	logging.LogDebugCtx(ctx, fmt.Sprintf("request: headers=%s; body=%s", utils.HeadersToStr(headers), string(body)))
}

func logHTTPResponse(ctx context.Context, resp *http.Response, respBody []byte) {
	logging.LogDebugCtx(ctx, fmt.Sprintf("response: %v; headers=%s; body=%s", resp.Status, utils.HeadersToStr(resp.Header), respBody))
}
