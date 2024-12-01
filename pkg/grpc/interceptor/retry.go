package interceptor

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ---------------------------------- client interceptor ----------------------------------

var (
	// default error code for triggering a retry
	defaultErrCodes = []codes.Code{codes.Internal}
)

// RetryOption set the retry retryOptions.
type RetryOption func(*retryOptions)

type retryOptions struct {
	times    uint
	interval time.Duration
	errCodes []codes.Code
}

func defaultRetryOptions() *retryOptions {
	return &retryOptions{
		times:    2,                      // default retry times
		interval: time.Millisecond * 100, // default retry interval 100 ms
		errCodes: defaultErrCodes,        // default error code for triggering a retry
	}
}

func (o *retryOptions) apply(opts ...RetryOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithRetryTimes set number of retries, max 10
func WithRetryTimes(n uint) RetryOption {
	return func(o *retryOptions) {
		if n > 10 {
			n = 10
		}
		o.times = n
	}
}

// WithRetryInterval set the retry interval from 1 ms to 10 seconds
func WithRetryInterval(t time.Duration) RetryOption {
	return func(o *retryOptions) {
		if t < time.Millisecond {
			t = time.Millisecond
		} else if t > 10*time.Second {
			t = 10 * time.Second
		}
		o.interval = t
	}
}

// WithRetryErrCodes set the trigger retry error code
func WithRetryErrCodes(errCodes ...codes.Code) RetryOption {
	for _, errCode := range errCodes {
		switch errCode {
		case codes.Internal, codes.DeadlineExceeded, codes.Unavailable:
		default:
			defaultErrCodes = append(defaultErrCodes, errCode)
		}
	}
	return func(o *retryOptions) {
		o.errCodes = defaultErrCodes
	}
}

// UnaryClientRetry client-side retry unary interceptor
func UnaryClientRetry(opts ...RetryOption) grpc.UnaryClientInterceptor {
	o := defaultRetryOptions()
	o.apply(opts...)

	return grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithMax(o.times), // set the number of retries
		grpc_retry.WithBackoff(func(ctx context.Context, attempt uint) time.Duration { // set retry interval
			return o.interval
		}),
		grpc_retry.WithCodes(o.errCodes...), // set retry error code
	)
}

// StreamClientRetry client-side retry stream interceptor
func StreamClientRetry(opts ...RetryOption) grpc.StreamClientInterceptor {
	o := defaultRetryOptions()
	o.apply(opts...)

	return grpc_retry.StreamClientInterceptor(
		grpc_retry.WithMax(o.times), // set the number of retries
		grpc_retry.WithBackoff(func(ctx context.Context, attempt uint) time.Duration { // set retry interval
			return o.interval
		}),
		grpc_retry.WithCodes(o.errCodes...), // set retry error code
	)
}
