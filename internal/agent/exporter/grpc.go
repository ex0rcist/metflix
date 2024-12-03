package exporter

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/ex0rcist/metflix/pkg/grpcapi"
	"github.com/ex0rcist/metflix/pkg/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// GRPCExporter sends collected metrics to metrics collector in single batch request.
type GRPCExporter struct {
	baseURL   *entities.Address
	publicKey security.PublicKey

	conn   *grpc.ClientConn
	buffer []*grpcapi.MetricExchange
	err    error
}

// Construct new GRPCEXporter.
func NewGRPCExporter(baseURL *entities.Address, publicKey security.PublicKey) *GRPCExporter {
	return &GRPCExporter{baseURL: baseURL, publicKey: publicKey}
}

// Add a metric to internal buffer.
func (e *GRPCExporter) Add(name string, value metrics.Metric) Exporter {
	if e.err != nil {
		return e
	}

	var req *grpcapi.MetricExchange
	switch v := value.(type) {
	case metrics.Counter:
		req = grpcapi.NewUpdateCounterMex(name, v)
	case metrics.Gauge:
		req = grpcapi.NewUpdateGaugeMex(name, v)
	default:
		e.err = entities.ErrMetricUnknown
		return e
	}

	e.buffer = append(e.buffer, req)

	return e
}

func (e *GRPCExporter) Error() error {
	if e.err == nil {
		return nil
	}

	return fmt.Errorf("metrics export failed: %w", e.err)
}

// Send metrics stored in internal buffer to metrics collector in single batch request.
func (e *GRPCExporter) Send() error {
	if e.err != nil {
		return e.err
	}

	if len(e.buffer) == 0 {
		return fmt.Errorf("cannot send empty buffer")
	}

	if e.conn == nil {
		e.conn, e.err = grpc.NewClient(e.baseURL.String(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		)

		if e.err != nil {
			return e.err
		}
	}

	md, err := e.prepareMetadata()
	if err != nil {
		logging.LogError(err)
	}

	ctx := setupLoggerCtx(md.Get("x-request-id")[0])

	ctx = metadata.NewOutgoingContext(ctx, md)
	client := grpcapi.NewMetricsClient(e.conn)

	if e.publicKey != nil {
		metricsBytes, err := proto.Marshal(&grpcapi.BatchUpdateRequest{Data: e.buffer})
		if err != nil {
			return err
		}

		payload, err := security.Encrypt(bytes.NewReader(metricsBytes), e.publicKey)
		if err != nil {
			return err
		}

		logging.LogDebugCtx(ctx, fmt.Sprintf("sending gRPC %s to %s...", "BatchUpdateEncryptedRequest", e.baseURL.String()))

		req := &grpcapi.BatchUpdateEncryptedRequest{EncryptedData: payload.Bytes()}
		_, e.err = client.BatchUpdateEncrypted(ctx, req)

		logResponseFromErr(ctx, e.err)
	} else {
		logging.LogDebugCtx(ctx, fmt.Sprintf("sending gRPC %s to %s...", "BatchUpdateRequest", e.baseURL.String()))

		req := &grpcapi.BatchUpdateRequest{Data: e.buffer}
		_, e.err = client.BatchUpdate(ctx, req)

		logResponseFromErr(ctx, e.err)
	}

	e.Reset()

	return nil
}

func logResponseFromErr(ctx context.Context, err error) {
	st, _ := status.FromError(err)
	logging.LogDebugCtx(ctx, fmt.Sprintf("got response status=%s", st.Code()))
}

func (e *GRPCExporter) Reset() {
	e.buffer = make([]*grpcapi.MetricExchange, 0)
	e.err = nil
}

func (e *GRPCExporter) Close() error {
	if e.conn == nil {
		return nil
	}

	return e.conn.Close()
}

func (e *GRPCExporter) prepareMetadata() (metadata.MD, error) {
	md := metadata.New(map[string]string{})

	clientIP, err := utils.GetOutboundIP()
	if err != nil {
		return md, err
	}

	md.Set("x-real-ip", clientIP.String())
	md.Set("x-request-id", utils.GenerateRequestID())

	return md, nil
}
