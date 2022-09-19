package metrics

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var requestAddr string

func initGin(r *gin.Engine, metricsFun gin.HandlerFunc) {
	addr := getAddr()

	r.Use(metricsFun)

	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "[get] hello")
	})

	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()
}

func TestMetricsPath(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	metricsFun := Metrics(r,
		WithMetricsPath("/test/metrics"),
	)
	initGin(r, metricsFun)

	resp, err := http.Get(requestAddr + "/test/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("code is %d", resp.StatusCode)
	}
}

func TestIgnoreStatusCodes(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	metricsFun := Metrics(r,
		WithIgnoreStatusCodes(http.StatusNotFound),
	)
	initGin(r, metricsFun)

	_, err := http.Get(requestAddr + "/xxxxxx")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(requestAddr + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(body), `status="404"`) {
		t.Fatal("ignore request status code [404] failed")
	}
}

func TestIgnoreRequestPaths(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	metricsFun := Metrics(r,
		WithIgnoreRequestPaths("/hello"),
	)
	initGin(r, metricsFun)

	_, err := http.Get(requestAddr + "/hello")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(requestAddr + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(body), `path="/hello"`) {
		t.Fatal("ignore request paths [/hello] failed")
	}
}

func TestIgnoreRequestMethods(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	metricsFun := Metrics(r,
		WithIgnoreRequestMethods(http.MethodGet),
	)
	initGin(r, metricsFun)

	_, err := http.Get(requestAddr + "/hello")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(requestAddr + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(body), `method="GET"`) {
		t.Fatal("ignore request method [GET] failed")
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
