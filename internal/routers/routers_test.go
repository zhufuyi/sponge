package routers

import (
	"testing"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/serverNameExample/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	t.Log(err)

	defer func() {
		if e := recover(); e != nil {
			t.Log("ignore connect mysql error info")
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := NewRouter()
	assert.NotNil(t, r)
}
