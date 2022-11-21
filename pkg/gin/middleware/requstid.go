package middleware

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/krand"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// ContextRequestIDKey context request id for context
	ContextRequestIDKey = "request_id"

	// HeaderXRequestIDKey http header request id key
	HeaderXRequestIDKey = "X-Request-ID"
)

// RequestID is an interceptor that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(HeaderXRequestIDKey)

		// Create request id
		if requestID == "" {
			requestID = krand.String(krand.R_All, 12) // generate a random string of length 12
			c.Request.Header.Set(HeaderXRequestIDKey, requestID)
			// Expose it for use in the application
			c.Set(ContextRequestIDKey, requestID)
		}

		// Set X-Request-ID header
		c.Writer.Header().Set(HeaderXRequestIDKey, requestID)

		c.Next()
	}
}

// GCtxRequestID get request id from gin.Context
func GCtxRequestID(c *gin.Context) string {
	if v, isExist := c.Get(ContextRequestIDKey); isExist {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}
	return ""
}

// GCtxRequestIDField get request id field from gin.Context
func GCtxRequestIDField(c *gin.Context) zap.Field {
	return zap.String(ContextRequestIDKey, GCtxRequestID(c))
}

// CtxRequestID get request id from context.Context
func CtxRequestID(ctx context.Context) string {
	v := ctx.Value(ContextRequestIDKey)
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

// CtxRequestIDField get request id field from context.Context
func CtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, CtxRequestID(ctx))
}

// HeaderRequestID get request id from the header
func HeaderRequestID(c *gin.Context) string {
	return c.Request.Header.Get(HeaderXRequestIDKey)
}

// HeaderRequestIDField get request id field from header
func HeaderRequestIDField(c *gin.Context) zap.Field {
	return zap.String(HeaderXRequestIDKey, HeaderRequestID(c))
}
