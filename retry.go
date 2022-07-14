package retry

import (
	"time"
)

func New(f RetryFunc, opts ...RetryOption) *Retry {
	r := &Retry{
		f: f,
		backoff: newExponentialBackoff(1, 64),
		numMaxRetry: 2,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type Retry struct {
	f RetryFunc
	numMaxRetry int
	backoff Backoff
	retryableErrs []error
}

func (r *Retry) Run() error {
	var err error
	for i := 0; i <= r.numMaxRetry; i++ {
		err = r.f()
		if err == nil {
			return nil
		}
		if !r.isRetryableErr(err) {
			return err
		}

		waitTime := r.backoff.GetWaitTime(float64(i))
		time.Sleep(waitTime)
	}

	return err
}

func (r *Retry) isRetryableErr(err error) bool {
	if len(r.retryableErrs) == 0 {
		return true
	}

	for _, rerr := range r.retryableErrs {
		if err.Error() == rerr.Error() {
			return true
		}
	}

	return false
}

type RetryFunc func() error

type RetryOption func(*Retry)

func WithNumMaxRetry(numMaxRetry int) RetryOption {
	return func(r *Retry) {
		r.numMaxRetry = numMaxRetry
	}
}

func WithBackoff(backoff Backoff) RetryOption {
	return func(r *Retry) {
		r.backoff = backoff
	}
}

func WithRetryableErrs(retryableErrs []error) RetryOption {
	return func(r *Retry) {
		r.retryableErrs = retryableErrs
	}
}