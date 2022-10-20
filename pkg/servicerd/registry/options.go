package registry

// Option service instance  options
type Option func(*options)

type options struct {
	version  string
	metadata map[string]string
}

func defaultOptions() *options {
	return &options{}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithVersion set server version
func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

// WithMetadata set metadata
func WithMetadata(metadata map[string]string) Option {
	return func(o *options) {
		o.metadata = metadata
	}
}
