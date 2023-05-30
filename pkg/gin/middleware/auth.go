package middleware

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	// HeaderAuthorizationKey http header authorization key
	HeaderAuthorizationKey = "Authorization"
)

// Auth authorization
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader(HeaderAuthorizationKey)
		if len(authorization) < 20 {
			logger.Warn("authorization is illegal", logger.String(HeaderAuthorizationKey, authorization))
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		token := authorization[7:] // remove Bearer prefix
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			logger.Warn("VerifyToken error", logger.Err(err))
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		c.Set("uid", claims.UID)

		c.Next()
	}
}

// AuthAdmin admin authentication
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader(HeaderAuthorizationKey)
		if len(authorization) < 20 {
			logger.Warn("authorization is illegal", logger.String(HeaderAuthorizationKey, authorization))
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		token := authorization[7:] // remove Bearer prefix
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			logger.Warn("VerifyToken error", logger.Err(err))
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}

		// determine if it is an administrator
		if claims.Role != "admin" {
			logger.Warn("prohibition of access", logger.String("uid", claims.UID), logger.String("role", claims.Role))
			response.Error(c, errcode.Forbidden)
			c.Abort()
			return
		}
		c.Set("uid", claims.UID)

		c.Next()
	}
}
