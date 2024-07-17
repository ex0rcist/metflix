package entities

import "errors"

var (
	ErrBadAddressFormat = errors.New("bad net address format")

	ErrMetricNotFound     = errors.New("metric not found")
	ErrMetricUnknown      = errors.New("unknown metric type")
	ErrMetricReport       = errors.New("metric report error")
	ErrMetricMissingName  = errors.New("metric name is missing")
	ErrMetricInvalidName  = errors.New("metric name contains invalid characters")
	ErrMetricLongName     = errors.New("metric name is too long")
	ErrMetricMissingValue = errors.New("metric value is missing")
	ErrMetricInvalidValue = errors.New("metric value is invalid")

	ErrStoragePush  = errors.New("failed to push record")
	ErrStorageFetch = errors.New("failed to get record")

	ErrUnexpected = errors.New("unexpected error")
)
