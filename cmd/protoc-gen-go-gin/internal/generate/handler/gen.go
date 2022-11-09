package handler

import (
	"bytes"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateFile generates handler and router files.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) (string, []byte, *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return "", nil, nil
	}

	filename := file.GeneratedFilenamePrefix + "_logic.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P(pkgImportTmplRaw)

	var fields []*tmplField
	for _, s := range file.Services {
		field := genService(g, s)
		fields = append(fields, field)
	}

	rf := &routerFields{ServiceNames: fields}

	return filename, rf.execute(), g
}

func genService(g *protogen.GeneratedFile, s *protogen.Service) *tmplField {
	field := &tmplField{
		Name:      s.GoName,
		LowerName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],
	}
	g.P(field.execute())

	for _, m := range s.Methods {
		funcCode := fmt.Sprintf(`func (h *%sHandler) %s(ctx context.Context, req *serverNameExampleV1.%s) (*serverNameExampleV1.%s, error) {
			// example:
			// 	reply, err := h.%sDao.%s(ctx, req)
			// 	if err != nil {
			//			logger.Warn("invoke error", logger.Err(err), middleware.CtxRequestIDField(ctx))
			//			return nil, ecode.InternalServerError.Err()
			//		}
			// 	return reply, nil
			//
			// If required, fill in the code for getting data from other dao here

			panic("implement me")
		}
`, field.LowerName, m.GoName, m.Input.GoIdent.GoName, m.Output.GoIdent.GoName, field.LowerName, m.GoName)
		g.P(m.Comments.Leading, funcCode)
	}

	return field
}

type tmplField struct {
	Name      string // Greeter
	LowerName string // greeter first character to lower
}

func (f *tmplField) execute() string {
	buf := new(bytes.Buffer)
	if err := handlerTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return buf.String()
}

type routerFields struct {
	ServiceNames []*tmplField
}

func (f *routerFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := routerTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
