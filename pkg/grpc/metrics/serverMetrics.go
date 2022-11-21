package metrics

import (
	"net/http"
	"sync"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// https://github.com/grpc-ecosystem/go-grpc-prometheus/tree/master/examples/grpc-server-with-prometheus

var (
	// create a Registry
	srvReg = prometheus.NewRegistry()

	// initialize server-side default metrics
	grpcServerMetrics = grpc_prometheus.NewServerMetrics()

	// go metrics
	goMetrics = collectors.NewGoCollector()

	// user-defined metrics https://prometheus.io/docs/concepts/metric_types/#histogram
	customizedCounterMetrics   = []*prometheus.CounterVec{}
	customizedSummaryMetrics   = []*prometheus.SummaryVec{}
	customizedGaugeMetrics     = []*prometheus.GaugeVec{}
	customizedHistogramMetrics = []*prometheus.HistogramVec{}

	srvOnce sync.Once
)

// Option set metrics
type Option func(*options)

type options struct{}

func defaultMetricsOptions() *options {
	return &options{}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithCounterMetrics add Counter type indicator
func WithCounterMetrics(metrics ...*prometheus.CounterVec) Option {
	return func(o *options) {
		customizedCounterMetrics = append(customizedCounterMetrics, metrics...)
	}
}

// WithSummaryMetrics add Summary type indicator
func WithSummaryMetrics(metrics ...*prometheus.SummaryVec) Option {
	return func(o *options) {
		customizedSummaryMetrics = append(customizedSummaryMetrics, metrics...)
	}
}

// WithGaugeMetrics add Gauge type indicator
func WithGaugeMetrics(metrics ...*prometheus.GaugeVec) Option {
	return func(o *options) {
		customizedGaugeMetrics = append(customizedGaugeMetrics, metrics...)
	}
}

// WithHistogramMetrics adding Histogram type indicators
func WithHistogramMetrics(metrics ...*prometheus.HistogramVec) Option {
	return func(o *options) {
		customizedHistogramMetrics = append(customizedHistogramMetrics, metrics...)
	}
}

func srvRegisterMetrics() {
	srvOnce.Do(func() {
		// enable time record
		grpcServerMetrics.EnableHandlingTimeHistogram()

		// register go metrics
		srvReg.MustRegister(goMetrics)

		// register metrics to capture, custom metrics also need to be registered
		srvReg.MustRegister(grpcServerMetrics)

		// register custom Counter metrics
		for _, metric := range customizedCounterMetrics {
			srvReg.MustRegister(metric)
		}
		for _, metric := range customizedSummaryMetrics {
			srvReg.MustRegister(metric)
		}
		for _, metric := range customizedGaugeMetrics {
			srvReg.MustRegister(metric)
		}
		for _, metric := range customizedHistogramMetrics {
			srvReg.MustRegister(metric)
		}
	})
}

// Register for http routing and grpc methods
func Register(mux *http.ServeMux, grpcServer *grpc.Server) {
	// register for http routing
	mux.Handle("/metrics", promhttp.HandlerFor(srvReg, promhttp.HandlerOpts{}))

	// register all gRPC methods to metrics
	grpcServerMetrics.InitializeMetrics(grpcServer)
}

// GoHTTPService initialize the prometheus exporter service on the server side and fetch data using http://ip:port/metrics
func GoHTTPService(addr string, grpcServer *grpc.Server) *http.Server {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: promhttp.HandlerFor(srvReg, promhttp.HandlerOpts{}),
	}

	// run http server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("listen and serve error: " + err.Error())
		}
	}()

	// initialising gRPC methods Metrics
	grpcServerMetrics.InitializeMetrics(grpcServer)

	return httpServer
}

// ---------------------------------- server interceptor ----------------------------------

// UnaryServerMetrics metrics unary interceptor
func UnaryServerMetrics(opts ...Option) grpc.UnaryServerInterceptor {
	o := defaultMetricsOptions()
	o.apply(opts...)
	srvRegisterMetrics()
	return grpcServerMetrics.UnaryServerInterceptor()
}

// StreamServerMetrics metrics stream interceptor
func StreamServerMetrics(opts ...Option) grpc.StreamServerInterceptor {
	o := defaultMetricsOptions()
	o.apply(opts...)
	srvRegisterMetrics()
	return grpcServerMetrics.StreamServerInterceptor()
}
