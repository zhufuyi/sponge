package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// Default request id name
	defaultRequestIDNameInHeader  = "X-Request-Id"
	defaultRequestIDNameInContext = "request_id"
)

// GetRequestIDFromContext returns 'RequestID' from the given context if present.
func GetRequestIDFromContext(c *gin.Context) string {
	if v, isExist := c.Get(defaultRequestIDNameInContext); isExist {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}
	return ""
}

// GetRequestIDFromHeaders returns 'RequestID' from the headers if present.
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(defaultRequestIDNameInHeader)
}

// FieldRequestIDFromContext zap logger request ID from context
func FieldRequestIDFromContext(c *gin.Context, name ...string) zap.Field {
	var requestIDName string
	if len(name) > 0 && name[0] != "" {
		requestIDName = name[0]
	} else {
		requestIDName = defaultRequestIDNameInContext
	}
	return zap.String(requestIDName, GetRequestIDFromContext(c))
}

// FieldRequestIDFromHeader zap logger request ID from header
func FieldRequestIDFromHeader(c *gin.Context, name ...string) zap.Field {
	var requestIDName string
	if len(name) > 0 && name[0] != "" {
		requestIDName = name[0]
	} else {
		requestIDName = defaultRequestIDNameInHeader
	}
	return zap.String(requestIDName, GetRequestIDFromHeaders(c))
}
