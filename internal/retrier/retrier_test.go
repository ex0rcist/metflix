package retrier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetrier_Success(t *testing.T) {
	payloadFn := func() error { return nil }
	retryIfFn := func(err error) bool { return false }

	delays := []time.Duration{time.Millisecond, time.Millisecond * 2}
	retrier := New(payloadFn, retryIfFn, WithDelays(delays))

	err := retrier.Run(context.Background())
	assert.NoError(t, err)
}

func TestRetrier_RetrySuccess(t *testing.T) {
	attempts := 0
	payloadFn := func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	}
	retryIfFn := func(err error) bool { return true }

	delays := []time.Duration{time.Millisecond, time.Millisecond * 2}
	retrier := New(payloadFn, retryIfFn, WithDelays(delays))

	err := retrier.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, attempts)
}

func TestRetrier_Failure(t *testing.T) {
	attempts := 0
	payloadFn := func() error {
		attempts++
		return errors.New("persistent error")
	}
	retryIfFn := func(err error) bool { return true }

	delays := []time.Duration{time.Millisecond, time.Millisecond * 2, time.Millisecond * 3}
	retrier := New(payloadFn, retryIfFn, WithDelays(delays))

	err := retrier.Run(context.Background())
	assert.Error(t, err)
	assert.Equal(t, len(delays)+1, attempts)
}
