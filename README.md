# retry-go

It helps to handle retries based on the backoff algorithm (w/ random jitter) as many times as you want.

## Reference

https://aws.amazon.com/ko/blogs/architecture/exponential-backoff-and-jitter/

## Install

```
go get -u github.com/jung-jin-lee/retry-go
```

## Usage

### default

```go
// function to run
f := func() error {
  numTry++
  return errors.New("error")
}
r := retry.New(f)
err := r.Run()
```

### w/ options

```go
f := func() error {
  return errors.New("error")
}
r := retry.New(f, retry.WithNumMaxRetry(10))
err := r.Run()
```

### w/ custom backoff

```go
// implement this interface
type Backoff interface {
	GetWaitTime(numRetry float64) time.Duration
}

//  and use retry w/ backoff option like this (default: exponential backoff w/ random jitter)
r := retry.New(f, retry.WithNumMaxRetry(10), retry.WithBackoff(&fakeBackoff{}))
```

### available options

```go

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
```
