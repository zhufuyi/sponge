package response

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/utils"
)

var (
	httpResponseCodes = []int{
		http.StatusOK, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden,
		http.StatusNotFound, http.StatusRequestTimeout, http.StatusConflict, http.StatusInternalServerError,
	}
	outs = []*errcode.Error{
		errcode.Success, errcode.InvalidParams, errcode.Unauthorized, errcode.InternalServerError, errcode.NotFound,
		errcode.AlreadyExists, errcode.Timeout, errcode.TooManyRequests, errcode.Forbidden,
		errcode.MethodNotAllowed, errcode.ServiceUnavailable,
	}
)

func runResponseHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/success", func(c *gin.Context) { Success(c, gin.H{"foo": "bar"}) })
	r.GET("/error", func(c *gin.Context) { Error(c, errcode.Unauthorized) })
	for _, code := range httpResponseCodes {
		code := code
		r.GET(fmt.Sprintf("/code/%d", code), func(c *gin.Context) { Output(c, code) })
	}
	for _, out := range outs {
		out := out
		r.GET(fmt.Sprintf("/out/code/%d", out.ToHTTPCode()), func(c *gin.Context) { Out(c, out) })
	}

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

func TestRespond(t *testing.T) {
	requestAddr := runResponseHTTPServer()

	result := &gohttp.StdResult{}
	err := gohttp.Get(result, requestAddr+"/success")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Data)

	result = &gohttp.StdResult{}
	err = gohttp.Get(result, requestAddr+"/error")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, result.Code)

	for _, code := range httpResponseCodes {
		result := &gohttp.StdResult{}
		url := fmt.Sprintf("%s/code/%d", requestAddr, code)
		err := gohttp.Get(result, url)
		if code == http.StatusOK {
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, result.Code)
			continue
		}
		assert.Error(t, err)
	}
	for _, out := range outs {
		result := &gohttp.StdResult{}
		url := fmt.Sprintf("%s/out/code/%d", requestAddr, out.ToHTTPCode())
		err := gohttp.Get(result, url)
		if out.ToHTTPCode() == http.StatusOK {
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, result.Code)
			continue
		}
		assert.Error(t, err)
	}
}
