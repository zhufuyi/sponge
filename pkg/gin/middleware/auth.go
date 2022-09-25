package middleware

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Auth 鉴权
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if len(authorization) < 20 {
			logger.Warn("authorization is illegal", logger.String("authorization", authorization))
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		token := authorization[7:] // 去掉Bearer 前缀
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

// AuthAdmin 管理员鉴权
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if len(authorization) < 20 {
			logger.Warn("authorization is illegal", logger.String("authorization", authorization))
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		token := authorization[7:] // 去掉Bearer 前缀
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			logger.Warn("VerifyToken error", logger.Err(err))
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}

		// 判断是否为管理员
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
