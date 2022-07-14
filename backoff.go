package retry

import (
	"math"
	"math/rand"
	"time"
)

type Backoff interface {
	GetWaitTime(numRetry float64) time.Duration
}

func newExponentialBackoff(base, cap float64) *exponentialBackoff {
	if base < 1 {
		base = 1
	}

	return &exponentialBackoff{base: base, cap: cap}
}

type exponentialBackoff struct {
	base float64
	cap float64
}

func (bo *exponentialBackoff) GetWaitTime(numRetry float64) time.Duration {
	if numRetry < 0 {
		numRetry = 0
	}

	backoff := bo.getBackoff(numRetry)
	randomJitter := bo.getRandomJitter()
	sleep := backoff + randomJitter
	waitDuration := time.Duration(sleep * 1000) * time.Millisecond

	return waitDuration
}

func (bo *exponentialBackoff) getBackoff(numRetry float64) float64 {
	return math.Min(bo.cap, bo.base * math.Pow(2, numRetry))
}

func (bo *exponentialBackoff) getRandomJitter() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	return r.Float64()
}