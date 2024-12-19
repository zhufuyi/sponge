package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/go-dev-frame/sponge/pkg/errcode"
	rl "github.com/go-dev-frame/sponge/pkg/shield/ratelimit"
)

// ---------------------------------- server interceptor ----------------------------------

// ErrLimitExceed is returned when the rate limiter is
// triggered and the request is rejected due to limit exceeded.
var ErrLimitExceed = rl.ErrLimitExceed

// RatelimitOption set the rate limits ratelimitOptions.
type RatelimitOption func(*ratelimitOptions)

type ratelimitOptions struct {
	window       time.Duration
	bucket       int
	cpuThreshold int64
	cpuQuota     float64
}

func defaultRatelimitOptions() *ratelimitOptions {
	return &ratelimitOptions{
		window:       time.Second * 10,
		bucket:       100,
		cpuThreshold: 800,
	}
}

func (o *ratelimitOptions) apply(opts ...RatelimitOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithWindow with window size.
func WithWindow(d time.Duration) RatelimitOption {
	return func(o *ratelimitOptions) {
		o.window = d
	}
}

// WithBucket with bucket size.
func WithBucket(b int) RatelimitOption {
	return func(o *ratelimitOptions) {
		o.bucket = b
	}
}

// WithCPUThreshold with cpu threshold
func WithCPUThreshold(threshold int64) RatelimitOption {
	return func(o *ratelimitOptions) {
		o.cpuThreshold = threshold
	}
}

// WithCPUQuota with real cpu quota(if it can not collect from process correct);
func WithCPUQuota(quota float64) RatelimitOption {
	return func(o *ratelimitOptions) {
		o.cpuQuota = quota
	}
}

// UnaryServerRateLimit server-side unary circuit breaker interceptor
func UnaryServerRateLimit(opts ...RatelimitOption) grpc.UnaryServerInterceptor {
	o := defaultRatelimitOptions()
	o.apply(opts...)
	limiter := rl.NewLimiter(
		rl.WithWindow(o.window),
		rl.WithBucket(o.bucket),
		rl.WithCPUThreshold(o.cpuThreshold),
		rl.WithCPUQuota(o.cpuQuota),
	)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		done, err := limiter.Allow()
		if err != nil {
			return nil, errcode.StatusLimitExceed.ToRPCErr(err.Error())
		}

		reply, err := handler(ctx, req)
		done(rl.DoneInfo{Err: err})
		return reply, err
	}
}

// StreamServerRateLimit server-side stream circuit breaker interceptor
func StreamServerRateLimit(opts ...RatelimitOption) grpc.StreamServerInterceptor {
	o := defaultRatelimitOptions()
	o.apply(opts...)
	limiter := rl.NewLimiter(
		rl.WithWindow(o.window),
		rl.WithBucket(o.bucket),
		rl.WithCPUThreshold(o.cpuThreshold),
		rl.WithCPUQuota(o.cpuQuota),
	)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		done, err := limiter.Allow()
		if err != nil {
			return errcode.StatusLimitExceed.ToRPCErr(err.Error())
		}

		err = handler(srv, ss)
		done(rl.DoneInfo{Err: err})
		return err
	}
}
