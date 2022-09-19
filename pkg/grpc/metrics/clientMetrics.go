package metrics

import (
	"net/http"
	"sync"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// https://github.com/grpc-ecosystem/go-grpc-prometheus/tree/master/examples/grpc-server-with-prometheus

var (
	// 创建一个Registry
	cliReg = prometheus.NewRegistry()

	// 初始化客户端默认的metrics
	grpcClientMetrics = grpc_prometheus.NewClientMetrics()

	// 执行一次
	cliOnce sync.Once
)

func cliRegisterMetrics() {
	cliOnce.Do(func() {
		// 注册metrics才能进行采集，自定义的metrics也需要注册
		cliReg.MustRegister(grpcClientMetrics)
	})
}

// ClientHTTPService 初始化客户端的prometheus的exporter服务，使用 http://ip:port/metrics 获取数据
func ClientHTTPService(addr string) *http.Server {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: promhttp.HandlerFor(cliReg, promhttp.HandlerOpts{}),
	}

	// 启动http服务
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("listen and serve error: " + err.Error())
		}
	}()

	return httpServer
}

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientMetrics metrics unary拦截器
func UnaryClientMetrics() grpc.UnaryClientInterceptor {
	cliRegisterMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcClientMetrics.UnaryClientInterceptor()
}

// StreamClientMetrics metrics stream拦截器
func StreamClientMetrics() grpc.StreamClientInterceptor {
	cliRegisterMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcClientMetrics.StreamClientInterceptor()
}
