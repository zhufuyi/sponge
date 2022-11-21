package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGRPCServer(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnablePprof = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableCircuitBreaker = true

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		server := NewGRPCServer(addr,
			WithGrpcReadTimeout(time.Second),
			WithGrpcWriteTimeout(time.Second),
			WithGrpcRegistry(nil, instance),
		)
		assert.NotNil(t, server)
		cancel()
	})
}

func TestGRPCServerMock(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}
	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnablePprof = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableCircuitBreaker = true

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})

	o := defaultGrpcOptions()
	o.apply(WithGrpcRegistry(&gRegistry{}, instance))

	s := &grpcServer{
		addr:      addr,
		iRegistry: o.iRegistry,
		instance:  o.instance,
	}

	s.listen, err = net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	s.server = grpc.NewServer(s.serverOptions()...)

	go func() {
		time.Sleep(time.Second * 3)
		s.server.Stop()
	}()

	str := s.String()
	assert.NotEmpty(t, str)
	err = s.Start()
	assert.NoError(t, err)
	err = s.Stop()
	assert.NoError(t, err)
}

type gRegistry struct{}

func (g gRegistry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

func (g gRegistry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}
