package swagger

import (
	"fmt"
	"os"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gofile"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

// DefaultRouter default swagger router, request url is http://<ip:port>/swagger/index.html
func DefaultRouter(r *gin.Engine, jsonContent []byte) {
	registerSwagger("swagger", jsonContent)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// DefaultRouterByFile  default swagger router base on file, request url is http://<ip:port>/swagger/index.html
func DefaultRouterByFile(r *gin.Engine, jsonFile string) {
	jsonContent, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Printf("\nos.ReadFile error: %v\n\n", err)
		return
	}
	registerSwagger("swagger", jsonContent)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// CustomRouter custom swagger routing, request url is http://<ip:port>/<name>/swagger/index.html
func CustomRouter(r *gin.Engine, name string, jsonContent []byte) {
	registerSwagger(name, jsonContent)
	r.GET(fmt.Sprintf("/%s/swagger/*any", name), ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName(name)))
}

// CustomRouterByFile custom swagger router base on file, request url is http://<ip:port>/<filename prefix>/swagger/index.html
func CustomRouterByFile(r *gin.Engine, jsonFile string) {
	jsonContent, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Printf("\nos.ReadFile error: %v\n\n", err)
		return
	}

	filename := gofile.GetFilename(jsonFile)
	name := strings.Split(filename, ".")[0]
	registerSwagger(name, jsonContent)

	r.GET(fmt.Sprintf("/%s/swagger/*any", name), ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName(name)))
}

func registerSwagger(infoInstanceName string, jsonContent []byte) {
	swaggerInfo := &swag.Spec{
		Schemes:          []string{"http", "https"},
		InfoInstanceName: infoInstanceName,
		SwaggerTemplate:  string(jsonContent),
	}

	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
}
