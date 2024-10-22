package utils

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/ex0rcist/metflix/internal/logging"
)

// Retrier service to make some retries
type Retrier struct {
	payloadFn retry.RetryableFunc
	retryIfFn retry.RetryIfFunc
	delays    []time.Duration
}

// Run, Forest
func (r Retrier) Run() error {
	return retry.Do(
		r.payloadFn,
		retry.RetryIf(r.retryIfFn),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			logging.LogWarnF("will retry after %v", r.delays[n])
			return r.delays[n]
		}),
		retry.Attempts(uint(len(r.delays))+1),
	)
}

// Retrier constructor
func NewRetrier(payloadFn func() error, retryIfFn func(err error) bool, delays []time.Duration) Retrier {
	return Retrier{
		payloadFn: payloadFn,
		retryIfFn: retryIfFn,
		delays:    delays,
	}
}
