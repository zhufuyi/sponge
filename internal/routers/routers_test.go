package routers

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
)

func TestNewRouter(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	config.Get().App.EnableMetrics = false
	config.Get().App.EnableTracing = true
	config.Get().App.EnablePprof = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableCircuitBreaker = true

	gin.SetMode(gin.ReleaseMode)
	r := NewRouter()

	userExampleRouter(r.Group("/"), &mock{})
}

type mock struct{}

func (u mock) Create(c *gin.Context)     { return }
func (u mock) DeleteByID(c *gin.Context) { return }
func (u mock) UpdateByID(c *gin.Context) { return }
func (u mock) GetByID(c *gin.Context)    { return }
func (u mock) ListByIDs(c *gin.Context)  { return }
func (u mock) List(c *gin.Context)       { return }
