package prof

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	r := gin.Default()
	Register(r, WithPrefix(""), WithPrefix("/myServer"), WithIOWaitTime())

	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("listen and serve error: " + err.Error())
		}
	}()
	time.Sleep(time.Millisecond * 200)

	resp, err := http.Get(requestAddr + "/myServer")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
