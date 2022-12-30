package server

import (
	"embed"
	"net/http"

	"github.com/zhufuyi/sponge/pkg/gin/handlerfunc"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//go:embed static
var staticFS embed.FS // index.html in the static directory

// NewRouter create a router
func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	//r.Use(middleware.Logging(middleware.WithLog(logger.Get())))
	binding.Validator = validator.Init()

	// solve vue using history route 404 problem
	r.NoRoute(handlerfunc.BrowserRefreshFS(staticFS, "static/index.html"))
	r.GET("/static/*filepath", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(staticFS))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/generate", GenerateCode)
	r.POST("/uploadFiles", UploadFiles)
	r.POST("/listTables", ListTables)
	r.GET("/record/:path", GetRecord)

	return r
}
