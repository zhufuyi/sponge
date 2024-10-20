// Package router is to generate gin router code.
package router

import (
	"bytes"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin/internal/parse"
)

// GenerateFiles generate gin router code.
func GenerateFiles(file *protogen.File) []byte {
	if len(file.Services) == 0 {
		return nil
	}

	pss := parse.ParseHTTPPbServices(file)
	return genGinRouterFile(pss, string(file.GoPackageName))
}

func genGinRouterFile(services parse.HTTPPbServices, goPackageName string) []byte {
	pkg := &importPkg{
		PackageName:  goPackageName,
		PackagePaths: services.MergeImportPkgPath(),
	}
	content := pkg.execute()

	for _, service := range services {
		rf := &ginRouterFields{service}
		content = append(content, rf.execute()...)
	}
	return content
}

type ginRouterFields struct {
	*parse.HTTPPbService
}

func (f *ginRouterFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := ginRouterTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type importPkg struct {
	PackageName  string
	PackagePaths string
}

func (f *importPkg) execute() []byte {
	buf := new(bytes.Buffer)
	if err := importPkgTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
