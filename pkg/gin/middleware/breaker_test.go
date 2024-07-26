package middleware

import (
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhufuyi/sponge/pkg/container/group"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/httpcli"
	"github.com/zhufuyi/sponge/pkg/shield/circuitbreaker"
	"github.com/zhufuyi/sponge/pkg/utils"
)

func runCircuitBreakerHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(CircuitBreaker(WithGroup(group.NewGroup(func() interface{} {
		return circuitbreaker.NewBreaker()
	})),
		WithValidCode(http.StatusForbidden),
	))

	r.GET("/hello", func(c *gin.Context) {
		if rand.Int()%2 == 0 {
			response.Output(c, http.StatusInternalServerError)
		} else {
			response.Success(c, "hello "+c.ClientIP())
		}
	})

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

func TestCircuitBreaker(t *testing.T) {
	requestAddr := runCircuitBreakerHTTPServer()

	var success, failures, countBreaker int32
	for j := 0; j < 5; j++ {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				result := &httpcli.StdResult{}
				if err := httpcli.Get(result, requestAddr+"/hello"); err != nil {
					if strings.Contains(err.Error(), ErrNotAllowed.Error()) {
						atomic.AddInt32(&countBreaker, 1)
					}
					atomic.AddInt32(&failures, 1)
				} else {
					atomic.AddInt32(&success, 1)
				}
			}
		}()

		wg.Wait()
		t.Logf("%s   success: %d, failures: %d, breakerOpen: %d\n",
			time.Now().Format(time.RFC3339Nano), success, failures, countBreaker)
	}
}
