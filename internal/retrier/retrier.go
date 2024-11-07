package retrier

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/ex0rcist/metflix/internal/logging"
)

// Retrier option to configure retrier
type RetryOption func(*Retrier)

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

// RetryOption to add delays setting
func WithDelays(delays []time.Duration) RetryOption {
	return func(r *Retrier) {
		r.delays = delays
	}
}

// Constructor
func New(payloadFn func() error, retryIfFn func(err error) bool, opts ...RetryOption) *Retrier {
	r := &Retrier{
		payloadFn: payloadFn,
		retryIfFn: retryIfFn,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}
