package handlerFile

import (
	"text/template"
)

func init() {
	var err error
	handlerTmpl, err = template.New("iHandler").Parse(handlerTmplRaw)
	if err != nil {
		panic(err)
	}
}

var (
	handlerTmpl    *template.Template
	handlerTmplRaw = `package handler

import (
	"context"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
)

var _ serverNameExampleV1.{{$.Name}}Logicer = (*{{$.LowerName}}Handler)(nil)

type {{$.LowerName}}Handler struct {
	// example: 
	// 	{{$.LowerName}}Dao dao.{{$.Name}}Dao

	// If required, fill in the definition of the other dao code here.
}

// New{{$.Name}}Handler creating handler
func New{{$.Name}}Handler() serverNameExampleV1.{{$.Name}}Logicer {
	return &{{$.LowerName}}Handler{
		// example:
		// {{$.LowerName}}Dao: dao.New{{$.Name}}Dao(
		// 	model.GetDB(),
		// 	cache.New{{$.Name}}Cache(model.GetCacheType()),
		// ),

		// If required, fill in the code to implement other dao here.
	}
}
`
)
