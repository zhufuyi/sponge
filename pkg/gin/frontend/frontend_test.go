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
		isUseEmbedFS   = true
		htmlDir        = "user/home"
		configFile     = "user/home/config.js"
		modifyConfigFn = func(content []byte) []byte {
			return content
		}
	)

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	err := New(staticFS, isUseEmbedFS, htmlDir, configFile, modifyConfigFn).SetRouter(r)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 100)

	r = gin.New()
	gin.SetMode(gin.ReleaseMode)
	isUseEmbedFS = false
	err = New(staticFS, isUseEmbedFS, htmlDir, configFile, modifyConfigFn).SetRouter(r)
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
