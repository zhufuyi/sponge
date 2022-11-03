package service

import (
	"bytes"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateFile generates a _service.pb.go file.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + "_service.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	for _, s := range file.Services {
		genService(g, s)
	}
	return g
}

func genService(g *protogen.GeneratedFile, s *protogen.Service) {
	field := &tmplField{
		Name:      s.GoName,
		LowerName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],
	}
	g.P(field.execute())

	for _, m := range s.Methods {
		funCode := fmt.Sprintf(`func (c *%sClient) %s(ctx context.Context, req *serverNameExampleV1.%s) (*serverNameExampleV1.%s, error) {
			// implement me
			// If required, fill in the code to fetch data from other microservices here.
			return c.%sCli.%s(ctx, req)
		}
`, field.LowerName, m.GoName, m.Input.GoIdent.GoName, m.Output.GoIdent.GoName, field.LowerName, m.GoName)
		g.P(m.Comments.Leading, funCode)
	}
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
