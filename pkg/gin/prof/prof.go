package prof

import (
	"net/http/pprof"

	"github.com/felixge/fgprof"
	"github.com/gin-gonic/gin"
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

// Register pprof for gin router
func Register(r *gin.Engine, opts ...Option) {
	o := &options{prefix: defaultPrefix}
	o.apply(opts...)

	group := r.Group(o.prefix)

	group.GET("/", gin.WrapF(pprof.Index))
	group.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	group.GET("/profile", gin.WrapF(pprof.Profile))
	group.POST("/symbol", gin.WrapF(pprof.Symbol))
	group.GET("/symbol", gin.WrapF(pprof.Symbol))
	group.GET("/trace", gin.WrapF(pprof.Trace))
	group.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	group.GET("/block", gin.WrapH(pprof.Handler("block")))
	group.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	group.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	group.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	group.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))

	if o.enableIOWaitTime {
		// Similar to /profile, add IO wait time,  https://github.com/felixge/fgprof
		group.GET("/profile-io", gin.WrapH(fgprof.Handler()))
	}
}
