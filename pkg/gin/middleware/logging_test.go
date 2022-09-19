package middleware

import (
	"fmt"
	"net"
	"testing"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
)

var requestAddr string

func initServer1() {
	addr := getAddr()
	r := gin.Default()
	r.Use(RequestID())

	// 默认打印日志
	//	r.Use(Logging())

	// 自定义打印日志
	r.Use(Logging(
		WithMaxLen(400),
		//WithRequestIDFromHeader(),
		WithRequestIDFromContext(),
		//WithIgnoreRoutes("/hello"), // 忽略/hello
	))

	// 自定义zap log
	//log, _ := logger.Init(logger.WithFormat("json"))
	//r.Use(Logging(
	//	WithLog(log),
	//))

	helloFun := func(c *gin.Context) {
		logger.Info("test request id", utils.FieldRequestIDFromContext(c))
		response.Success(c, "hello world")
	}

	r.GET("/hello", helloFun)
	r.DELETE("/hello", helloFun)
	r.POST("/hello", helloFun)
	r.PUT("/hello", helloFun)
	r.PATCH("/hello", helloFun)

	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()
}

// ------------------------------------------------------------------------------------------

func TestRequest(t *testing.T) {
	initServer1()

	wantHello := "hello world"
	result := &gohttp.StdResult{}
	type User struct {
		Name string `json:"name"`
	}

	t.Run("get hello", func(t *testing.T) {
		err := gohttp.Get(result, requestAddr+"/hello", gohttp.KV{"id": "100"})
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
		err := gohttp.Delete(result, requestAddr+"/hello", gohttp.KV{"id": "100"})
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
		err := gohttp.Post(result, requestAddr+"/hello", &User{"foo"})
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
		err := gohttp.Put(result, requestAddr+"/hello", &User{"foo"})
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
		err := gohttp.Patch(result, requestAddr+"/hello", &User{"foo"})
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

func getAddr() string {
	port, _ := getAvailablePort()
	requestAddr = fmt.Sprintf("http://localhost:%d", port)
	return fmt.Sprintf(":%d", port)
}

func getAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()

	return port, err
}
