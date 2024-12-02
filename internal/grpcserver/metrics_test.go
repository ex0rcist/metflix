package grpcserver

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	// "github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/ex0rcist/metflix/pkg/grpcapi"
	"github.com/ex0rcist/metflix/pkg/metrics"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestBatchUpdate(t *testing.T) {
	root, _ := utils.GetProjectRoot()
	prvKey, _ := security.NewPrivateKey(entities.FilePath(filepath.Join(root, "example_key.pem")))

	batchReq := []*grpcapi.MetricExchange{
		grpcapi.NewUpdateCounterMex("PollCount", 10),
		grpcapi.NewUpdateGaugeMex("Alloc", 11.23),
	}

	batchResp := []storage.Record{
		{Name: "Alloc", Value: metrics.Gauge(11.23)},
		{Name: "PollCount", Value: metrics.Counter(10)},
	}

	type expected struct {
		code     codes.Code
		response []*grpcapi.MetricExchange
	}

	tt := []struct {
		name       string
		data       []*grpcapi.MetricExchange
		prvKey     security.PrivateKey
		serviceRsp []storage.Record
		serviceErr error
		expected   expected
	}{
		{
			name:       "Successful batch update",
			data:       batchReq,
			serviceRsp: batchResp,
			expected: expected{
				code: codes.OK,
				response: []*grpcapi.MetricExchange{
					grpcapi.NewUpdateGaugeMex("Alloc", 11.23),
					grpcapi.NewUpdateCounterMex("PollCount", 10),
				},
			},
		},
		{
			name:       "Batch update on RSA-protected server should fail",
			data:       batchReq,
			prvKey:     prvKey,
			serviceRsp: []storage.Record{},
			expected:   expected{code: codes.InvalidArgument},
		},
		{
			name: "Batch update fails on empty list",
			data: make([]*grpcapi.MetricExchange, 0),
			expected: expected{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Batch update fails if unknown metric kind found in list",
			data: []*grpcapi.MetricExchange{
				{Id: "xxx", Mtype: "unknown"},
			},
			expected: expected{
				code: codes.InvalidArgument,
			},
		},
		{
			name:       "Batch update fails if service is broken",
			data:       batchReq,
			serviceErr: entities.ErrUnexpected,
			expected: expected{
				code: codes.Internal,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			m := new(services.MetricServiceMock)
			m.On("PushList", mock.Anything, mock.Anything).Return(tc.serviceRsp, tc.serviceErr)

			conn, closer := createTestServer(t, m, nil, tc.prvKey)
			t.Cleanup(closer)

			client := grpcapi.NewMetricsClient(conn)
			req := &grpcapi.BatchUpdateRequest{Data: tc.data}
			resp, err := client.BatchUpdate(context.Background(), req)

			rv, ok := status.FromError(err)

			require.True(t, ok)
			require.Equal(t, tc.expected.code, rv.Code())

			if tc.expected.code == codes.OK {
				require.Equal(t, tc.expected.response, resp.Data)
			}
		})
	}
}

func TestBatchUpdateEncrypted(t *testing.T) {
	root, _ := utils.GetProjectRoot()
	prvKey, _ := security.NewPrivateKey(entities.FilePath(filepath.Join(root, "example_key.pem")))
	pubKey, _ := security.NewPublicKey(entities.FilePath(filepath.Join(root, "example_key.pub.pem")))

	batchReq := []*grpcapi.MetricExchange{
		grpcapi.NewUpdateCounterMex("PollCount", 10),
		grpcapi.NewUpdateGaugeMex("Alloc", 11.23),
	}

	batchResp := []storage.Record{
		{Name: "Alloc", Value: metrics.Gauge(11.23)},
		{Name: "PollCount", Value: metrics.Counter(10)},
	}

	type expected struct {
		code     codes.Code
		response []*grpcapi.MetricExchange
	}

	tt := []struct {
		name       string
		data       []*grpcapi.MetricExchange
		prvKey     security.PrivateKey
		pubKey     security.PublicKey
		serviceRsp []storage.Record
		serviceErr error
		expected   expected
	}{
		{
			name:       "Successful batch update",
			data:       batchReq,
			prvKey:     prvKey,
			pubKey:     pubKey,
			serviceRsp: batchResp,
			expected: expected{
				code: codes.OK,
				response: []*grpcapi.MetricExchange{
					grpcapi.NewUpdateGaugeMex("Alloc", 11.23),
					grpcapi.NewUpdateCounterMex("PollCount", 10),
				},
			},
		},
		{
			name:   "Batch update fails on empty list",
			prvKey: prvKey,
			pubKey: pubKey,
			data:   make([]*grpcapi.MetricExchange, 0),
			expected: expected{
				code: codes.InvalidArgument,
			},
		},
		{
			name:   "Batch update fails if unknown metric kind found in list",
			prvKey: prvKey,
			pubKey: pubKey,
			data: []*grpcapi.MetricExchange{
				{Id: "xxx", Mtype: "unknown"},
			},
			expected: expected{
				code: codes.InvalidArgument,
			},
		},
		{
			name:       "Batch update fails if service is broken",
			prvKey:     prvKey,
			pubKey:     pubKey,
			data:       batchReq,
			serviceErr: entities.ErrUnexpected,
			expected: expected{
				code: codes.Internal,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			m := new(services.MetricServiceMock)
			m.On("PushList", mock.Anything, mock.Anything).Return(tc.serviceRsp, tc.serviceErr)

			conn, closer := createTestServer(t, m, nil, tc.prvKey)
			t.Cleanup(closer)

			client := grpcapi.NewMetricsClient(conn)

			encrypted, err := encrypt(tc.data, tc.pubKey)
			require.NoError(t, err)

			req := &grpcapi.BatchUpdateEncryptedRequest{EncryptedData: encrypted}
			resp, err := client.BatchUpdateEncrypted(context.Background(), req)

			rv, ok := status.FromError(err)

			require.True(t, ok)
			require.Equal(t, tc.expected.code, rv.Code())

			if tc.expected.code == codes.OK {
				require.Equal(t, tc.expected.response, resp.Data)
			}
		})
	}
}

func encrypt(data []*grpcapi.MetricExchange, key security.PublicKey) ([]byte, error) {
	metricsBytes, err := proto.Marshal(&grpcapi.BatchUpdateRequest{Data: data})
	if err != nil {
		return nil, err
	}

	payload, err := security.Encrypt(bytes.NewReader(metricsBytes), key)
	if err != nil {
		return nil, err
	}

	return payload.Bytes(), nil
}
