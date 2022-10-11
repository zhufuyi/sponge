package ratelimiter

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func runRateLimiterHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// e.g. (1) path limit, qps=500, burst=1000
	// r.Use(QPS())

	// e.g. (2) path limit, qps=50, burst=100
	r.Use(QPS(
		WithPath(),
		WithQPS(50),
		WithBurst(100),
	))

	// e.g. (3) ip limit, qps=20, burst=40
	//	r.Use(QPS(
	//		WithIP(),
	//		WithQPS(20),
	//		WithBurst(40),
	//	))

	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, "pong "+c.ClientIP())
	})

	r.GET("/hello", func(c *gin.Context) {
		response.Success(c, "hello "+c.ClientIP())
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

func TestLimiter_QPS(t *testing.T) {
	requestAddr := runRateLimiterHTTPServer()

	success, failure := 0, 0
	start := time.Now()
	for i := 0; i < 150; i++ {
		result := &gohttp.StdResult{}
		err := gohttp.Get(result, requestAddr+"/hello")
		if err != nil {
			failure++
			if failure%10 == 0 {
				fmt.Printf("%d  %v\n", i, err)
			}
		} else {
			success++
		}
	}

	end := time.Now().Sub(start).Seconds()
	t.Logf("time=%.3fs,  success=%d, failure=%d, qps=%.1f", end, success, failure, float64(success)/end)
}

func TestRateLimiter(t *testing.T) {
	requestAddr := runRateLimiterHTTPServer()

	var pingSuccess, pingFailures int32
	var helloSuccess, helloFailures int32

	for j := 0; j < 5; j++ {
		wg := &sync.WaitGroup{}
		for i := 0; i < 40; i++ {

			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				result := &gohttp.StdResult{}
				if err := gohttp.Get(result, requestAddr+"/ping"); err != nil {
					atomic.AddInt32(&pingFailures, 1)
				} else {
					atomic.AddInt32(&pingSuccess, 1)
				}
			}(i)

			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				result := &gohttp.StdResult{}
				if err := gohttp.Get(result, requestAddr+"/hello"); err != nil {
					atomic.AddInt32(&helloFailures, 1)
				} else {
					atomic.AddInt32(&helloSuccess, 1)
				}
			}(i)
		}

		wg.Wait()
		fmt.Printf("%s   helloSuccess: %d, helloFailures: %d  pingSuccess: %d, pingFailures: %d\n", time.Now().Format(time.RFC3339Nano), helloSuccess, helloFailures, pingSuccess, pingFailures)

		//time.Sleep(time.Millisecond * 200)
	}
}

func TestLimiter_GetQPSLimiterStatus(t *testing.T) {
	requestAddr := runRateLimiterHTTPServer()

	var pingSuccess, pingFailures int32

	for j := 0; j < 5; j++ {
		wg := &sync.WaitGroup{}
		for i := 0; i < 40; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				result := &gohttp.StdResult{}
				if err := gohttp.Get(result, requestAddr+"/ping"); err != nil {
					atomic.AddInt32(&pingFailures, 1)
				} else {
					atomic.AddInt32(&pingSuccess, 1)
				}
			}(i)
		}

		wg.Wait()

		qps, _ := GetLimiter().GetQPSLimiterStatus("/ping")
		fmt.Printf("%s    pingSuccess: %d, pingFailures: %d    limit:%.f\n", time.Now().Format(time.RFC3339Nano), pingSuccess, pingFailures, qps)
		//time.Sleep(time.Millisecond * 200)
	}
}

func TestLimiter_UpdateQPSLimiter(t *testing.T) {
	requestAddr := runRateLimiterHTTPServer()

	var pingSuccess, pingFailures int32

	for j := 0; j < 5; j++ {
		wg := &sync.WaitGroup{}
		for i := 0; i < 40; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				result := &gohttp.StdResult{}
				if err := gohttp.Get(result, requestAddr+"/ping"); err != nil {
					atomic.AddInt32(&pingFailures, 1)
				} else {
					atomic.AddInt32(&pingSuccess, 1)
				}
			}(i)
		}

		wg.Wait()

		limit, burst := GetLimiter().GetQPSLimiterStatus("/ping")
		GetLimiter().UpdateQPSLimiter("/ping", limit+rate.Limit(j), burst)
		fmt.Printf("%s    pingSuccess: %d, pingFailures: %d    limit:%.f\n", time.Now().Format(time.RFC3339Nano), pingSuccess, pingFailures, limit)
		//time.Sleep(time.Millisecond * 200)
	}
}

func runRateLimiterHTTPServer2() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()
	l := NewLimiter()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(l.SetLimiter(10, 20))

	r.Use(QPS(
		WithIP(),
		WithQPS(10),
		WithBurst(20),
	))

	r.GET("/hello", func(c *gin.Context) {
		response.Success(c, "hello "+c.ClientIP())
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

func TestRateLimiter2(t *testing.T) {
	requestAddr := runRateLimiterHTTPServer2()

	var pingSuccess, pingFailures int32
	var helloSuccess, helloFailures int32

	for j := 0; j < 3; j++ {
		wg := &sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				result := &gohttp.StdResult{}
				if err := gohttp.Get(result, requestAddr+"/hello"); err != nil {
					atomic.AddInt32(&helloFailures, 1)
				} else {
					atomic.AddInt32(&helloSuccess, 1)
				}
			}(i)
		}

		wg.Wait()
		fmt.Printf("%s   helloSuccess: %d, helloFailures: %d  pingSuccess: %d, pingFailures: %d\n", time.Now().Format(time.RFC3339Nano), helloSuccess, helloFailures, pingSuccess, pingFailures)

		//time.Sleep(time.Millisecond * 200)
	}
}
