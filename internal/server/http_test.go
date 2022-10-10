package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServer(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	t.Log(err)

	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()

	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnableProfile = true
	config.Get().App.EnableLimit = true

	port, _ := utils.GetAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	gin.SetMode(gin.ReleaseMode)

	server := NewHTTPServer(addr,
		WithHTTPReadTimeout(time.Second),
		WithHTTPWriteTimeout(time.Second),
		WithHTTPIsProd(true),
	)
	assert.NotNil(t, server)

	str := server.String()
	assert.NotEmpty(t, str)

	go server.Start()

	time.Sleep(time.Second)
	err = server.Stop()
	assert.NoError(t, err)
}
