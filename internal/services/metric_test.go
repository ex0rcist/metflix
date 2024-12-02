package services

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/ex0rcist/metflix/pkg/metrics"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Get(t *testing.T) {
	type args struct {
		name string
		kind string
	}

	tests := []struct {
		name     string
		mock     func(m *storage.StorageMock)
		args     args
		expected storage.Record
		wantErr  bool
	}{
		{
			name: "existing record",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			args:     args{name: "test", kind: metrics.KindCounter},
			expected: storage.Record{Name: "test", Value: metrics.Counter(42)},
			wantErr:  false,
		},

		{
			name: "non-existing storage.Record",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(storage.Record{}, entities.ErrRecordNotFound)
			},
			args:     args{name: "test", kind: metrics.KindCounter},
			expected: storage.Record{},
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(storage.StorageMock)
			service := NewMetricService(m)

			if tt.mock != nil {
				tt.mock(m)
			}

			result, err := service.Get(ctx, tt.args.name, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got %v", tt.wantErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestService_Push(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(m *storage.StorageMock)
		record   storage.Record
		expected storage.Record
		wantErr  bool
	}{
		{
			name: "new counter storage.Record",
			mock: func(m *storage.StorageMock) {
				r := storage.Record{Name: "test", Value: metrics.Counter(42)}

				m.On("Get", mock.Anything, "test_counter").Return(storage.Record{}, entities.ErrRecordNotFound)
				m.On("Push", mock.Anything, "test_counter", r).Return(nil) // no error, successful push
			},
			record:   storage.Record{Name: "test", Value: metrics.Counter(42)},
			expected: storage.Record{Name: "test", Value: metrics.Counter(42)},
			wantErr:  false,
		},
		{
			name: "update counter storage.Record",
			mock: func(m *storage.StorageMock) {
				oldr := storage.Record{Name: "test", Value: metrics.Counter(42)}
				newr := storage.Record{Name: "test", Value: metrics.Counter(84)}

				m.On("Get", mock.Anything, "test_counter").Return(oldr, nil)
				m.On("Push", mock.Anything, "test_counter", newr).Return(nil) // no error, successful push
			},
			record:   storage.Record{Name: "test", Value: metrics.Counter(42)},
			expected: storage.Record{Name: "test", Value: metrics.Counter(84)},
			wantErr:  false,
		},
		{
			name: "new gauge record",
			mock: func(m *storage.StorageMock) {
				r := storage.Record{Name: "test", Value: metrics.Gauge(42.42)}

				m.On("Get", mock.Anything, "test_gauge").Return(storage.Record{}, entities.ErrRecordNotFound)
				m.On("Push", mock.Anything, "test_gauge", r).Return(nil) // no error, successful push
			},
			record:   storage.Record{Name: "test", Value: metrics.Gauge(42.42)},
			expected: storage.Record{Name: "test", Value: metrics.Gauge(42.42)},
			wantErr:  false,
		},
		{
			name: "update gauge record",
			mock: func(m *storage.StorageMock) {
				oldr := storage.Record{Name: "test", Value: metrics.Gauge(42.42)}
				newr := storage.Record{Name: "test", Value: metrics.Gauge(43.43)}

				m.On("Get", mock.Anything, "test_gauge").Return(oldr, nil)
				m.On("Push", mock.Anything, "test_gauge", newr).Return(nil) // no error, successful push
			},
			record:   storage.Record{Name: "test", Value: metrics.Gauge(43.43)},
			expected: storage.Record{Name: "test", Value: metrics.Gauge(43.43)},
			wantErr:  false,
		},
		{
			name: "underlying error",
			mock: func(m *storage.StorageMock) {
				m.On("Push", mock.Anything, "test_gauge", mock.AnythingOfType("storage.Record")).Return(entities.ErrUnexpected)
			},
			record:   storage.Record{Name: "test", Value: metrics.Gauge(43.43)},
			expected: storage.Record{},
			wantErr:  true,
		},
		{
			name:     "missing record name",
			record:   storage.Record{Name: "", Value: metrics.Counter(43)},
			expected: storage.Record{},
			wantErr:  true,
		},
		{
			name: "storage get error",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(storage.Record{}, entities.ErrUnexpected)
			},
			record:   storage.Record{Name: "test", Value: metrics.Counter(43)},
			expected: storage.Record{},
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(storage.StorageMock)
			service := NewMetricService(m)

			if tt.mock != nil {
				tt.mock(m)
			}

			result, err := service.Push(ctx, tt.record)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got %v", tt.wantErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestService_PushList(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(m *storage.StorageMock)
		records  []storage.Record
		expected []storage.Record
		wantErr  bool
	}{
		{
			name: "should push list",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "existedCounter_counter").Return(storage.Record{Name: "existedCounter", Value: metrics.Counter(42)}, nil)
				m.On("Get", mock.Anything, "newCounter_counter").Return(storage.Record{}, entities.ErrRecordNotFound)

				m.On("PushList", mock.Anything, mock.AnythingOfType("map[string]storage.Record")).Return(nil) // no error, successful push
			},
			records: []storage.Record{
				{Name: "existedCounter", Value: metrics.Counter(42)},
				{Name: "newCounter", Value: metrics.Counter(42)},
				{Name: "newGauge", Value: metrics.Gauge(42.42)},
			},
			expected: []storage.Record{
				{Name: "existedCounter", Value: metrics.Counter(84)},
				{Name: "newCounter", Value: metrics.Counter(42)},
				{Name: "newGauge", Value: metrics.Gauge(42.42)},
			},
			wantErr: false,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(storage.StorageMock)
			service := NewMetricService(m)

			if tt.mock != nil {
				tt.mock(m)
			}

			result, err := service.PushList(ctx, tt.records)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got %v", tt.wantErr, err)
			}

			require.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(m *storage.StorageMock)
		expected []storage.Record
		wantErr  bool
	}{
		{
			name: "normal list",
			mock: func(m *storage.StorageMock) {
				m.On("List", mock.Anything).Return([]storage.Record{
					{Name: "metricX", Value: metrics.Counter(42)},
					{Name: "metricA", Value: metrics.Gauge(42.42)},
				}, nil)
			},
			expected: []storage.Record{ // sorted
				{Name: "metricA", Value: metrics.Gauge(42.42)},
				{Name: "metricX", Value: metrics.Counter(42)},
			},
			wantErr: false,
		},

		{
			name: "had error",
			mock: func(m *storage.StorageMock) {
				m.On("List", mock.Anything).Return([]storage.Record{}, entities.ErrUnexpected)
			},
			expected: []storage.Record{},
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(storage.StorageMock)
			service := NewMetricService(m)

			if tt.mock != nil {
				tt.mock(m)
			}

			result, err := service.List(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got %v", tt.wantErr, err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d records, got %d", len(tt.expected), len(result))
			}

			for i, record := range result {
				if record != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected[i], record)
				}
			}

		})
	}

}

