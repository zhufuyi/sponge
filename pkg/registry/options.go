package registry

// Option represents the etcd options
type Option func(*options)

type options struct {
	id       string
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

// WithID set server id
func WithID(id string) Option {
	return func(o *options) {
		o.id = id
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
