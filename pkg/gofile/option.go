package gofile

const (
	prefix  = "prefix"
	suffix  = "suffix"
	contain = "contain"
)

var (
	defaultFilterType = "" // with prefix, suffix, contain, no filter by default
)

type options struct {
	filter string
	name   string

	noAbsolutePath bool
}

func defaultOptions() *options {
	return &options{
		filter: defaultFilterType,
	}
}

// Option set the file options.
type Option func(*options)

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithSuffix set suffix matching
func WithSuffix(name string) Option {
	return func(o *options) {
		o.filter = suffix
		o.name = name
	}
}

// WithPrefix set prefix matching
func WithPrefix(name string) Option {
	return func(o *options) {
		o.filter = prefix
		o.name = name
	}
}

// WithContain set contain matching
func WithContain(name string) Option {
	return func(o *options) {
		o.filter = contain
		o.name = name
	}
}

// WithNoAbsolutePath set no absolute path
func WithNoAbsolutePath() Option {
	return func(o *options) {
		o.noAbsolutePath = true
	}
}
