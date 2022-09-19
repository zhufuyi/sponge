package ratelimiter

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

var requestAddr string

func init() {
	addr := getAddr()
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
		c.JSON(200, "pong "+c.ClientIP())
	})

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, "hello "+c.ClientIP())
	})

	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()
}

func TestLimiter_QPS(t *testing.T) {
	success, failure := 0, 0
	start := time.Now()
	for i := 0; i < 1000; i++ {
		err := get(requestAddr + "/hello")
		if err != nil {
			failure++
			if failure%10 == 0 {
				fmt.Printf("%d  %v\n", i, err)
			}
		} else {
			success++
		}
		time.Sleep(time.Millisecond) // 间隔1毫秒
	}
	time := time.Now().Sub(start).Seconds()
	t.Logf("time=%.3fs,  success=%d, failure=%d, qps=%.1f", time, success, failure, float64(success)/time)
}

func TestRateLimiter(t *testing.T) {
	var pingSuccess, pingFailures int32
	var helloSuccess, helloFailures int32

	for j := 0; j < 10; j++ {
		wg := &sync.WaitGroup{}
		for i := 0; i < 20; i++ {

			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if err := get(requestAddr + "/ping"); err != nil {
					atomic.AddInt32(&pingFailures, 1)
				} else {
					atomic.AddInt32(&pingSuccess, 1)
				}
			}(i)

			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if err := get(requestAddr + "/hello"); err != nil {
					atomic.AddInt32(&helloFailures, 1)
				} else {
					atomic.AddInt32(&helloSuccess, 1)
				}
			}(i)
		}

		wg.Wait()
		fmt.Printf("%s   helloSuccess: %d, helloFailures: %d  pingSuccess: %d, pingFailures: %d\n", time.Now().Format(time.RFC3339Nano), helloSuccess, helloFailures, pingSuccess, pingFailures)

		time.Sleep(time.Millisecond * 200)
	}
}

func TestLimiter_GetQPSLimiterStatus(t *testing.T) {
	var pingSuccess, pingFailures int32

	for j := 0; j < 10; j++ {
		wg := &sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if err := get(requestAddr + "/ping"); err != nil {
					atomic.AddInt32(&pingFailures, 1)
				} else {
					atomic.AddInt32(&pingSuccess, 1)
				}
			}(i)
		}

		wg.Wait()

		qps, _ := GetLimiter().GetQPSLimiterStatus("/ping")
		fmt.Printf("%s    pingSuccess: %d, pingFailures: %d    limit:%.f\n", time.Now().Format(time.RFC3339Nano), pingSuccess, pingFailures, qps)
		time.Sleep(time.Millisecond * 200)
	}
}

func TestLimiter_UpdateQPSLimiter(t *testing.T) {
	var pingSuccess, pingFailures int32

	for j := 0; j < 10; j++ {
		wg := &sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if err := get(requestAddr + "/ping"); err != nil {
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
		time.Sleep(time.Millisecond * 200)
	}
}

func getAddr() string {
	port, _ := getAvailablePort()
	requestAddr = fmt.Sprintf("http://localhost:%d", port)
	return fmt.Sprintf(":%d", port)
}

func getAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()

	return port, err
}

func get(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(resp.StatusCode))
	}

	return nil
}
