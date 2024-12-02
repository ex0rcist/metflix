package grpcserver

import (
	"bytes"
	"context"

	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/pkg/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// MetricsServer allows to store and retrieve metrics.
type MetricsServer struct {
	grpcapi.UnimplementedMetricsServer

	privateKey    security.PrivateKey
	metricService services.MetricProvider
}

// RegisterMetricsServer creates new instance of gRPC serving Metrics API and attaches it to the server.
func RegisterMetricsServer(server *grpc.Server, metricService services.MetricProvider, privateKey security.PrivateKey) {
	s := &MetricsServer{metricService: metricService, privateKey: privateKey}

	grpcapi.RegisterMetricsServer(server, s)
}

// BatchUpdate pushes list of metrics data.
func (s MetricsServer) BatchUpdate(ctx context.Context, req *grpcapi.BatchUpdateRequest) (*grpcapi.BatchUpdateResponse, error) {
	// do not allow if server is configured with RSA encoding
	if s.privateKey != nil {
		return nil, status.Errorf(codes.InvalidArgument, "please use encrypted endpoint")
	}

	return s.batchUpdate(ctx, req)
}

// BatchUpdateEncrypted decodes encrypted data and pushes list of metrics data.
func (s MetricsServer) BatchUpdateEncrypted(ctx context.Context, encReq *grpcapi.BatchUpdateEncryptedRequest) (*grpcapi.BatchUpdateResponse, error) {
	buff, err := security.Decrypt(bytes.NewReader(encReq.EncryptedData), s.privateKey)
	if err != nil {
		return nil, err
	}

	req := &grpcapi.BatchUpdateRequest{}
	err = proto.Unmarshal(buff.Bytes(), req)
	if err != nil {
		return nil, err
	}

	return s.batchUpdate(ctx, req)
}

func (s MetricsServer) batchUpdate(ctx context.Context, req *grpcapi.BatchUpdateRequest) (*grpcapi.BatchUpdateResponse, error) {
	records, err := toRecordsList(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	records, err = s.metricService.PushList(ctx, records)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data, err := toMetricExchangeList(records)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcapi.BatchUpdateResponse{Data: data}, nil
}
