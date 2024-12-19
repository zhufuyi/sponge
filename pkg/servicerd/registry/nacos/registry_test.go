package nacos

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/servicerd/registry"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

func TestNewRegistry(t *testing.T) {
	nacosIPAddr := "192.168.3.37"
	nacosPort := 8848
	nacosNamespaceID := "3454d2b5-2455-4d0e-bf6d-e033b086bb4c"

	id := "serverName_192.168.3.37"
	instanceName := "serverName"
	instanceEndpoints := []string{"grpc://192.168.3.27:8282"}

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		iRegistry, instance, err := NewRegistry(nacosIPAddr, nacosPort, nacosNamespaceID, id, instanceName, instanceEndpoints)
		if err != nil {
			t.Log(err)
			return
		}
		t.Log(iRegistry, instance)
	})
}

func newNacosRegistry() *Registry {
	return New(getCli(),
		WithPrefix("/micro"),
		WithWeight(1),
		WithCluster("cluster"),
		WithGroup("dev"),
		WithDefaultKind("grpc"),
	)
}

func TestRegistry(t *testing.T) {
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})
	r := &Registry{}
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		r = newNacosRegistry()
	})

	go func() {
		defer func() { recover() }()
		_, err := r.Watch(context.Background(), "foo")
		t.Log(err)
	}()

	defer func() { recover() }()
	time.Sleep(time.Millisecond * 10)
	err := r.Register(context.Background(), instance)
	t.Log(err)
}

func TestDeregister(t *testing.T) {
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})
	r := &Registry{}
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		r = newNacosRegistry()
	})

	defer func() { recover() }()
	err := r.Deregister(context.Background(), instance)
	t.Log(err)
}

func TestGetService(t *testing.T) {
	r := &Registry{}
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		r = newNacosRegistry()
	})

	defer func() { recover() }()
	_, err := r.GetService(context.Background(), "foo")
	t.Log(err)
}

func TestRegistry_RegisterError(t *testing.T) {
	instance := registry.NewServiceInstance("", "", []string{"grpc://127.0.0.1:8282"})
	r := &Registry{}
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		r = newNacosRegistry()
	})
	defer func() { recover() }()

	err := r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.2:8282"},
		registry.WithMetadata(map[string]string{
			"foo2": "bar2",
		}))
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", "bar", []string{"127.0.0.1:port"})
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", "bar", []string{"127.0.0.1"})
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)
}
