package storage

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
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
		mock     func(m *StorageMock)
		args     args
		expected Record
		wantErr  bool
	}{
		{
			name: "existing record",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(Record{Name: "test", Value: metrics.Counter(42)}, nil)

			},
			args:     args{name: "test", kind: metrics.KindCounter},
			expected: Record{Name: "test", Value: metrics.Counter(42)},
			wantErr:  false,
		},

		{
			name: "non-existing record",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(Record{}, entities.ErrRecordNotFound)
			},
			args:     args{name: "test", kind: metrics.KindCounter},
			expected: Record{},
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(StorageMock)
			service := NewService(m)

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
		mock     func(m *StorageMock)
		record   Record
		expected Record
		wantErr  bool
	}{
		{
			name: "new counter record",
			mock: func(m *StorageMock) {
				r := Record{Name: "test", Value: metrics.Counter(42)}

				m.On("Get", mock.Anything, "test_counter").Return(Record{}, entities.ErrRecordNotFound)
				m.On("Push", mock.Anything, "test_counter", r).Return(nil) // no error, successful push
			},
			record:   Record{Name: "test", Value: metrics.Counter(42)},
			expected: Record{Name: "test", Value: metrics.Counter(42)},
			wantErr:  false,
		},
		{
			name: "update counter record",
			mock: func(m *StorageMock) {
				oldr := Record{Name: "test", Value: metrics.Counter(42)}
				newr := Record{Name: "test", Value: metrics.Counter(84)}

				m.On("Get", mock.Anything, "test_counter").Return(oldr, nil)
				m.On("Push", mock.Anything, "test_counter", newr).Return(nil) // no error, successful push
			},
			record:   Record{Name: "test", Value: metrics.Counter(42)},
			expected: Record{Name: "test", Value: metrics.Counter(84)},
			wantErr:  false,
		},
		{
			name: "new gauge record",
			mock: func(m *StorageMock) {
				r := Record{Name: "test", Value: metrics.Gauge(42.42)}

				m.On("Get", mock.Anything, "test_gauge").Return(Record{}, entities.ErrRecordNotFound)
				m.On("Push", mock.Anything, "test_gauge", r).Return(nil) // no error, successful push
			},
			record:   Record{Name: "test", Value: metrics.Gauge(42.42)},
			expected: Record{Name: "test", Value: metrics.Gauge(42.42)},
			wantErr:  false,
		},
		{
			name: "update gauge record",
			mock: func(m *StorageMock) {
				oldr := Record{Name: "test", Value: metrics.Gauge(42.42)}
				newr := Record{Name: "test", Value: metrics.Gauge(43.43)}

				m.On("Get", mock.Anything, "test_gauge").Return(oldr, nil)
				m.On("Push", mock.Anything, "test_gauge", newr).Return(nil) // no error, successful push
			},
			record:   Record{Name: "test", Value: metrics.Gauge(43.43)},
			expected: Record{Name: "test", Value: metrics.Gauge(43.43)},
			wantErr:  false,
		},
		{
			name: "underlying error",
			mock: func(m *StorageMock) {
				m.On("Push", mock.Anything, "test_gauge", mock.AnythingOfType("Record")).Return(entities.ErrUnexpected)
			},
			record:   Record{Name: "test", Value: metrics.Gauge(43.43)},
			expected: Record{},
			wantErr:  true,
		},
		{
			name:     "missing record name",
			record:   Record{Name: "", Value: metrics.Counter(43)},
			expected: Record{},
			wantErr:  true,
		},
		{
			name: "storage get error",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(Record{}, entities.ErrUnexpected)
			},
			record:   Record{Name: "test", Value: metrics.Counter(43)},
			expected: Record{},
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(StorageMock)
			service := NewService(m)

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
		mock     func(m *StorageMock)
		records  []Record
		expected []Record
		wantErr  bool
	}{
		{
			name: "should push list",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "existedCounter_counter").Return(Record{Name: "existedCounter", Value: metrics.Counter(42)}, nil)
				m.On("Get", mock.Anything, "newCounter_counter").Return(Record{}, entities.ErrRecordNotFound)

				m.On("PushList", mock.Anything, mock.AnythingOfType("map[string]storage.Record")).Return(nil) // no error, successful push
			},
			records: []Record{
				{Name: "existedCounter", Value: metrics.Counter(42)},
				{Name: "newCounter", Value: metrics.Counter(42)},
				{Name: "newGauge", Value: metrics.Gauge(42.42)},
			},
			expected: []Record{
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
			m := new(StorageMock)
			service := NewService(m)

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
		mock     func(m *StorageMock)
		expected []Record
		wantErr  bool
	}{
		{
			name: "normal list",
			mock: func(m *StorageMock) {
				m.On("List", mock.Anything).Return([]Record{
					{Name: "metricX", Value: metrics.Counter(42)},
					{Name: "metricA", Value: metrics.Gauge(42.42)},
				}, nil)
			},
			expected: []Record{ // sorted
				{Name: "metricA", Value: metrics.Gauge(42.42)},
				{Name: "metricX", Value: metrics.Counter(42)},
			},
			wantErr: false,
		},

		{
			name: "had error",
			mock: func(m *StorageMock) {
				m.On("List", mock.Anything).Return([]Record{}, entities.ErrUnexpected)
			},
			expected: []Record{},
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(StorageMock)
			service := NewService(m)

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
		mock     func(m *StorageMock)
		record   Record
		expected metrics.Metric
		wantErr  bool
	}{
		{
			name: "new counter record",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(Record{}, entities.ErrRecordNotFound)
			},
			record:   Record{Name: "test", Value: metrics.Counter(42)},
			expected: metrics.Counter(42),
			wantErr:  false,
		},
		{
			name: "existing counter record",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(Record{Name: "test", Value: metrics.Counter(42)}, nil)
			},
			record:   Record{Name: "test", Value: metrics.Counter(42)},
			expected: metrics.Counter(84),
			wantErr:  false,
		},
		{
			name: "new gauge record",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_gauge").Return(Record{}, entities.ErrRecordNotFound)
			},
			record:   Record{Name: "test", Value: metrics.Gauge(42.42)},
			expected: metrics.Gauge(42.42),
			wantErr:  false,
		},
		{
			name: "existing gauge record",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_gauge").Return(Record{Name: "test", Value: metrics.Gauge(42.42)})
			},
			record:   Record{Name: "test", Value: metrics.Gauge(43.43)},
			expected: metrics.Gauge(43.43),
			wantErr:  false,
		},
		{
			name: "underlying error",
			mock: func(m *StorageMock) {
				m.On("Get", mock.Anything, "test_counter").Return(Record{}, entities.ErrUnexpected)
			},
			record:   Record{Name: "test", Value: metrics.Counter(42)},
			expected: nil,
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := new(StorageMock)
			service := NewService(m)

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
