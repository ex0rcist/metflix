package entities

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrBadAddressFormat = errors.New("bad net address format")

	ErrRecordNotFound     = errors.New("metric not found")
	ErrMetricUnknown      = errors.New("unknown metric type")
	ErrMetricReport       = errors.New("metric report error")
	ErrMetricMissingName  = errors.New("metric name is missing")
	ErrMetricInvalidName  = errors.New("metric name contains invalid characters")
	ErrMetricLongName     = errors.New("metric name is too long")
	ErrMetricMissingValue = errors.New("metric value is missing")
	ErrMetricInvalidValue = errors.New("metric value is invalid")

	ErrStoragePush       = errors.New("failed to push record")
	ErrStorageFetch      = errors.New("failed to get record")
	ErrStorageUnpingable = errors.New("healthcheck is not supported")

	ErrEncodingInternal    = errors.New("internal encoding error")
	ErrEncodingUnsupported = errors.New("requsted encoding is not supported")

	ErrNoSignature = errors.New("no signature provided")

	ErrUnexpected = errors.New("unexpected error")
)

// Constructor wrapper.
func NewStackError(err error) error {
	return errors.New(err.Error())
}

var _ error = (*RetriableError)(nil)

// Error to handle retries.
type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

// Return readable representation.
func (e RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}
