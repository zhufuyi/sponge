package ratelimiter

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var l *Limiter

// Limiter is a controller for the request rate.
type Limiter struct {
	qpsLimiter sync.Map
}

// NewLimiter instantiation
func NewLimiter() *Limiter {
	return &Limiter{}
}

// GetLimiter get Limiter object, can be updated or query
func GetLimiter() *Limiter {
	return l
}

// SetLimiter set limiter parameters
// "limit" indicates the number of token buckets to be added at a rate = value/second (e.g. 10 means 1 token every 100 ms)
// "burst" the maximum instantaneous request spike allowed
func (l *Limiter) SetLimiter(limit rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		l.qpsLimiter.LoadOrStore(path, rate.NewLimiter(limit, burst))
		if !l.allow(path) {
			c.JSON(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
			return
		}
		c.Next()
	}
}

func (l *Limiter) allow(path string) bool {
	if limiter, exist := l.qpsLimiter.Load(path); exist {
		if ql, ok := limiter.(*rate.Limiter); ok && !ql.Allow() {
			return false
		}
	}

	return true
}

// UpdateQPSLimiter updates the settings for a given path's QPS limiter.
func (l *Limiter) UpdateQPSLimiter(path string, limit rate.Limit, burst int) {
	if limiter, exist := l.qpsLimiter.Load(path); exist {
		limiter.(*rate.Limiter).SetLimit(limit)
		limiter.(*rate.Limiter).SetBurst(burst)
	} else {
		l.qpsLimiter.Store(path, rate.NewLimiter(limit, burst))
	}
}

// GetQPSLimiterStatus returns the status of a given path's QPS limiter.
func (l *Limiter) GetQPSLimiterStatus(path string) (limit rate.Limit, burst int) {
	if limiter, exist := l.qpsLimiter.Load(path); exist {
		return limiter.(*rate.Limiter).Limit(), limiter.(*rate.Limiter).Burst()
	}

	return 0, 0
}

// QPS set limit qps parameters
func QPS(opts ...Option) gin.HandlerFunc {
	o := defaultOptions()
	o.apply(opts...)
	l = NewLimiter()

	return func(c *gin.Context) {
		var path string
		if !o.isIP {
			path = c.FullPath()
		} else {
			path = c.ClientIP()
		}

		l.qpsLimiter.LoadOrStore(path, rate.NewLimiter(o.qps, o.burst))
		if !l.allow(path) {
			c.Abort()
			c.JSON(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
			return
		}
		c.Next()
	}
}
