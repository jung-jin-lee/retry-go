package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	// given
	numTry := 0
	f := func() error {
		numTry++
		return errors.New("error")
	}
	r := New(f)
	expected := 3

	// when
	err := r.Run()

	// then
	assert.Error(t, err)
	assert.Equal(t, expected, numTry)
}

func TestRetryWith(t *testing.T) {
	// given
	numTry := 0
	f := func() error {
		numTry++
		return errors.New("error")
	}
	r := New(f, WithNumMaxRetry(10), WithBackoff(&fakeBackoff{}))
	expected := 11

	// when
	err := r.Run()

	// then
	assert.Error(t, err)
	assert.Equal(t, expected, numTry)
}

func TestRetryWithRetryableErrs(t *testing.T) {
	// given
	retryableErr1 := errors.New("retryable err1")
	retryableErr2 := errors.New("retryable err2")
	retryableErrs := []error{retryableErr1, retryableErr2}
	numTry := 0
	f := func() error {
		numTry++
		return retryableErr2
	}
	r := New(f, WithRetryableErrs(retryableErrs), WithNumMaxRetry(3), WithBackoff(&fakeBackoff{}))
	expected := 4

	// when
	err := r.Run()

	// then
	assert.Error(t, err)
	assert.Equal(t, expected, numTry)
}

func TestRetryWithNonRetryableErrs(t *testing.T) {
	// given
	nonRetryableErr := errors.New("non retryable err")
	retryableErr1 := errors.New("retryable err1")
	retryableErr2 := errors.New("retryable err2")
	retryableErrs := []error{retryableErr1, retryableErr2}
	numTry := 0
	f := func() error {
		numTry++
		return nonRetryableErr
	}
	r := New(f, WithRetryableErrs(retryableErrs), WithNumMaxRetry(3), WithBackoff(&fakeBackoff{}))
	expected := 1

	// when
	err := r.Run()

	// then
	assert.Error(t, err)
	assert.Equal(t, expected, numTry)
}

type fakeBackoff struct {}

func (bo *fakeBackoff) GetWaitTime(numRetry float64) time.Duration {
	return 1*time.Microsecond
}