package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 需要连接连接真实数据库测试
func TestHTTPServer(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}
	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnableProfile = true
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
	)
	assert.NotNil(t, server)
}

func TestHTTPServer2(t *testing.T) {
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

	o := defaultHTTPOptions()
	s := &httpServer{
		addr: addr,
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
