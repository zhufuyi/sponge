package prof

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/utils"
)

func TestRegister(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, WithPrefix(""), WithPrefix("/myServer"), WithIOWaitTime())

	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
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
