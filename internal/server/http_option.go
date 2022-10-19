package server

import (
	"time"
)

// HTTPOption setting up http
type HTTPOption func(*httpOptions)

type httpOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	isProd       bool
}

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

// WithHTTPReadTimeout setting up read timeout
func WithHTTPReadTimeout(timeout time.Duration) HTTPOption {
	return func(o *httpOptions) {
		o.readTimeout = timeout
	}
}

// WithHTTPWriteTimeout setting up writer timeout
func WithHTTPWriteTimeout(timeout time.Duration) HTTPOption {
	return func(o *httpOptions) {
		o.writeTimeout = timeout
	}
}

// WithHTTPIsProd setting up production environment markers
func WithHTTPIsProd(IsProd bool) HTTPOption {
	return func(o *httpOptions) {
		o.isProd = IsProd
	}
}
