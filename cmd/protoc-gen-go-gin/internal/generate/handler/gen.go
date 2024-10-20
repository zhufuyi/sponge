// Package handler is to generate template code, router code, and error code.
package handler

import (
	"bytes"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin/internal/parse"
)

// GenerateFiles generate handler logic, router, error code files.
func GenerateFiles(file *protogen.File, isMixType bool, moduleName string) ([]byte, []byte, []byte) {
	if len(file.Services) == 0 {
		return nil, nil, nil
	}

	pss := parse.GetServices(file, moduleName)

	var logicContent, routerFileContent, errCodeFileContent []byte

	if !isMixType {
		logicContent = genHandlerLogicFile(pss)
		routerFileContent = genRouterFile(pss)
		errCodeFileContent = genErrCodeFile(pss)
	} else {
		logicContent = genMixLogicFile(pss)
		routerFileContent = genMixRouterFile(pss)
	}

	return logicContent, routerFileContent, errCodeFileContent
}

func genHandlerLogicFile(fields []*parse.PbService) []byte {
	hlf := &handlerLogicFields{PbServices: fields}
	return hlf.execute()
}

func genRouterFile(fields []*parse.PbService) []byte {
	rf := &routerFields{PbServices: fields}
	return rf.execute()
}

func genErrCodeFile(fields []*parse.PbService) []byte {
	cf := &errCodeFields{PbServices: fields}
	return cf.execute()
}

func genMixLogicFile(fields []*parse.PbService) []byte {
	mlf := &mixLogicFields{PbServices: fields}
	return mlf.execute()
}

func genMixRouterFile(fields []*parse.PbService) []byte {
	mrf := &mixRouterFields{PbServices: fields}
	return mrf.execute()
}

type handlerLogicFields struct {
	PbServices []*parse.PbService
}

func (f *handlerLogicFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := handlerLogicTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	content := handleSplitLineMark(buf.Bytes())
	return bytes.ReplaceAll(content, []byte(importPkgPathMark), parse.GetImportPkg(f.PbServices))
}

type routerFields struct {
	PbServices []*parse.PbService
}

func (f *routerFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := routerTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	content := handleSplitLineMark(buf.Bytes())
	return bytes.ReplaceAll(content, []byte(importPkgPathMark), parse.GetSourceImportPkg(f.PbServices))
}

type errCodeFields struct {
	PbServices []*parse.PbService
}

func (f *errCodeFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := httpErrCodeTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	data := bytes.ReplaceAll(buf.Bytes(), []byte("// --blank line--"), []byte{})
	return handleSplitLineMark(data)
}

type mixLogicFields struct {
	PbServices []*parse.PbService
}

func (f *mixLogicFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := mixLogicTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	content := handleSplitLineMark(buf.Bytes())
	importPkgs := parse.GetImportPkg(f.PbServices)
	mark := []byte("ctx = middleware.AdaptCtx(ctx)")
	if bytes.Contains(content, mark) {
		importPkgs = append(importPkgs, []byte("\n\t")...)
		importPkgs = append(importPkgs, []byte(`"github.com/zhufuyi/sponge/pkg/gin/middleware"`)...)
	}
	return bytes.ReplaceAll(content, []byte(importPkgPathMark), importPkgs)
}

type mixRouterFields struct {
	PbServices []*parse.PbService
}

func (f *mixRouterFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := mixRouterTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	content := handleSplitLineMark(buf.Bytes())
	return bytes.ReplaceAll(content, []byte(importPkgPathMark), parse.GetSourceImportPkg(f.PbServices))
}

const importPkgPathMark = "// import api service package here"

var splitLineMark = []byte(`// ---------- Do not delete or move this split line, this is the merge code marker ----------`)

func handleSplitLineMark(data []byte) []byte {
	ss := bytes.Split(data, splitLineMark)
	if len(ss) <= 2 {
		return ss[0]
	}

	var out []byte
	for i, s := range ss {
		out = append(out, s...)
		if i < len(ss)-2 {
			out = append(out, splitLineMark...)
		}
	}
	return out
}
