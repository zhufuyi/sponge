package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func runRequestIDHTTPServer(fn func(c *gin.Context)) string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(RequestID())
	r.GET("/ping", func(c *gin.Context) {
		fn(c)
		c.String(200, "pong")
	})

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

func TestFieldRequestIDFromContext(t *testing.T) {
	requestAddr := runRequestIDHTTPServer(func(c *gin.Context) {
		str := GCtxRequestID(c)
		t.Log(str)
		field := GCtxRequestIDField(c)
		t.Log(field)

		str = HeaderRequestID(c)
		t.Log(str)
		field = HeaderRequestIDField(c)
		t.Log(field)

		str = CtxRequestID(c)
		t.Log(str)
		field = CtxRequestIDField(c)
		t.Log(field)
	})

	_, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
}

func TestGetRequestIDFromContext(t *testing.T) {
	str := GCtxRequestID(&gin.Context{})
	assert.Equal(t, "", str)
	str = CtxRequestID(context.Background())
	assert.Equal(t, "", str)
}
