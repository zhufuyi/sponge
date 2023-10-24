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
	r.Use(RequestID(
	//WithContextRequestIDKey("my_req_id"),
	//WithHeaderRequestIDKey("My-X-Req-Id"),
	))
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

		c.Set("foo", "bar")

		ctx := WrapCtx(c)
		t.Log(ctx.Value(ContextRequestIDKey))
		t.Log(GetFromCtx(ctx, "foo"))
		t.Log(CtxRequestIDField(ctx))
		t.Log(GetFromCtx(ctx, "not-exist"))
		t.Log(GetFromHeader(ctx, HeaderXRequestIDKey))
		t.Log(GetFromHeader(ctx, "Accept"))
		t.Log(GetFromHeader(ctx, "not-exist"))
		t.Log(GetFromHeaders(ctx, "Accept"))
		t.Log(GetFromHeaders(ctx, "not-exist"))
	})

	_, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)

	defer func() { recover() }()
	req, _ := http.NewRequest("GET", requestAddr+"/ping", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "text/html")
	req.Header.Set(HeaderXRequestIDKey, "2ab996de-cc03-412d-ba0a-79596efa6947")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
}

func TestGetRequestIDFromContext(t *testing.T) {
	str := GCtxRequestID(&gin.Context{})
	assert.Equal(t, "", str)
	str = CtxRequestID(context.Background())
	assert.Equal(t, "", str)
}

func TestRequestIDKeyOptions(t *testing.T) {
	opts := []RequestIDOption{
		WithContextRequestIDKey("xx"), // invalid settings
		WithContextRequestIDKey("my_req_id"),
		WithHeaderRequestIDKey("xx"), // invalid settings
		WithHeaderRequestIDKey("My-X-Req-Id"),
	}

	o := defaultRequestIDOptions()
	o.apply(opts...)
	o.setRequestIDKey()

	t.Log(ContextRequestIDKey, HeaderXRequestIDKey)

	assert.Equal(t, "my_req_id", ContextRequestIDKey)
	assert.Equal(t, "My-X-Req-Id", HeaderXRequestIDKey)
}
