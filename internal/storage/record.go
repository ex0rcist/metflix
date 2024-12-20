package storage

import (
	"encoding/json"
	"fmt"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/pkg/metrics"
)

// Struct to store metrics in storage
type Record struct {
	Name  string
	Value metrics.Metric
}

// Calculate record ID for ease of store and search
func CalculateRecordID(name, kind string) string {
	if len(name) == 0 || len(kind) == 0 {
		return ""
	}

	return name + "_" + kind
}

// Calculate record ID for ease of store and search
func (r Record) CalculateRecordID() string {
	return CalculateRecordID(r.Name, r.Value.Kind())
}

// Serialize to JSON
func (r Record) MarshalJSON() ([]byte, error) {
	jv, err := json.Marshal(map[string]string{
		"name":  r.Name,
		"kind":  r.Value.Kind(),
		"value": r.Value.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("record marshaling fail: %w", err)
	}

	return jv, nil
}

// Deserialize from JSON
func (r *Record) UnmarshalJSON(src []byte) error {
	var data map[string]string

	if err := json.Unmarshal(src, &data); err != nil {
		return fmt.Errorf("record unmarshaling failed: %w", err)
	}

	r.Name = data["name"]

	switch data["kind"] {
	case metrics.KindCounter:
		value, err := metrics.ToCounter(data["value"])
		if err != nil {
			return fmt.Errorf("record unmarshaling failed: %w", err)
		}

		r.Value = value
	case metrics.KindGauge:
		value, err := metrics.ToGauge(data["value"])
		if err != nil {
			return fmt.Errorf("record unmarshaling failed: %w", err)
		}

		r.Value = value
	default:
		return fmt.Errorf("record unmarshaling failed: %w", entities.ErrMetricUnknown)
	}

	return nil
}
