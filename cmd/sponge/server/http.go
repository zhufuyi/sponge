package server

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/zhufuyi/sponge/pkg/gin/handlerfunc"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/validator"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//go:embed static
var staticFS embed.FS // index.html in the static directory

// NewRouter create a router
func NewRouter(isLog ...bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	if len(isLog) > 0 && isLog[0] {
		r.Use(middleware.Logging(middleware.WithLog(logger.Get())))
	}
	binding.Validator = validator.Init()

	// solve vue using history route 404 problem
	r.NoRoute(handlerfunc.BrowserRefreshFS(staticFS, "static/index.html"))
	r.GET("/static/*filepath", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(staticFS))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	apiV1 := r.Group("/api/v1")
	apiV1.POST("/generate", GenerateCode)
	apiV1.POST("/uploadFiles", UploadFiles)
	apiV1.POST("/listTables", ListTables)
	apiV1.GET("/record/:path", GetRecord)

	return r
}

// RunHTTPServer run http server
func RunHTTPServer(addr string, isLog ...bool) {
	initRecord()

	router := NewRouter(isLog...)
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("ListenAndServe error: %v", err))
	}
}
