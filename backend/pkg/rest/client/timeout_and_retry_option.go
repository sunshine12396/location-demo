package client

import (
	"time"

	pkgerrors "github.com/pkg/errors"
)

// timeoutAndRetryOption holds timeout & retry info for client. This is optional. If not provided, it will
// pick up the default config.
type timeoutAndRetryOption struct {
	// Max num of retries. Setting to <= 0 means no retry
	// Default: 0
	maxRetries uint64
	// Max execution wait time per try.
	// Default: 15s
	maxWaitPerTry time.Duration
	// Max execution wait time, regardless of retries.
	// Default: 15s
	maxWaitInclRetries time.Duration
	// Set false to exclude retry on timeout errors.
	// Good for non-idempotent resources (i.e. push notifications)
	// Default: false
	retryOnTimeout bool
	// Retry on certain http status code
	// Default: empty
	retryOnStatusCodes map[int]bool
}

func (t *timeoutAndRetryOption) validate() error {
	if t.maxWaitPerTry > t.maxWaitInclRetries {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxWaitPerTry > maxWaitInclRetries")
	}
	if t.maxWaitPerTry < 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxWaitPerTry should not be less than zero")
	}
	if t.maxWaitInclRetries < 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxWaitInclRetries should not be less than zero")
	}
	if t.retryOnTimeout && t.maxRetries == 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxRetries should not be zero when retry onTimeout is true")
	}
	if len(t.retryOnStatusCodes) > 0 && t.maxRetries == 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxRetries should not be zero when retry onStatusCode not empty")
	}
	return nil
}
