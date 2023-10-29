package middleware

import (
	"context"
	"net/http"

	"github.com/zhufuyi/sponge/pkg/krand"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// ContextRequestIDKey request id for context
	ContextRequestIDKey = "request_id"

	// HeaderXRequestIDKey header request id key
	HeaderXRequestIDKey = "X-Request-Id"
)

// RequestIDOption set the request id  options.
type RequestIDOption func(*requestIDOptions)

type requestIDOptions struct {
	contextRequestIDKey string
	headerXRequestIDKey string
}

func defaultRequestIDOptions() *requestIDOptions {
	return &requestIDOptions{
		contextRequestIDKey: ContextRequestIDKey,
		headerXRequestIDKey: HeaderXRequestIDKey,
	}
}

func (o *requestIDOptions) apply(opts ...RequestIDOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *requestIDOptions) setRequestIDKey() {
	if o.contextRequestIDKey != ContextRequestIDKey {
		ContextRequestIDKey = o.contextRequestIDKey
	}
	if o.headerXRequestIDKey != HeaderXRequestIDKey {
		HeaderXRequestIDKey = o.headerXRequestIDKey
	}
}

// WithContextRequestIDKey set context request id key, minimum length of 4
func WithContextRequestIDKey(key string) RequestIDOption {
	return func(o *requestIDOptions) {
		if len(key) < 4 {
			return
		}
		o.contextRequestIDKey = key
	}
}

// WithHeaderRequestIDKey set header request id key, minimum length of 4
func WithHeaderRequestIDKey(key string) RequestIDOption {
	return func(o *requestIDOptions) {
		if len(key) < 4 {
			return
		}
		o.headerXRequestIDKey = key
	}
}

// CtxKeyString for context.WithValue key type
type CtxKeyString string

// RequestIDKey request_id
var RequestIDKey = CtxKeyString(ContextRequestIDKey)

// -------------------------------------------------------------------------------------------

// RequestID is an interceptor that injects a 'request id' into the context and request/response header of each request.
func RequestID(opts ...RequestIDOption) gin.HandlerFunc {
	// customized request id key
	o := defaultRequestIDOptions()
	o.apply(opts...)
	o.setRequestIDKey()

	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(HeaderXRequestIDKey)

		// Create request id
		if requestID == "" {
			requestID = krand.String(krand.R_All, 10)
			c.Request.Header.Set(HeaderXRequestIDKey, requestID)
		}

		// Expose it for use in the application
		c.Set(ContextRequestIDKey, requestID)

		// Set X-Request-Id header
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

// HeaderRequestID get request id from the header
func HeaderRequestID(c *gin.Context) string {
	return c.Request.Header.Get(HeaderXRequestIDKey)
}

// HeaderRequestIDField get request id field from header
func HeaderRequestIDField(c *gin.Context) zap.Field {
	return zap.String(HeaderXRequestIDKey, HeaderRequestID(c))
}

// -------------------------------------------------------------------------------------------

// RequestHeaderKey request header key
var RequestHeaderKey = "request_header_key"

// WrapCtx wrap context, put the Keys and Header of gin.Context into context
func WrapCtx(c *gin.Context) context.Context {
	ctx := context.WithValue(c.Request.Context(), ContextRequestIDKey, c.GetString(ContextRequestIDKey)) //nolint
	return context.WithValue(ctx, RequestHeaderKey, c.Request.Header)                                    //nolint
}

// GetFromCtx get value from context
func GetFromCtx(ctx context.Context, key string) interface{} {
	return ctx.Value(key)
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

// GetFromHeader get value from header
func GetFromHeader(ctx context.Context, key string) string {
	header, ok := ctx.Value(RequestHeaderKey).(http.Header)
	if !ok {
		return ""
	}
	return header.Get(key)
}

// GetFromHeaders get values from header
func GetFromHeaders(ctx context.Context, key string) []string {
	header, ok := ctx.Value(RequestHeaderKey).(http.Header)
	if !ok {
		return []string{}
	}
	return header.Values(key)
}
