package storage

import "github.com/ex0rcist/metflix/internal/metrics"

type RecordID string
type Record struct {
	Name  string
	Value metrics.Metric
}

func CalculateRecordID(name, kind string) RecordID {
	if len(name) == 0 || len(kind) == 0 {
		return RecordID("")
	}

	return RecordID(name + "_" + kind)
}

func (r Record) CalculateRecordID() RecordID {
	return CalculateRecordID(r.Name, r.Value.Kind())
}
