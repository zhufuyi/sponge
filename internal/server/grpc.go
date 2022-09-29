package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/zhufuyi/sponge/internal/serverNameExample/config"
	"github.com/zhufuyi/sponge/internal/serverNameExample/service"
	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/grpc/metrics"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/registry"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var _ app.IServer = (*grpcServer)(nil)

type grpcServer struct {
	addr   string
	server *grpc.Server
	listen net.Listener

	metricsHTTPServer   *http.Server
	goRunPromHTTPServer func() error

	iRegistry       registry.Registry
	serviceInstance *registry.ServiceInstance
}

// Start grpc service
func (s *grpcServer) Start() error {
	if s.iRegistry != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) //nolint
		if err := s.iRegistry.Register(ctx, s.serviceInstance); err != nil {
			return err
		}
	}

	if s.goRunPromHTTPServer != nil {
		if err := s.goRunPromHTTPServer(); err != nil {
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
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) //nolint
		if err := s.iRegistry.Deregister(ctx, s.serviceInstance); err != nil {
			return err
		}
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
	return "grpc service, addr = " + s.addr
}

// InitServerOptions 初始化rpc设置
func (s *grpcServer) serverOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		interceptor.UnaryServerRecovery(),
		interceptor.UnaryServerCtxTags(),
	}

	streamServerInterceptors := []grpc.StreamServerInterceptor{}

	// logger 拦截器
	unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerLog(
		logger.Get(),
	))

	// metrics 拦截器
	if config.Get().App.EnableMetrics {
		unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerMetrics())
		s.goRunPromHTTPServer = func() error {
			if s == nil || s.server == nil {
				return errors.New("grpc server is nil")
			}
			promAddr := fmt.Sprintf(":%d", config.Get().Metrics.Port)
			s.metricsHTTPServer = metrics.GoHTTPService(promAddr, s.server)
			logger.Infof("start up grpc metrics service, addr = %s", promAddr)
			return nil
		}
	}

	// limit 拦截器
	if config.Get().App.EnableLimit {
		unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerRateLimit(
			interceptor.WithRateLimitQPS(config.Get().RateLimiter.QPSLimit),
		))
	}

	// trace 拦截器
	if config.Get().App.EnableTracing {
		unaryServerInterceptors = append(unaryServerInterceptors, interceptor.UnaryServerTracing())
	}

	unaryServer := grpc_middleware.WithUnaryServerChain(unaryServerInterceptors...)
	streamServer := grpc_middleware.WithStreamServerChain(streamServerInterceptors...)

	options = append(options, unaryServer, streamServer)

	return options
}

// NewGRPCServer 创建一个grpc服务
func NewGRPCServer(addr string, opts ...GRPCOption) app.IServer {
	var err error
	o := defaultGRPCOptions()
	o.apply(opts...)
	s := &grpcServer{
		addr:            addr,
		iRegistry:       o.iRegistry,
		serviceInstance: o.instance,
	}

	// 监听TCP端口
	s.listen, err = net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// 创建grpc server对象，拦截器可以在这里注入
	s.server = grpc.NewServer(s.serverOptions()...)

	// 注册所有服务
	service.RegisterAllService(s.server)

	return s
}
