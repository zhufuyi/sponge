package prof

import (
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
)

var defaultPrefix = "/debug/pprof"

// Option set defaultPrefix func
type Option func(o *options)

type options struct {
	prefix           string
	enableIOWaitTime bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithPrefix set route defaultPrefix
func WithPrefix(prefix string) Option {
	return func(o *options) {
		if prefix == "" {
			return
		}
		o.prefix = prefix
	}
}

// WithIOWaitTime enable IO wait time
func WithIOWaitTime() Option {
	return func(o *options) {
		o.enableIOWaitTime = true
	}
}

// Register pprof server mux
func Register(mux *http.ServeMux, opts ...Option) {
	o := &options{prefix: defaultPrefix}
	o.apply(opts...)

	mux.Handle(o.prefix+"/", http.HandlerFunc(pprof.Index))
	mux.Handle(o.prefix+"/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle(o.prefix+"/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle(o.prefix+"/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle(o.prefix+"/trace", http.HandlerFunc(pprof.Trace))
	mux.Handle(o.prefix+"/heap", pprof.Handler("heap"))
	mux.Handle(o.prefix+"/goroutine", pprof.Handler("goroutine"))
	mux.Handle(o.prefix+"/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle(o.prefix+"/block", pprof.Handler("block"))
	mux.Handle(o.prefix+"/mutex", pprof.Handler("mutex"))

	if o.enableIOWaitTime {
		// Similar to /profile, add IO wait time,  https://github.com/felixge/fgprof
		mux.Handle(o.prefix+"/profile-io", fgprof.Handler())
	}
}
