package handler

import (
	"math/rand"
	"text/template"
	"time"
)

func init() {
	var err error
	handlerLogicTmpl, err = template.New("handlerLogic").Parse(handlerLogicTmplRaw)
	if err != nil {
		panic(err)
	}
	routerTmpl, err = template.New("handlerRouter").Parse(routerTmplRaw)
	if err != nil {
		panic(err)
	}
	httpErrCodeTmpl, err = template.New("httpErrCode").Parse(httpErrCodeTmplRaw)
	if err != nil {
		panic(err)
	}

	mixLogicTmpl, err = template.New("mixLogic").Parse(mixLogicTmplRaw)
	if err != nil {
		panic(err)
	}
	mixRouterTmpl, err = template.New("mixRouter").Parse(mixRouterTmplRaw)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano()) //nolint
}

var (
	handlerLogicTmpl    *template.Template
	handlerLogicTmplRaw = `// Code generated by https://github.com/zhufuyi/sponge

package handler

import (
	"context"

	//"github.com/zhufuyi/sponge/pkg/gin/middleware"

	// import api service package here
)

{{- range .PbServices}}

var _ {{.ProtoPkgName}}.{{.Name}}Logicer = (*{{.LowerName}}Handler)(nil)

type {{.LowerName}}Handler struct {
	// example: 
	// 	{{.LowerName}}Dao dao.{{.Name}}Dao
}

// New{{.Name}}Handler create a handler
func New{{.Name}}Handler() {{.ProtoPkgName}}.{{.Name}}Logicer {
	return &{{.LowerName}}Handler{
		// example:
		// 	{{.LowerName}}Dao: dao.New{{.Name}}Dao(
		// 		model.GetDB(),
		// 		cache.New{{.Name}}Cache(model.GetCacheType()),
		// 	),
	}
}

{{- range .Methods}}

{{if eq .InvokeType 0}}{{if .Path}}{{.Comment}}
func (h *{{.LowerServiceName}}Handler) {{.MethodName}}(ctx context.Context, req *{{.RequestImportPkgName}}.{{.Request}}) (*{{.ReplyImportPkgName}}.{{.Reply}}, error) {
	panic("implement me")

	// fill in the business logic code here
	// example:
	//	    {{if .IsIgnoreShouldBind}}c, ctx := middleware.AdaptCtx(ctx)
	//	    if err = c.ShouldBindJSON(req); err != nil {
	//	    	logger.Warn("ShouldBindJSON error", logger.Error(err), middleware.CtxRequestIDField(ctx))
	//	    	return nil, ecode.InvalidParams.Err()
	//	    }{{else}}{{if .IsPassGinContext}}c, ctx := middleware.AdaptCtx(ctx){{end}}{{end}}
	//	    err := req.Validate()
	//	    if err != nil {
	//		    logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
	//		    return nil, ecode.InvalidParams.Err()
	//	    }
	//
	// 	reply, err := h.{{.LowerServiceName}}Dao.{{.MethodName}}(ctx, &model.{{.ServiceName}}{
{{- range .RequestFields}}
	//     	{{.Name}}: req.{{.Name}},
{{- end}}
	//     })
	// 	if err != nil {
	//			logger.Warn("{{.MethodName}} error", logger.Err(err), middleware.CtxRequestIDField(ctx))
	//			return nil, ecode.InternalServerError.Err()
	//		}
	//
	//     return &{{.ReplyImportPkgName}}.{{.Reply}}{
{{- range .ReplyFields}}
	//     	{{.Name}}: reply.{{.Name}},
{{- end}}
	//     }, nil
}{{end}}{{end}}

{{- end}}

// ---------- Do not delete or move this split line, this is the merge code marker ----------

{{- end}}
`

	routerTmpl    *template.Template
	routerTmplRaw = `// Code generated by https://github.com/zhufuyi/sponge

package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/zhufuyi/sponge/pkg/logger"
	//"github.com/zhufuyi/sponge/pkg/middleware"

	// import api service package here
	"moduleNameExample/internal/handler"
)

func init() {
	allMiddlewareFns = append(allMiddlewareFns, func(c *middlewareConfig) {
{{- range .PbServices}}
		{{.LowerName}}Middlewares(c)
{{- end}}
	})

	allRouteFns = append(allRouteFns,
		func(r *gin.Engine, groupPathMiddlewares map[string][]gin.HandlerFunc, singlePathMiddlewares map[string][]gin.HandlerFunc) {
{{- range .PbServices}}
			{{.LowerName}}Router(r, groupPathMiddlewares, singlePathMiddlewares, handler.New{{.Name}}Handler())
{{- end}}
		})
}

{{- range .PbServices}}

func {{.LowerName}}Router(
	r *gin.Engine,
	groupPathMiddlewares map[string][]gin.HandlerFunc,
	singlePathMiddlewares map[string][]gin.HandlerFunc,
	iService {{.ProtoPkgName}}.{{.Name}}Logicer) {
	{{.ProtoPkgName}}.Register{{.Name}}Router(
		r,
		groupPathMiddlewares,
		singlePathMiddlewares,
		iService,
		{{.ProtoPkgName}}.With{{.Name}}Logger(logger.Get()),
		{{.ProtoPkgName}}.With{{.Name}}HTTPResponse(),
		{{.ProtoPkgName}}.With{{.Name}}ErrorToHTTPCode(
			// Set some error codes to standard http return codes,
			// by default there is already ecode.InternalServerError and ecode.ServiceUnavailable
			// example:
			// 	ecode.Forbidden, ecode.LimitExceed,
		),
	)
}

// you can set the middleware of a route group, or set the middleware of a single route, 
// or you can mix them, pay attention to the duplication of middleware when mixing them, 
// it is recommended to set the middleware of a single route in preference
func {{.LowerName}}Middlewares(c *middlewareConfig) {
	// set up group route middleware, group path is left prefix rules,
	// if the left prefix is hit, the middleware will take effect, e.g. group route is /api/v1, route /api/v1/{{.LowerName}}/:id  will take effect
	// c.setGroupPath("/api/v1/{{.LowerName}}", middleware.Auth())

	// set up single route middleware, just uncomment the code and fill in the middlewares, nothing else needs to be changed
{{- range .Methods}}
	{{if eq .InvokeType 0}}{{if .Path}}//c.setSinglePath("{{.Method}}", "{{.Path}}", middleware.Auth()){{end}}{{end}}
{{- end}}
}

// ---------- Do not delete or move this split line, this is the merge code marker ----------

{{- end}}
`

	mixLogicTmpl    *template.Template
	mixLogicTmplRaw = `// Code generated by https://github.com/zhufuyi/sponge

package handler

import (
	"context"

	// import api service package here
	"moduleNameExample/internal/service"
)

{{- range .PbServices}}

var _ {{.ProtoPkgName}}.{{.Name}}Logicer = (*{{.LowerName}}Handler)(nil)

type {{.LowerName}}Handler struct {
	server {{.ProtoPkgName}}.{{.Name}}Server
}

// New{{.Name}}Handler create a handler
func New{{.Name}}Handler() {{.ProtoPkgName}}.{{.Name}}Logicer {
	return &{{.LowerName}}Handler{
		server: service.New{{.Name}}Server(),
	}
}

{{- range .Methods}}

{{if eq .InvokeType 0}}{{if .Path}}{{.Comment}}
func (h *{{.LowerServiceName}}Handler) {{.MethodName}}(ctx context.Context, req *{{.RequestImportPkgName}}.{{.Request}}) (*{{.ReplyImportPkgName}}.{{.Reply}}, error) {
	{{if eq true .IsIgnoreShouldBind .IsPassGinContext}}_, ctx = middleware.AdaptCtx(ctx){{end}}
	return h.server.{{.MethodName}}(ctx, req)
}{{end}}{{end}}

{{- end}}

// ---------- Do not delete or move this split line, this is the merge code marker ----------

{{- end}}
`

	mixRouterTmpl    *template.Template
	mixRouterTmplRaw = `// Code generated by https://github.com/zhufuyi/sponge

package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/logger"

	// import api service package here
	"moduleNameExample/internal/handler"
)

func init() {
	allMiddlewareFns = append(allMiddlewareFns, func(c *middlewareConfig) {
{{- range .PbServices}}
		{{.LowerName}}Middlewares(c)
{{- end}}
	})

	allRouteFns = append(allRouteFns,
		func(r *gin.Engine, groupPathMiddlewares map[string][]gin.HandlerFunc, singlePathMiddlewares map[string][]gin.HandlerFunc) {
{{- range .PbServices}}
			{{.LowerName}}Router(r, groupPathMiddlewares, singlePathMiddlewares, handler.New{{.Name}}Handler())
{{- end}}
		})
}

{{- range .PbServices}}

func {{.LowerName}}Router(
	r *gin.Engine,
	groupPathMiddlewares map[string][]gin.HandlerFunc,
	singlePathMiddlewares map[string][]gin.HandlerFunc,
	iService {{.ProtoPkgName}}.{{.Name}}Logicer) {
	ctxFn := func(c *gin.Context) context.Context {
		md := metadata.New(map[string]string{
			middleware.ContextRequestIDKey: middleware.GCtxRequestID(c),
		})
		return metadata.NewIncomingContext(c.Request.Context(), md)
	}
	{{.ProtoPkgName}}.Register{{.Name}}Router(
		r,
		groupPathMiddlewares,
		singlePathMiddlewares,
		iService,
		{{.ProtoPkgName}}.With{{.Name}}Logger(logger.Get()),
		{{.ProtoPkgName}}.With{{.Name}}RPCResponse(),
		{{.ProtoPkgName}}.With{{.Name}}WrapCtx(ctxFn),
		{{.ProtoPkgName}}.With{{.Name}}ErrorToHTTPCode(
			// Set some error codes to standard http return codes,
			// by default there is already ecode.InternalServerError and ecode.ServiceUnavailable
			// example:
			// 	ecode.Forbidden, ecode.LimitExceed,
		),
	)
}

// you can set the middleware of a route group, or set the middleware of a single route, 
// or you can mix them, pay attention to the duplication of middleware when mixing them, 
// it is recommended to set the middleware of a single route in preference
func {{.LowerName}}Middlewares(c *middlewareConfig) {
	// set up group route middleware, group path is left prefix rules,
	// if the left prefix is hit, the middleware will take effect, e.g. group route is /api/v1, route /api/v1/{{.LowerName}}/:id  will take effect
	// c.setGroupPath("/api/v1/{{.LowerName}}", middleware.Auth())

	// set up single route middleware, just uncomment the code and fill in the middlewares, nothing else needs to be changed
{{- range .Methods}}
	{{if eq .InvokeType 0}}{{if .Path}}//c.setSinglePath("{{.Method}}", "{{.Path}}", middleware.Auth()){{end}}{{end}}
{{- end}}
}

// ---------- Do not delete or move this split line, this is the merge code marker ----------

{{- end}}
`

	httpErrCodeTmpl *template.Template
	//nolint
	httpErrCodeTmplRaw = `// Code generated by https://github.com/zhufuyi/sponge

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

{{- range .PbServices}}

// {{.LowerName}} business-level http error codes.
// the {{.LowerName}}NO value range is 1~100, if the same error code is used, it will cause panic.
var (
	{{.LowerName}}NO       = {{.RandNumber}}
	{{.LowerName}}Name     = "{{.LowerName}}"
	{{.LowerName}}BaseCode = errcode.HCode({{.LowerName}}NO)
// --blank line--
{{- range $i, $v := .Methods}}
	{{if eq .InvokeType 0}}{{if .Path}}Err{{.MethodName}}{{.ServiceName}}   = errcode.NewError({{.LowerServiceName}}BaseCode+{{$v.AddOne $i}}, "failed to {{.MethodName}} "+{{.LowerServiceName}}Name){{end}}{{end}}
{{- end}}

	// error codes are globally unique, adding 1 to the previous error code
)

// ---------- Do not delete or move this split line, this is the merge code marker ----------

{{- end}}
`
)
