package router

import (
	"text/template"
)

func init() {
	var err error
	handlerTmpl, err = template.New("iRouter").Parse(handlerTmplRaw)
	if err != nil {
		panic(err)
	}
}

var (
	handlerTmpl    *template.Template
	handlerTmplRaw = `
type {{$.Name}}Logicer interface {
{{range .MethodSet}}{{.Name}}(ctx context.Context, req *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

type {{$.Name}}Option func(*{{$.LowerName}}Options)

type {{$.LowerName}}Options struct {
	isFromRPC bool
	responser errcode.Responser
	zapLog    *zap.Logger
}

func (o *{{$.LowerName}}Options) apply(opts ...{{$.Name}}Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func With{{$.Name}}HTTPResponse() {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.isFromRPC = false
	}
}

func With{{$.Name}}RPCResponse() {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.isFromRPC = true
	}
}

func With{{$.Name}}Responser(responser errcode.Responser) {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.responser = responser
	}
}

func With{{$.Name}}Logger(zapLog *zap.Logger) {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.zapLog = zapLog
	}
}

func Register{{$.Name}}Router(iRouter gin.IRouter, iLogic {{$.Name}}Logicer, opts ...{{$.Name}}Option) {
	o := &{{$.LowerName}}Options{}
	o.apply(opts...)

	if o.responser == nil {
		o.responser = errcode.NewResponse(o.isFromRPC)
	}
	if o.zapLog == nil {
		o.zapLog,_ = zap.NewProduction()
	}

	r := &{{$.LowerName}}Router {
		iRouter:   iRouter,
		iLogic:    iLogic,
		iResponse: o.responser,
		zapLog:    o.zapLog,
	}
	r.register()
}

type {{$.LowerName}}Router struct {
	iRouter   gin.IRouter
	iLogic    {{$.Name}}Logicer
	iResponse errcode.Responser
	zapLog    *zap.Logger
}

func (r *{{$.LowerName}}Router) register() {
{{range .Methods}}r.iRouter.Handle("{{.Method}}", "{{.Path}}", r.{{ .HandlerName }})
{{end}}
}

{{range .Methods}}
func (r *{{$.LowerName}}Router) {{ .HandlerName }} (c *gin.Context) {
	req := &{{.Request}}{}
{{if .HasPathParams }}
	if err := c.ShouldBindUri(req); err != nil {
		r.zapLog.Warn("ShouldBindUri error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{end}}
{{if eq .Method "GET" "DELETE" }}
	if err := c.ShouldBindQuery(req); err != nil {
		r.zapLog.Warn("ShouldBindQuery error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
	if err := c.ShouldBindJSON(req); err != nil {
		r.zapLog.Warn("ShouldBindJSON error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{else}}
	if err := c.ShouldBind(req); err != nil {
		r.zapLog.Warn("ShouldBind error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{end}}
	out, err := r.iLogic.{{.Name}}(c.Request.Context(), req)
	if err != nil {
		isIgnore := r.iResponse.Error(c, err)
		if !isIgnore {
			r.zapLog.Error("{{.Name}} error", zap.Error(err), middleware.GCtxRequestIDField(c))
		}
		return
	}

	r.iResponse.Success(c, out)
}
{{end}}
`
)
