package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zhufuyi/sponge/config"
	"github.com/zhufuyi/sponge/pkg/utils"
	"testing"
	"time"
)

func TestHTTPServer(t *testing.T) {
	err := config.Init(config.Path("conf.yml"))
	t.Log(err)

	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()

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
