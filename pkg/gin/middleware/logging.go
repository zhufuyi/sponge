package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var contentMark = []byte(" ...... ")

var (
	// Print body max length
	defaultMaxLength = 300

	// default zap log
	defaultLogger, _ = zap.NewProduction()

	// Ignore route list
	defaultIgnoreRoutes = map[string]struct{}{
		"/ping":   {},
		"/pong":   {},
		"/health": {},
	}
)

// Option set the gin logger options.
type Option func(*options)

func defaultOptions() *options {
	return &options{
		maxLength:     defaultMaxLength,
		log:           defaultLogger,
		ignoreRoutes:  defaultIgnoreRoutes,
		requestIDFrom: 0,
	}
}

type options struct {
	maxLength     int
	log           *zap.Logger
	ignoreRoutes  map[string]struct{}
	requestIDFrom int // 0: ignore, 1: from context, 2: from header
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithMaxLen logger content max length
func WithMaxLen(maxLen int) Option {
	return func(o *options) {
		o.maxLength = maxLen
	}
}

// WithLog set log
func WithLog(log *zap.Logger) Option {
	return func(o *options) {
		if log != nil {
			o.log = log
		}
	}
}

// WithIgnoreRoutes no logger content routes
func WithIgnoreRoutes(routes ...string) Option {
	return func(o *options) {
		for _, route := range routes {
			o.ignoreRoutes[route] = struct{}{}
		}
	}
}

// WithRequestIDFromContext name is field in context, default value is request_id
func WithRequestIDFromContext() Option {
	return func(o *options) {
		o.requestIDFrom = 1
	}
}

// WithRequestIDFromHeader name is field in header, default value is X-Request-Id
func WithRequestIDFromHeader() Option {
	return func(o *options) {
		o.requestIDFrom = 2
	}
}

// ------------------------------------------------------------------------------------------

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// If there is sensitive information in the body, you can use WithIgnoreRoutes set the route to ignore logging
func getBodyData(buf *bytes.Buffer, maxLen int) []byte {
	l := buf.Len()
	if l == 0 {
		return []byte("")
	} else if l <= maxLen {
		return buf.Bytes()[:l-1]
	}
	return append(bytes.Clone(buf.Bytes()[:maxLen]), contentMark...)
}

// Logging print request and response info
func Logging(opts ...Option) gin.HandlerFunc {
	o := defaultOptions()
	o.apply(opts...)

	return func(c *gin.Context) {
		start := time.Now()

		// ignore printing of the specified route
		if _, ok := o.ignoreRoutes[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		// print input information before processing
		buf := bytes.Buffer{}
		_, _ = buf.ReadFrom(c.Request.Body)

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
		}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch || c.Request.Method == http.MethodDelete {
			fields = append(fields,
				zap.Int("size", buf.Len()),
				zap.ByteString("body", getBodyData(&buf, o.maxLength)),
			)
		}

		reqID := ""
		if o.requestIDFrom == 1 {
			if v, isExist := c.Get(ContextRequestIDKey); isExist {
				if requestID, ok := v.(string); ok {
					reqID = requestID
					fields = append(fields, zap.String(ContextRequestIDKey, reqID))
				}
			}
		} else if o.requestIDFrom == 2 {
			reqID = c.Request.Header.Get(HeaderXRequestIDKey)
			fields = append(fields, zap.String(ContextRequestIDKey, reqID))
		}
		o.log.Info("<<<<", fields...)

		c.Request.Body = io.NopCloser(&buf)

		// replace writer
		newWriter := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = newWriter

		// processing requests
		c.Next()

		// print return message after processing
		fields = []zap.Field{
			zap.Int("code", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.Path),
			zap.Int64("time_us", time.Since(start).Microseconds()),
			zap.Int("size", newWriter.body.Len()),
			zap.ByteString("body", getBodyData(newWriter.body, o.maxLength)),
		}
		if reqID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, reqID))
		}
		o.log.Info(">>>>", fields...)
	}
}

// SimpleLog print response info
func SimpleLog(opts ...Option) gin.HandlerFunc {
	o := defaultOptions()
	o.apply(opts...)

	return func(c *gin.Context) {
		start := time.Now()

		// ignore printing of the specified route
		if _, ok := o.ignoreRoutes[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		reqID := ""
		if o.requestIDFrom == 1 {
			if v, isExist := c.Get(ContextRequestIDKey); isExist {
				if requestID, ok := v.(string); ok {
					reqID = requestID
				}
			}
		} else if o.requestIDFrom == 2 {
			reqID = c.Request.Header.Get(HeaderXRequestIDKey)
		}

		// processing requests
		c.Next()

		// print return message after processing
		fields := []zap.Field{
			zap.Int("code", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.Int64("time_us", time.Since(start).Microseconds()),
			zap.Int("size", c.Writer.Size()),
		}
		if reqID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, reqID))
		}
		o.log.Info("[GIN]", fields...)
	}
}
