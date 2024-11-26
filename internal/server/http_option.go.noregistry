package server

// HTTPOption setting up http
type HTTPOption func(*httpOptions)

type httpOptions struct {
	isProd bool
}

func defaultHTTPOptions() *httpOptions {
	return &httpOptions{
		isProd: false,
	}
}

func (o *httpOptions) apply(opts ...HTTPOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHTTPIsProd setting up production environment markers
func WithHTTPIsProd(isProd bool) HTTPOption {
	return func(o *httpOptions) {
		o.isProd = isProd
	}
}
