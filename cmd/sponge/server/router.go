package server

import (
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/validator"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// NewRouter create a router
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.Logging(middleware.WithLog(logger.Get())))
	binding.Validator = validator.Init()

	r.POST("/generate", GenerateCode)
	r.POST("/uploadFiles", UploadFiles)
	r.POST("/listTables", ListTables)
	r.GET("/record/:path", GetRecord)

	return r
}
