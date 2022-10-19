package routers

import (
	"testing"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/handler"

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
	config.Get().App.EnablePprof = true
	config.Get().App.EnableLimit = true

	gin.SetMode(gin.ReleaseMode)
	r := NewRouter()

	userExampleRouter(r.Group("/"), handler.NewUserExampleHandler())
}

type mock struct{}

func (u mock) Create(c *gin.Context) { return }

func (u mock) DeleteByID(c *gin.Context) { return }

func (u mock) UpdateByID(c *gin.Context) { return }

func (u mock) GetByID(c *gin.Context) { return }

func (u mock) ListByIDs(c *gin.Context) { return }

func (u mock) List(c *gin.Context) { return }

func Test_userExampleRouter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	userExampleRouter(r.Group("/"), &mock{})
}
