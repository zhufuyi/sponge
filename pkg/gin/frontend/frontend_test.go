package frontend

import (
	"embed"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed README.md
var staticFS embed.FS

func TestFrontEnd_SetRouter(t *testing.T) {
	var (
		htmlPath       = "user/home"
		addrConfigFile = "user/home/config.js"

		defaultAddr = "http://localhost:8080"
		customAddr  = ""
	)

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	err := New(htmlPath, defaultAddr, customAddr, addrConfigFile, staticFS).SetRouter(r)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 100)

	customAddr = "http://127.0.0.1:8080"
	r = gin.New()
	gin.SetMode(gin.ReleaseMode)
	err = New(htmlPath, defaultAddr, customAddr, addrConfigFile, staticFS).SetRouter(r)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 500)

	path, _ := os.Getwd()
	err = os.RemoveAll(path + "/frontend")
	if err != nil {
		t.Error(err)
	}
}

func TestAutoOpenBrowser(t *testing.T) {
	_ = AutoOpenBrowser("http://localhost:8080")
}
