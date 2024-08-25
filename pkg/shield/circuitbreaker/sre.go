package circuitbreaker

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zhufuyi/sponge/pkg/shield/window"
)

// Option is sre breaker option function.
type Option func(*options)

const (
	// StateOpen when circuit breaker open, request not allowed, after sleep
	// some duration, allow one single request for testing the health, if ok
	// then state reset to closed, if not continue the step.
	StateOpen int32 = iota
	// StateClosed when circuit breaker closed, request allowed, the breaker
	// calc the succeed ratio, if request num greater request setting and
	// ratio lower than the setting ratio, then reset state to open.
	StateClosed
)

var (
	_ CircuitBreaker = &Breaker{}
)

// options is a breaker options.
type options struct {
	success float64
	request int64
	bucket  int
	window  time.Duration
}

// WithSuccess with the K = 1 / Success value of sre breaker, default success is 0.5
// Reducing the K will make adaptive throttling behave more aggressively,
// Increasing the K will make adaptive throttling behave less aggressively.
func WithSuccess(s float64) Option {
	return func(c *options) {
		c.success = s
	}
}

// WithRequest with the minimum number of requests allowed.
func WithRequest(r int64) Option {
	return func(c *options) {
		c.request = r
	}
}

// WithWindow with the duration size of the statistical window.
func WithWindow(d time.Duration) Option {
	return func(c *options) {
		c.window = d
	}
}

// WithBucket set the bucket number in a window duration.
func WithBucket(b int) Option {
	return func(c *options) {
		c.bucket = b
	}
}

// Breaker is a sre CircuitBreaker pattern.
type Breaker struct {
	stat window.RollingCounter
	r    *rand.Rand
	// rand.New(...) returns a non thread safe object
	randLock sync.Mutex

	// Reducing the k will make adaptive throttling behave more aggressively,
	// Increasing the k will make adaptive throttling behave less aggressively.
	k       float64
	request int64

	state int32
}

// NewBreaker return a sreBresker with options
func NewBreaker(opts ...Option) CircuitBreaker {
	opt := options{
		success: 0.6,
		request: 100,
		bucket:  10,
		window:  3 * time.Second,
	}
	for _, o := range opts {
		o(&opt)
	}
	counterOpts := window.RollingCounterOpts{
		Size:           opt.bucket,
		BucketDuration: time.Duration(int64(opt.window) / int64(opt.bucket)),
	}
	stat := window.NewRollingCounter(counterOpts)
	return &Breaker{
		stat:    stat,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		request: opt.request,
		k:       1 / opt.success,
		state:   StateClosed,
	}
}

func (b *Breaker) summary() (success int64, total int64) {
	b.stat.Reduce(func(iterator window.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			total += bucket.Count
			for _, p := range bucket.Points {
				success += int64(p)
			}
		}
		return 0
	})
	return //nolint
}

// Allow request if error returns nil.
func (b *Breaker) Allow() error {
	// The number of requests accepted by the backend
	accepts, total := b.summary()
	// The number of requests attempted by the application layer(at the client, on top of the adaptive throttling system)
	requests := b.k * float64(accepts)
	// check overflow requests = K * accepts
	if total < b.request || float64(total) < requests {
		atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClosed)
		return nil
	}
	atomic.CompareAndSwapInt32(&b.state, StateClosed, StateOpen)
	dr := math.Max(0, (float64(total)-requests)/float64(total+1))
	drop := b.trueOnProba(dr)
	if drop {
		return ErrNotAllowed
	}
	return nil
}

// MarkSuccess mark request is success.
func (b *Breaker) MarkSuccess() {
	b.stat.Add(1)
}

// MarkFailed mark request is failed.
func (b *Breaker) MarkFailed() {
	// NOTE: when client reject request locally, keep adding counter let the drop ratio higher.
	b.stat.Add(0)
}

func (b *Breaker) trueOnProba(proba float64) (truth bool) {
	b.randLock.Lock()
	truth = b.r.Float64() < proba
	b.randLock.Unlock()
	return truth
}
