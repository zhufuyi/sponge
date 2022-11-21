package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors cross domain
func Cors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
			ExposeHeaders:    []string{"Content-Length", "text/plain", "Authorization", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	)
}
