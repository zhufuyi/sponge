package middleware

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/httpcli"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"
)

func init() {
	_, _ = logger.Init()
}

func runLogHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(RequestID())

	// default Print Log
	//	r.Use(Logging())

	// custom print log
	r.Use(Logging(
		WithLog(logger.Get()),
		WithMaxLen(40),
		WithRequestIDFromHeader(),
		WithRequestIDFromContext(),
		WithIgnoreRoutes("/ping"), // ignore path /ping
	))

	// custom zap log
	//log, _ := logger.Init(logger.WithFormat("json"))
	//r.Use(Logging(
	//	WithLog(log),
	//))

	helloFun := func(c *gin.Context) {
		logger.Info("test request id", GCtxRequestIDField(c))
		response.Success(c, "hello world")
	}

	pingFun := func(c *gin.Context) {
		response.Success(c, "ping")
	}

	r.GET("/hello", helloFun)
	r.GET("/ping", pingFun)
	r.DELETE("/hello", helloFun)
	r.POST("/hello", helloFun)
	r.PUT("/hello", helloFun)
	r.PATCH("/hello", helloFun)

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)

	return requestAddr
}

func TestRequest(t *testing.T) {
	requestAddr := runLogHTTPServer()

	wantHello := "hello world"
	result := &httpcli.StdResult{}
	type User struct {
		Name string `json:"name"`
	}

	t.Run("get ping", func(t *testing.T) {
		err := httpcli.Get(result, requestAddr+"/ping")
		if err != nil {
			t.Error(err)
			return
		}
		got := result.Data.(string)
		if got != "ping" {
			t.Errorf("got: %s, want: ping", got)
		}
	})

	t.Run("get hello", func(t *testing.T) {
		err := httpcli.Get(result, requestAddr+"/hello", httpcli.WithParams(map[string]interface{}{"id": "100"}))
		if err != nil {
			t.Error(err)
			return
		}
		got := result.Data.(string)
		if got != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("delete hello", func(t *testing.T) {
		err := httpcli.Delete(result, requestAddr+"/hello", httpcli.WithParams(map[string]interface{}{"id": "100"}))
		if err != nil {
			t.Error(err)
			return
		}
		got := result.Data.(string)
		if got != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("post hello", func(t *testing.T) {
		err := httpcli.Post(result, requestAddr+"/hello", &User{"foo"})
		if err != nil {
			t.Error(err)
			return
		}
		got := result.Data.(string)
		if got != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("put hello", func(t *testing.T) {
		err := httpcli.Put(result, requestAddr+"/hello", &User{"foo"})
		if err != nil {
			t.Error(err)
			return
		}
		got := result.Data.(string)
		if got != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("patch hello", func(t *testing.T) {
		err := httpcli.Patch(result, requestAddr+"/hello", &User{"foo"})
		if err != nil {
			t.Error(err)
			return
		}
		got := result.Data.(string)
		if got != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})
}

func runLogHTTPServer2() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(RequestID())
	r.Use(SimpleLog(
		WithLog(logger.Get()),
		WithMaxLen(200),
		WithRequestIDFromContext(),
		WithRequestIDFromHeader(),
	))

	pingFun := func(c *gin.Context) {
		response.Success(c, "ping")
	}

	r.GET("/ping", pingFun)

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)

	return requestAddr
}

func TestRequest2(t *testing.T) {
	requestAddr := runLogHTTPServer2()
	result := &httpcli.StdResult{}
	t.Run("get ping", func(t *testing.T) {
		err := httpcli.Get(result, requestAddr+"/ping")
		if err != nil {
			t.Error(err)
			return
		}
		got := result.Data.(string)
		if got != "ping" {
			t.Errorf("got: %s, want: ping", got)
		}
	})
}
