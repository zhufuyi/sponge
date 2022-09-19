package server

import (
	"time"
)

// HTTPOption 设置http
type HTTPOption func(*httpOptions)

type httpOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	isProd       bool
}

// 默认设置
func defaultHTTPOptions() *httpOptions {
	return &httpOptions{
		readTimeout:  time.Second * 60,
		writeTimeout: time.Second * 60,
		isProd:       false,
	}
}

func (o *httpOptions) apply(opts ...HTTPOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHTTPReadTimeout 设置read timeout
func WithHTTPReadTimeout(timeout time.Duration) HTTPOption {
	return func(o *httpOptions) {
		o.readTimeout = timeout
	}
}

// WithHTTPWriteTimeout 设置writer timeout
func WithHTTPWriteTimeout(timeout time.Duration) HTTPOption {
	return func(o *httpOptions) {
		o.writeTimeout = timeout
	}
}

// WithHTTPIsProd 设置是否为生产环境
func WithHTTPIsProd(IsProd bool) HTTPOption {
	return func(o *httpOptions) {
		o.isProd = IsProd
	}
}
