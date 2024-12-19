package metrics

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gin/handlerfunc"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

func TestMetrics(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(Metrics(r,
		WithMetricsPath("/metrics"),
		WithIgnoreStatusCodes(http.StatusNotFound),
		WithIgnoreRequestPaths("/hello-ignore"),
		WithIgnoreRequestMethods(http.MethodDelete),
	))
	r.GET("ping", handlerfunc.Ping)
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "[get] hello")
	})

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	resp, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	resp, err = http.Get(requestAddr + "/hello")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	resp, err = http.Get(requestAddr + "/metrics")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
