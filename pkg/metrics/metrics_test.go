package metrics

import (
	"strconv"
	"testing"
)

// Test ToCounter function
func TestToCounter(t *testing.T) {
	tests := []struct {
		input       string
		expected    Counter
		expectError bool
	}{
		{"42", Counter(42), false},
		{"-10", Counter(-10), false},
		{"not_a_number", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ToCounter(tt.input)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

// Test ToGauge function
func TestToGauge(t *testing.T) {
	tests := []struct {
		input       string
		expected    Gauge
		expectError bool
	}{
		{"42.5", Gauge(42.5), false},
		{"-10.3", Gauge(-10.3), false},
		{"not_a_number", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ToGauge(tt.input)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

// Test Metric methods for Counter
func TestCounterMethods(t *testing.T) {
	var c Counter = 100
	if c.Kind() != KindCounter {
		t.Errorf("expected: %v, got: %v", KindCounter, c.Kind())
	}
	if c.String() != strconv.FormatInt(int64(c), 10) {
		t.Errorf("expected: %v, got: %v", strconv.FormatInt(int64(c), 10), c.String())
	}
}

// Test Metric methods for Gauge
func TestGaugeMethods(t *testing.T) {
	var g Gauge = 100.25
	if g.Kind() != KindGauge {
		t.Errorf("expected: %v, got: %v", KindGauge, g.Kind())
	}
	if g.String() != strconv.FormatFloat(float64(g), 'f', -1, 64) {
		t.Errorf("expected: %v, got: %v", strconv.FormatFloat(float64(g), 'f', -1, 64), g.String())
	}
}

// Test NewUpdateCounterMex function
func TestNewUpdateCounterMex(t *testing.T) {
	var value Counter = 123
	mex := NewUpdateCounterMex("test_counter", value)
	if mex.ID != "test_counter" {
		t.Errorf("expected ID: %v, got: %v", "test_counter", mex.ID)
	}
	if mex.MType != KindCounter {
		t.Errorf("expected MType: %v, got: %v", KindCounter, mex.MType)
	}
	if mex.Delta == nil || *mex.Delta != value {
		t.Errorf("expected Delta: %v, got: %v", value, mex.Delta)
	}
	if mex.Value != nil {
		t.Errorf("expected Value to be nil, got: %v", mex.Value)
	}
}

// Test NewUpdateGaugeMex function
func TestNewUpdateGaugeMex(t *testing.T) {
	var value Gauge = 123.45
	mex := NewUpdateGaugeMex("test_gauge", value)
	if mex.ID != "test_gauge" {
		t.Errorf("expected ID: %v, got: %v", "test_gauge", mex.ID)
	}
	if mex.MType != KindGauge {
		t.Errorf("expected MType: %v, got: %v", KindGauge, mex.MType)
	}
	if mex.Value == nil || *mex.Value != value {
		t.Errorf("expected Value: %v, got: %v", value, mex.Value)
	}
	if mex.Delta != nil {
		t.Errorf("expected Delta to be nil, got: %v", mex.Delta)
	}
}

// Test NewGetCounterMex function
func TestNewGetCounterMex(t *testing.T) {
	mex := NewGetCounterMex("get_counter")
	if mex.ID != "get_counter" {
		t.Errorf("expected ID: %v, got: %v", "get_counter", mex.ID)
	}
	if mex.MType != KindCounter {
		t.Errorf("expected MType: %v, got: %v", KindCounter, mex.MType)
	}
	if mex.Delta != nil {
		t.Errorf("expected Delta to be nil, got: %v", mex.Delta)
	}
	if mex.Value != nil {
		t.Errorf("expected Value to be nil, got: %v", mex.Value)
	}
}

// Test NewGetGaugeMex function
func TestNewGetGaugeMex(t *testing.T) {
	mex := NewGetGaugeMex("get_gauge")
	if mex.ID != "get_gauge" {
		t.Errorf("expected ID: %v, got: %v", "get_gauge", mex.ID)
	}
	if mex.MType != KindGauge {
		t.Errorf("expected MType: %v, got: %v", KindGauge, mex.MType)
	}
	if mex.Delta != nil {
		t.Errorf("expected Delta to be nil, got: %v", mex.Delta)
	}
	if mex.Value != nil {
		t.Errorf("expected Value to be nil, got: %v", mex.Value)
	}
}