func TestService_calculateNewValue(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(m *storage.StorageMock)
		record   storage.Record
		expected metrics.Metric
		wantErr  bool
	}{
		{
			name: "new counter record",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(storage.Record{}, entities.ErrRecordNotFound)
			},
			record:   storage.Record{Name: "test", Value: metrics.Counter(42)},
			expected: metrics.Counter(42),
			wantErr:  false,
		},
		{
			name: "existing counter record",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(storage.Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			record:   storage.Record{Name: "test", Value: metrics.Counter(42)},
			expected: metrics.Counter(84),
			wantErr:  false,
		},
		{
			name: "new gauge record",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_gauge").Return(storage.Record{}, entities.ErrRecordNotFound)
			},
			record:   storage.Record{Name: "test", Value: metrics.Gauge(42.42)},
			expected: metrics.Gauge(42.42),
			wantErr:  false,
		},
		{
			name: "existing gauge record",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_gauge").Return(storage.Record{Name: "test", Value: metrics.Gauge(42.42)})
			},
			record:   storage.Record{Name: "test", Value: metrics.Gauge(43.43)},
			expected: metrics.Gauge(43.43),
			wantErr:  false,
		},
		{
			name: "underlying error",
			mock: func(m *storage.StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(storage.Record{}, entities.ErrUnexpected)
			},
			record:   storage.Record{Name: "test", Value: metrics.Counter(42)},
			expected: nil,
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := new(storage.StorageMock)
			service := NewMetricService(m)

			if tt.mock != nil {
				tt.mock(m)
			}

			result, err := service.calculateNewValue(ctx, tt.record)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got %v", tt.wantErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
