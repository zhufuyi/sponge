package service

import (
	"bytes"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateFile generates a service.go file.
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
		funCode := fmt.Sprintf(`func (c *%sClient) %s(ctx context.Context, req *serverNameExampleV1.%s) (*serverNameExampleV1.%s, error) {
			// implement me
			// If required, fill in the code to fetch data from other rpc servers here.
			return c.%sCli.%s(ctx, req)
		}
`, field.LowerName, m.GoName, m.Input.GoIdent.GoName, m.Output.GoIdent.GoName, field.LowerName, m.GoName)
		g.P(m.Comments.Leading, funCode)
	}

	return field
}

type tmplField struct {
	Name      string // Greeter
	LowerName string // greeter first character to lower
}

func (f *tmplField) execute() string {
	buf := new(bytes.Buffer)
	if err := serviceTmpl.Execute(buf, f); err != nil {
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
