package gofile

const (
	prefix  = "prefix"
	suffix  = "suffix"
	contain = "contain"
)

var (
	defaultFilterType = "" // 有prefix、suffix、contain，默认不过滤
)

type options struct {
	filter string
	name   string
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

// WithSuffix 后缀匹配
func WithSuffix(name string) Option {
	return func(o *options) {
		o.filter = suffix
		o.name = name
	}
}

// WithPrefix 前缀匹配
func WithPrefix(name string) Option {
	return func(o *options) {
		o.filter = prefix
		o.name = name
	}
}

// WithContain 包含字符串
func WithContain(name string) Option {
	return func(o *options) {
		o.filter = contain
		o.name = name
	}
}
