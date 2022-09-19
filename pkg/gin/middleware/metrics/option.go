package metrics

import (
	"strings"
)

// Option set the metrics options.
type Option func(*options)

type options struct {
	metricsPath          string
	ignoreStatusCodes    map[int]struct{}
	ignoreRequestPaths   map[string]struct{}
	ignoreRequestMethods map[string]struct{}
}

// defaultOptions default value
func defaultOptions() *options {
	return &options{
		metricsPath:          "/metrics",
		ignoreStatusCodes:    nil,
		ignoreRequestPaths:   nil,
		ignoreRequestMethods: nil,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithMetricsPath set metrics path
func WithMetricsPath(metricsPath string) Option {
	return func(o *options) {
		o.metricsPath = metricsPath
	}
}

// WithIgnoreStatusCodes ignore status codes
func WithIgnoreStatusCodes(statusCodes ...int) Option {
	return func(o *options) {
		codeMaps := make(map[int]struct{}, len(statusCodes))
		for _, code := range statusCodes {
			codeMaps[code] = struct{}{}
		}
		o.ignoreStatusCodes = codeMaps
	}
}

// WithIgnoreRequestPaths ignore request paths
func WithIgnoreRequestPaths(paths ...string) Option {
	return func(o *options) {
		pathMaps := make(map[string]struct{}, len(paths))
		for _, path := range paths {
			pathMaps[path] = struct{}{}
		}
		o.ignoreRequestPaths = pathMaps
	}
}

// WithIgnoreRequestMethods ignore request methods
func WithIgnoreRequestMethods(methods ...string) Option {
	return func(o *options) {
		methodMaps := make(map[string]struct{}, len(methods))
		for _, method := range methods {
			methodMaps[strings.ToUpper(method)] = struct{}{}
		}
		o.ignoreRequestMethods = methodMaps
	}
}

func (o *options) isIgnoreCodeStatus(statusCode int) bool {
	if o.ignoreStatusCodes == nil {
		return false
	}
	_, ok := o.ignoreStatusCodes[statusCode]
	return ok
}

func (o *options) isIgnorePath(path string) bool {
	if o.ignoreRequestPaths == nil {
		return false
	}
	_, ok := o.ignoreRequestPaths[path]
	return ok
}

func (o *options) checkIgnoreMethod(method string) bool {
	if o.ignoreRequestMethods == nil {
		return false
	}
	_, ok := o.ignoreRequestMethods[strings.ToUpper(method)]
	return ok
}
