package rpcclient

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNewServerNameExampleRPCConn(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer func() { recover() }()
		time.Sleep(time.Millisecond * 100)
		config.Get().GrpcClient[0].RegistryDiscoveryType = "consul"
		NewServerNameExampleRPCConn()
	}()

	go func() {
		defer func() { recover() }()
		time.Sleep(time.Millisecond * 200)
		config.Get().GrpcClient[0].RegistryDiscoveryType = "etcd"
		NewServerNameExampleRPCConn()
	}()

	go func() {
		defer func() { recover() }()
		time.Sleep(time.Millisecond * 300)
		config.Get().GrpcClient[0].RegistryDiscoveryType = "nacos"
		NewServerNameExampleRPCConn()
	}()

	go func() {
		defer func() { recover() }()
		time.Sleep(time.Millisecond * 400)
		config.Get().GrpcClient[0].Name = "unknown name"
		NewServerNameExampleRPCConn()
	}()

	time.Sleep(time.Second * 6)

	go func() {
		defer func() { recover() }()
		conn := GetServerNameExampleRPCConn()
		assert.NotNil(t, conn)
	}()

	err = CloseServerNameExampleRPCConn()
	assert.NoError(t, err)
}

func TestGetServerNameExampleRPCConn(t *testing.T) {
	serverNameExampleConn = nil
	err := CloseServerNameExampleRPCConn()
	assert.NoError(t, err)

	// nil error
	defer func() { recover() }()
	conn := GetServerNameExampleRPCConn()
	assert.NotNil(t, conn)
}
