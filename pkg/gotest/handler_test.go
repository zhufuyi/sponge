package gotest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func newHandler() *Handler {
	var testData = map[string]interface{}{
		"1": "foo",
		"2": "bar",
	}

	// 初始化mock cache
	c := NewCache(map[string]interface{}{"no cache": testData})
	c.ICache = struct{}{} // instantiated cache interface

	// 初始化mock dao
	d := NewDao(c, testData)
	d.IDao = struct{}{} // instantiated dao interface

	// 初始化mock handler
	h := NewHandler(d, testData)
	h.IHandler = struct{}{} // instantiated handler interface

	return h
}

func TestNewHandler(t *testing.T) {
	h := newHandler()
	defer h.Close()
}

func TestHandler_GetRequestURL(t *testing.T) {
	h := newHandler()
	defer h.Close()

	h.GetRequestURL("/path")
}

func TestHandler_GoRunHttpServer(t *testing.T) {
	h := newHandler()
	defer h.Close()

	h.GoRunHttpServer([]RouterInfo{
		{
			FuncName: "Hello",
			Method:   http.MethodGet,
			Path:     "/hello",
			HandlerFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "hello world!")
			},
		},
	})

}
