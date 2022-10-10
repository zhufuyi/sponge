package routers

import (
	"github.com/zhufuyi/sponge/internal/handler"
	"testing"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/gin-gonic/gin"
)

func TestNewRouter(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()

	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnableProfile = true
	config.Get().App.EnableLimit = true

	gin.SetMode(gin.ReleaseMode)
	r := NewRouter()

	userExampleRouter(r.Group("/"), handler.NewUserExampleHandler())
}
