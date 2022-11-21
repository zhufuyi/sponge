package swagger

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func runHTTPServer(registerSwaggerFn func(r *gin.Engine)) string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	registerSwaggerFn(r)

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

func TestDefaultHandler(t *testing.T) {
	registerSwaggerFn := func(r *gin.Engine) {
		data, err := os.ReadFile("swagger_test.json")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { recover() }()
		DefaultRouter(r, data)
	}

	requestAddr := runHTTPServer(registerSwaggerFn)
	resp, _ := http.Get(requestAddr + "/swagger/index.html")
	t.Logf("code = %d", resp.StatusCode)
}

func TestDefaultFileHandler(t *testing.T) {
	registerSwaggerFn := func(r *gin.Engine) {
		defer func() { recover() }()
		DefaultRouterByFile(r, "swagger_test.json")
	}

	requestAddr := runHTTPServer(registerSwaggerFn)
	resp, _ := http.Get(requestAddr + "/swagger/index.html")
	t.Logf("code = %d", resp.StatusCode)

	r := gin.Default()
	DefaultRouterByFile(r, "not_found.json")
}

func TestHandlers(t *testing.T) {
	registerSwaggerFn := func(r *gin.Engine) {
		data, err := os.ReadFile("swagger_test.json")
		if err != nil {
			t.Fatal(err)
		}
		CustomRouter(r, "swagger_test_1", data)
	}

	requestAddr := runHTTPServer(registerSwaggerFn)
	resp, _ := http.Get(requestAddr + "/swagger_test_1/swagger/index.html")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFileHandler(t *testing.T) {
	registerSwaggerFn := func(r *gin.Engine) {
		CustomRouterByFile(r, "swagger_test.json")
	}

	requestAddr := runHTTPServer(registerSwaggerFn)
	resp, _ := http.Get(requestAddr + "/swagger_test/swagger/index.html")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	r := gin.Default()
	CustomRouterByFile(r, "not_found.json")
}
