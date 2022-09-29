package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/serverNameExample/config"
	"github.com/zhufuyi/sponge/pkg/registry"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestGRPCServer(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	t.Log(err)

	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:9090"})

	server := NewGRPCServer(addr,
		WithGRPCReadTimeout(time.Second),
		WithGRPCWriteTimeout(time.Second),
		WithRegistry(nil, instance),
	)
	assert.NotNil(t, server)

	str := server.String()
	assert.NotEmpty(t, str)

	go server.Start()

	time.Sleep(time.Second)
	err = server.Stop()
	assert.NoError(t, err)
}
