package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/service"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/grpc/metrics"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/prof"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var _ app.IServer = (*grpcServer)(nil)

type grpcServer struct {
	addr   string
	server *grpc.Server
	listen net.Listener

	mux                             *http.ServeMux
	httpServer                      *http.Server
	registerMetricsMuxAndMethodFunc func() error

	iRegistry registry.Registry
	instance  *registry.ServiceInstance
}

// Start grpc service
func (s *grpcServer) Start() error {
	// registration Services
	if s.iRegistry != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) //nolint
		if err := s.iRegistry.Register(ctx, s.instance); err != nil {
			return err
		}
	}

	if s.registerMetricsMuxAndMethodFunc != nil {
		if err := s.registerMetricsMuxAndMethodFunc(); err != nil {
			return err
		}
	}

	// if either pprof or metrics is enabled, the http service will be started
	if s.mux != nil {
		addr := fmt.Sprintf(":%d", config.Get().Grpc.HTTPPort)
		s.httpServer = &http.Server{
			Addr:    addr,
			Handler: s.mux,
		}
		go func() {
			fmt.Printf("http address of pprof and metrics %s\n", addr)
			if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic("listen and serve error: " + err.Error())
			}
		}()
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

	if s.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) //nolint
		if err := s.httpServer.Shutdown(ctx); err != nil {
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
		s.registerMetricsMuxAndMethodFunc = s.registerMetricsMuxAndMethod()
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

func (s *grpcServer) registerMetricsMuxAndMethod() func() error {
	return func() error {
		if s.mux == nil {
			s.mux = http.NewServeMux()
		}
		metrics.Register(s.mux, s.server)
		return nil
	}
}

func (s *grpcServer) registerProfMux() {
	if s.mux == nil {
		s.mux = http.NewServeMux()
	}
	prof.Register(s.mux, prof.WithIOWaitTime())
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
		s.registerProfMux()
	}

	s.listen, err = net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s.server = grpc.NewServer(s.serverOptions()...)
	service.RegisterAllService(s.server) // register for all services
	return s
}
