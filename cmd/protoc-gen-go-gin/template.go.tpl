type {{ $.InterfaceName }} interface {
{{range .MethodSet}}{{.Name}}(ctx context.Context, req *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

func Register{{$.Name}}HTTPServer(r gin.IRouter, iService {{ $.InterfaceName }}, resps ...errcode.Responser) {
	var resp = errcode.NewResponse()
	if len(resps) > 0 {
		resp = resps[0] // replace the default response methods
	}

	s := &{{$.Name}}{
		router: r,
		iService: iService,
		resp:   resp,
	}
	s.RegisterService()
}

type {{$.Name}} struct {
	router gin.IRouter
	iService {{ $.InterfaceName }}
	resp   errcode.Responser
}

func (s *{{$.Name}}) RegisterService() {
{{range .Methods}}s.router.Handle("{{.Method}}", "{{.Path}}", s.{{ .HandlerName }})
{{end}}
}

{{range .Methods}}
func (s *{{$.Name}}) {{ .HandlerName }} (c *gin.Context) {
	req := &{{.Request}}{}
{{if .HasPathParams }}
	if err := c.ShouldBindUri(req); err != nil {
		logger.Warn("ShouldBindUri error", logger.Err(err))
		s.resp.ParamError(c, err)
		return
	}
{{end}}
{{if eq .Method "GET" "DELETE" }}
	if err := c.ShouldBindQuery(req); err != nil {
		logger.Warn("ShouldBindQuery error", logger.Err(err))
		s.resp.ParamError(c, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Warn("ShouldBindJSON error", logger.Err(err))
		s.resp.ParamError(c, err)
		return
	}
{{else}}
	if err := c.ShouldBind(req); err != nil {
		logger.Warn("ShouldBind error", logger.Err(err))
		s.resp.ParamError(c, err)
		return
	}
{{end}}
	md := metadata.New(nil)
	for k, v := range c.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(c, md)
	out, err := s.iService.{{.Name}}(newCtx, req)
	if err != nil {
		logger.Error("{{.Name}} error", logger.Err(err))
		s.resp.Error(c, err)
		return
	}

	s.resp.Success(c, out)
}
{{end}}
