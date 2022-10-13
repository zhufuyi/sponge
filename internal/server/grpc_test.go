package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/registry"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// 需要连接连接真实数据库测试
func TestGRPCServer(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnableProfile = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableRegistryDiscovery = true

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})

	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()
	server := NewGRPCServer(addr,
		WithGRPCReadTimeout(time.Second),
		WithGRPCWriteTimeout(time.Second),
		WithRegistry(nil, instance),
	)
	assert.NotNil(t, server)
}

func TestGRPCServer2(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}
	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnableProfile = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableRegistryDiscovery = true

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})

	o := defaultGRPCOptions()
	o.apply(WithRegistry(&gRegistry{}, instance))

	s := &grpcServer{
		addr:            addr,
		iRegistry:       o.iRegistry,
		serviceInstance: o.instance,
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
