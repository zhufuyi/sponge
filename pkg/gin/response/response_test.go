package response

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var httpResponseCodes = []int{
	http.StatusOK, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden,
	http.StatusNotFound, http.StatusRequestTimeout, http.StatusConflict, http.StatusInternalServerError,
}

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

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

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
}
