package handlerfunc

import (
	"net/http"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
)

// checkHealthResponse check health result
type checkHealthResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// CheckHealth check healthy.
// @Summary check health
// @Description check health
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} checkHealthResponse{}
// @Router /health [get]
func CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, checkHealthResponse{Status: "UP", Hostname: utils.GetHostname()})
}

// Ping ping
// @Summary ping
// @Description ping
// @Tags system
// @Accept  json
// @Produce  json
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
