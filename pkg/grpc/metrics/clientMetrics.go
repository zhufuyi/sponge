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
	// client side default router
	clientPattern = "/rpc_client/metrics"

	// create a Registry
	cliReg = prometheus.NewRegistry()

	// initialise the client's default metrics
	grpcClientMetrics = grpc_prometheus.NewClientMetrics()

	cliOnce sync.Once
)

func cliRegisterMetrics() {
	cliOnce.Do(func() {
		// register metrics, including custom metrics
		cliReg.MustRegister(grpcClientMetrics)
	})
}

// SetClientPattern set the client pattern
func SetClientPattern(pattern string) {
	if pattern != "" {
		clientPattern = pattern
	}
}

// ClientRegister for http routing and grpc methods
func ClientRegister(mux *http.ServeMux) {
	// register for http routing
	mux.Handle(clientPattern, promhttp.HandlerFor(cliReg, promhttp.HandlerOpts{}))
}

// ClientHTTPService initialize the client's prometheus exporter service and use http://ip:port/metrics to fetch data
func ClientHTTPService(addr string) *http.Server {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: promhttp.HandlerFor(cliReg, promhttp.HandlerOpts{}),
	}

	// run http server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("listen and serve error: " + err.Error())
		}
	}()

	return httpServer
}

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientMetrics metrics unary interceptor
func UnaryClientMetrics() grpc.UnaryClientInterceptor {
	cliRegisterMetrics()
	return grpcClientMetrics.UnaryClientInterceptor()
}

// StreamClientMetrics metrics stream interceptor
func StreamClientMetrics() grpc.StreamClientInterceptor {
	cliRegisterMetrics()
	return grpcClientMetrics.StreamClientInterceptor()
}
