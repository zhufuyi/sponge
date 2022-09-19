package middleware

import (
	"github.com/zhufuyi/sponge/pkg/krand"

	"github.com/gin-gonic/gin"
)

const (
	// ContextRequestIDKey context request id for context
	ContextRequestIDKey = "request_id"

	// HeaderXRequestIDKey http header request ID key
	HeaderXRequestIDKey = "X-Request-ID"
)

// RequestID is an interceptor that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(HeaderXRequestIDKey)

		// Create request id
		if requestID == "" {
			requestID = krand.String(krand.R_All, 12) // 生成长度为12的随机字符串
			c.Request.Header.Set(HeaderXRequestIDKey, requestID)
			// Expose it for use in the application
			c.Set(ContextRequestIDKey, requestID)
		}

		// Set X-Request-ID header
		c.Writer.Header().Set(HeaderXRequestIDKey, requestID)

		c.Next()
	}
}

// GetRequestIDFromContext returns 'RequestID' from the given context if present.
func GetRequestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(ContextRequestIDKey); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}

// GetRequestIDFromHeaders returns 'RequestID' from the headers if present.
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(HeaderXRequestIDKey)
}
