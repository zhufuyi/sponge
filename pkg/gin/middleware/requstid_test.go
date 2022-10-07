package middleware

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(RequestID())

	r.GET("/hello", func(c *gin.Context) {
		response.Success(c, gin.H{"reqID": GetRequestIDFromContext(c)})
	})

	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"reqID": GetRequestIDFromHeaders(c)})
	})

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	result := &gohttp.StdResult{}
	err := gohttp.Get(result, requestAddr+"/hello")
	assert.NoError(t, err)
	t.Log(result)

	result = &gohttp.StdResult{}
	err = gohttp.Get(result, requestAddr+"/ping")
	assert.NoError(t, err)
	t.Log(result)
}
