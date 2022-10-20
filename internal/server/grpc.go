package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/service"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/grpc/metrics"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var _ app.IServer = (*grpcServer)(nil)
var _ = pprof.Cmdline

type grpcServer struct {
	addr   string
	server *grpc.Server
	listen net.Listener

	metricsHTTPServer     *http.Server
	metricsHTTPServerFunc func() error
	pprofHTTPServerFunc   func() error

	iRegistry registry.Registry
	instance  *registry.ServiceInstance
}

// Start grpc service
func (s *grpcServer) Start() error {
	if s.iRegistry != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) //nolint
		if err := s.iRegistry.Register(ctx, s.instance); err != nil {
			return err
		}
	}

	if s.metricsHTTPServerFunc != nil {
		if err := s.metricsHTTPServerFunc(); err != nil {
			return err
		}
	}

	if s.pprofHTTPServerFunc != nil {
		if err := s.pprofHTTPServerFunc(); err != nil {
			return err
		}
	}

	if err := s.server.Serve(s.listen); err != nil { // block
		return err
	}

	return nil
}

// Stop grpc service
func (s *grpcServer) Stop() error {
	if s.iRegistry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		go func() {
			_ = s.iRegistry.Deregister(ctx, s.instance)
			cancel()
		}()
		<-ctx.Done()
	}

	s.server.GracefulStop()

	if s.metricsHTTPServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) //nolint
		if err := s.metricsHTTPServer.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

// String comment
func (s *grpcServer) String() string {
	return "grpc service address " + s.addr
}

// InitServerOptions setting up interceptors
func (s *grpcServer) serverOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		interceptor.UnaryServerRecovery(),
		interceptor.UnaryServerCtxTags(),
	}

	streamServerInterceptors := []grpc.StreamServerInterceptor{}

	// logger interceptor
	unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerLog(
		logger.Get(),
	))

	// metrics interceptor
	if config.Get().App.EnableMetrics {
		unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerMetrics())
		s.metricsHTTPServerFunc = s.metricsServer()
	}

	// limit interceptor
	if config.Get().App.EnableLimit {
		unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerRateLimit())
	}

	// circuit breaker interceptor
	if config.Get().App.EnableCircuitBreaker {
		unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerCircuitBreaker())
	}

	// trace interceptor
	if config.Get().App.EnableTracing {
		unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerTracing())
	}

	unaryServer := grpc_middleware.WithUnaryServerChain(unaryServerInterceptors...)
	streamServer := grpc_middleware.WithStreamServerChain(streamServerInterceptors...)

	options = append(options, unaryServer, streamServer)

	return options
}

func (s *grpcServer) metricsServer() func() error {
	return func() error {
		if s == nil || s.server == nil {
			return errors.New("grpc server is nil")
		}
		promAddr := fmt.Sprintf(":%d", config.Get().Grpc.MetricsPort)
		fmt.Printf("grpc metrics address %s\n", promAddr)
		s.metricsHTTPServer = metrics.GoHTTPService(promAddr, s.server)
		return nil
	}
}

func (s *grpcServer) pprofServer() func() error {
	return func() error {
		pprofAddr := fmt.Sprintf(":%d", config.Get().Grpc.PprofPort)
		fmt.Printf("grpc pprof address %s\n", pprofAddr)
		go func() {
			if err := http.ListenAndServe(pprofAddr, nil); err != nil { // default route is /debug/pprof
				panic("listen and serve error: " + err.Error())
			}
		}()
		return nil
	}
}

// NewGRPCServer creates a new grpc server
func NewGRPCServer(addr string, opts ...GrpcOption) app.IServer {
	var err error
	o := defaultGrpcOptions()
	o.apply(opts...)
	s := &grpcServer{
		addr:      addr,
		iRegistry: o.iRegistry,
		instance:  o.instance,
	}
	if config.Get().App.EnablePprof {
		s.pprofHTTPServerFunc = s.pprofServer()
	}

	s.listen, err = net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s.server = grpc.NewServer(s.serverOptions()...)
	service.RegisterAllService(s.server) // register for all services
	return s
}
