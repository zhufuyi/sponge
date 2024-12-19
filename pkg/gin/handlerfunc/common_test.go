package handlerfunc

import (
	"embed"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/httpcli"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

func TestCommonHandlers(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/health", CheckHealth)
	r.GET("/ping", Ping)
	r.GET("/codes", ListCodes)

	go func() {
		_ = r.Run(serverAddr)
	}()

	time.Sleep(time.Millisecond * 200)
	resp, err := http.Get(requestAddr + "/health")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	resp, err = http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	resp, err = http.Get(requestAddr + "/codes")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	time.Sleep(time.Second)
}

func TestBrowserRefresh(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(BrowserRefresh("README.md"))

	go func() {
		_ = r.Run(serverAddr)
	}()

	time.Sleep(time.Millisecond * 200)
	resp, err := http.Get(requestAddr + "/notfound")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	req := httpcli.New()
	req.SetURL(requestAddr + "/notfound")
	req.SetHeader("Accept", "text/html")
	_, err = req.GET()
	assert.NoError(t, err)
}

func TestBrowserRefresh2(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(BrowserRefresh("notfound"))

	go func() {
		_ = r.Run(serverAddr)
	}()

	time.Sleep(time.Millisecond * 200)

	req := httpcli.New()
	req.SetURL(requestAddr + "/notfound")
	req.SetHeader("Accept", "text/html")
	_, err := req.GET()
	assert.NoError(t, err)
}

//go:embed README.md
var readmeFS embed.FS

func TestBrowserRefreshFS(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(BrowserRefreshFS(readmeFS, "README.md"))

	go func() {
		_ = r.Run(serverAddr)
	}()

	time.Sleep(time.Millisecond * 200)
	resp, err := http.Get(requestAddr + "/notfound")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	req := httpcli.New()
	req.SetURL(requestAddr + "/notfound")
	req.SetHeader("Accept", "text/html")
	_, err = req.GET()
	assert.NoError(t, err)
}

func TestBrowserRefreshFS2(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(BrowserRefreshFS(readmeFS, "notfound"))

	go func() {
		_ = r.Run(serverAddr)
	}()

	time.Sleep(time.Millisecond * 200)

	req := httpcli.New()
	req.SetURL(requestAddr + "/notfound")
	req.SetHeader("Accept", "text/html")
	_, err := req.GET()
	assert.NoError(t, err)
}
