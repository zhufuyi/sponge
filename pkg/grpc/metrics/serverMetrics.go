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
	// 创建一个Registry
	srvReg = prometheus.NewRegistry()

	// 初始化服务端默认的metrics
	grpcServerMetrics = grpc_prometheus.NewServerMetrics()

	// go metrics
	goMetrics = collectors.NewGoCollector()

	// 用户自定义指标 https://prometheus.io/docs/concepts/metric_types/#histogram
	customizedCounterMetrics   = []*prometheus.CounterVec{}
	customizedSummaryMetrics   = []*prometheus.SummaryVec{}
	customizedGaugeMetrics     = []*prometheus.GaugeVec{}
	customizedHistogramMetrics = []*prometheus.HistogramVec{}

	// 执行一次
	srvOnce sync.Once
)

// Option 设置metrics
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

// WithCounterMetrics 添加Counter类型指标
func WithCounterMetrics(metrics ...*prometheus.CounterVec) Option {
	return func(o *options) {
		customizedCounterMetrics = append(customizedCounterMetrics, metrics...)
	}
}

// WithSummaryMetrics 添加Summary类型指标
func WithSummaryMetrics(metrics ...*prometheus.SummaryVec) Option {
	return func(o *options) {
		customizedSummaryMetrics = append(customizedSummaryMetrics, metrics...)
	}
}

// WithGaugeMetrics 添加Gauge类型指标
func WithGaugeMetrics(metrics ...*prometheus.GaugeVec) Option {
	return func(o *options) {
		customizedGaugeMetrics = append(customizedGaugeMetrics, metrics...)
	}
}

// WithHistogramMetrics 添加Histogram类型指标
func WithHistogramMetrics(metrics ...*prometheus.HistogramVec) Option {
	return func(o *options) {
		customizedHistogramMetrics = append(customizedHistogramMetrics, metrics...)
	}
}

func srvRegisterMetrics() {
	srvOnce.Do(func() {
		// 开启了对RPCs处理时间的记录
		grpcServerMetrics.EnableHandlingTimeHistogram()

		// 注册go metrics
		srvReg.MustRegister(goMetrics)

		// 注册metrics才能进行采集，自定义的metrics也需要注册
		srvReg.MustRegister(grpcServerMetrics)

		// 注册自定义counter metric
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

// Register 注册http路由和grpc方法
func Register(mux *http.ServeMux, grpcServer *grpc.Server) {
	// 注册http路由
	mux.Handle("/metrics", promhttp.HandlerFor(srvReg, promhttp.HandlerOpts{}))

	// 注册所有gRPC方法到metrics
	grpcServerMetrics.InitializeMetrics(grpcServer)
}

// GoHTTPService 初始化服务端的prometheus的exporter服务，使用 http://ip:port/metrics 获取数据
func GoHTTPService(addr string, grpcServer *grpc.Server) *http.Server {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: promhttp.HandlerFor(srvReg, promhttp.HandlerOpts{}),
	}

	// 启动http服务
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("listen and serve error: " + err.Error())
		}
	}()

	// 所有gRPC方法初始化Metrics
	grpcServerMetrics.InitializeMetrics(grpcServer)

	return httpServer
}

// ---------------------------------- server interceptor ----------------------------------

// UnaryServerMetrics metrics unary拦截器
func UnaryServerMetrics(opts ...Option) grpc.UnaryServerInterceptor {
	o := defaultMetricsOptions()
	o.apply(opts...)
	srvRegisterMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcServerMetrics.UnaryServerInterceptor()
}

// StreamServerMetrics metrics stream拦截器
func StreamServerMetrics(opts ...Option) grpc.StreamServerInterceptor {
	o := defaultMetricsOptions()
	o.apply(opts...)
	srvRegisterMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcServerMetrics.StreamServerInterceptor()
}
