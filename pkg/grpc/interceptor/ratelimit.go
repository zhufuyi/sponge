package interceptor

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/reugn/equalizer"
	"google.golang.org/grpc"
)

// ---------------------------------- server interceptor ----------------------------------

// RateLimitOption 日志设置
type RateLimitOption func(*rateLimitOptions)

type rateLimitOptions struct {
	qps            int           // 允许请求速度
	capacity       int           // 重新填充容量
	refillInterval time.Duration // 填充token速度，refillInterval=time.Second/qps*capacity
}

func defaultRateLimitOptions() *rateLimitOptions {
	return &rateLimitOptions{
		qps:            1000,
		capacity:       50,
		refillInterval: time.Second / 1000 * 50,
	}
}

func (o *rateLimitOptions) apply(opts ...RateLimitOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithRateLimitQPS 设置请求qps
func WithRateLimitQPS(qps int) RateLimitOption {
	return func(o *rateLimitOptions) {
		o.qps = qps
		if qps < 10 {
			o.capacity = qps
		} else if qps < 100 {
			o.capacity = 10
		} else if qps < 500 {
			o.capacity = 40
		} else if qps < 1000 {
			o.capacity = 80
		} else if qps < 2000 {
			o.capacity = 100
		} else if qps < 4000 {
			o.capacity = 200
		} else if qps < 10000 {
			o.capacity = 400
		} else {
			o.capacity = 500
		}
		o.refillInterval = time.Second / time.Duration(o.qps) * time.Duration(o.capacity)
	}
}

type myLimiter struct {
	TB *equalizer.TokenBucket // 令牌桶
}

func (m *myLimiter) Limit() bool {
	if m.TB.Ask() {
		return false
	}

	return true
}

// UnaryServerRateLimit 限流unary拦截器
func UnaryServerRateLimit(opts ...RateLimitOption) grpc.UnaryServerInterceptor {
	o := defaultRateLimitOptions()
	o.apply(opts...)

	limiter := &myLimiter{TB: equalizer.NewTokenBucket(int32(o.capacity), o.refillInterval)}
	return ratelimit.UnaryServerInterceptor(limiter)
}

// StreamServerRateLimit 限流stream拦截器
func StreamServerRateLimit(opts ...RateLimitOption) grpc.StreamServerInterceptor {
	o := defaultRateLimitOptions()
	o.apply(opts...)

	limiter := &myLimiter{equalizer.NewTokenBucket(int32(o.capacity), o.refillInterval)}
	return ratelimit.StreamServerInterceptor(limiter)
}
