package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWaitTime(t *testing.T) {
	backoff := newExponentialBackoff(1, 64)
	second := time.Second

	waitTime1 := backoff.GetWaitTime(1)
	assert.Condition(t, func() bool {
		return 2*second <= waitTime1 && waitTime1 < 3*second
	})

	waitTime2 := backoff.GetWaitTime(2)

	assert.Condition(t, func() bool {
		return 4*second <= waitTime2 && waitTime2 < 5*second
	})

	waitTime3 := backoff.GetWaitTime(3)

	assert.Condition(t, func() bool {
		return 8*second <= waitTime3 && waitTime3 < 9*second
	})
}