package storage

import "github.com/ex0rcist/metflix/internal/metrics"

type Record struct {
	Name  string
	Value metrics.Metric
}

func CalculateRecordID(name, kind string) string {
	if len(name) == 0 || len(kind) == 0 {
		return ""
	}

	return name + "_" + kind
}

func (r Record) CalculateRecordID() string {
	return CalculateRecordID(r.Name, r.Value.Kind())
}
