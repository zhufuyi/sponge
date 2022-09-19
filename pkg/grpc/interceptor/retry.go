package interceptor

import (
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ---------------------------------- client interceptor ----------------------------------

var (
	// 默认触发重试的错误码
	defaultErrCodes = []codes.Code{codes.Unavailable}
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
		times:    2,                      // 重试次数
		interval: time.Millisecond * 100, // 重试间隔100毫秒
		errCodes: defaultErrCodes,        // 默认触发重试的错误码
	}
}

func (o *retryOptions) apply(opts ...RetryOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithRetryTimes 设置重试次数，最大10次
func WithRetryTimes(n uint) RetryOption {
	return func(o *retryOptions) {
		if n > 10 {
			n = 10
		}
		o.times = n
	}
}

// WithRetryInterval 设置重试时间间隔，范围1毫秒到10秒
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

// WithRetryErrCodes 设置触发重试错误码
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

// UnaryClientRetry 重试unary拦截器
func UnaryClientRetry(opts ...RetryOption) grpc.UnaryClientInterceptor {
	o := defaultRetryOptions()
	o.apply(opts...)

	return grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithMax(o.times), // 设置重试次数
		grpc_retry.WithBackoff(func(attempt uint) time.Duration { // 设置重试间隔
			return o.interval
		}),
		grpc_retry.WithCodes(o.errCodes...), // 设置重试错误码
	)
}

// StreamClientRetry 重试stream拦截器
func StreamClientRetry(opts ...RetryOption) grpc.StreamClientInterceptor {
	o := defaultRetryOptions()
	o.apply(opts...)

	return grpc_retry.StreamClientInterceptor(
		grpc_retry.WithMax(o.times), // 设置重试次数
		grpc_retry.WithBackoff(func(attempt uint) time.Duration { // 设置重试间隔
			return o.interval
		}),
		grpc_retry.WithCodes(o.errCodes...), // 设置重试错误码
	)
}
