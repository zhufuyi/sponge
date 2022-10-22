package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// need real database to test
func TestHTTPServer(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}
	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnablePprof = true
	config.Get().App.EnableLimit = true

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	gin.SetMode(gin.ReleaseMode)

	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()
	server := NewHTTPServer(addr,
		WithHTTPReadTimeout(time.Second),
		WithHTTPWriteTimeout(time.Second),
		WithHTTPIsProd(true),
		WithHTTPRegistry(&iRegistry{}, &registry.ServiceInstance{}),
	)
	assert.NotNil(t, server)
}

func TestHTTPServerMock(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}
	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnablePprof = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableRegistryDiscovery = true

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)

	o := defaultHTTPOptions()
	s := &httpServer{
		addr:      addr,
		instance:  &registry.ServiceInstance{},
		iRegistry: &iRegistry{},
	}
	s.server = &http.Server{
		Addr:           addr,
		Handler:        http.NewServeMux(),
		ReadTimeout:    o.readTimeout,
		WriteTimeout:   o.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		time.Sleep(time.Second * 3)
		_ = s.server.Shutdown(context.Background())
	}()

	str := s.String()
	assert.NotEmpty(t, str)
	err = s.Start()
	assert.NoError(t, err)
	err = s.Stop()
	assert.NoError(t, err)
}

type iRegistry struct{}

func (i *iRegistry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

func (i *iRegistry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}
