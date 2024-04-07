// Package handler is to generate template code, router code, and error code.
package handler

import (
	"bytes"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin/internal/parse"
)

// GenerateFiles generate handler logic, router, error code files.
func GenerateFiles(file *protogen.File, isMixType bool) ([]byte, []byte, []byte) {
	if len(file.Services) == 0 {
		return nil, nil, nil
	}

	pss := parse.GetServices(file)

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
	lf := &handlerLogicFields{PbServices: fields}
	return lf.execute()
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
	return handleSplitLineMark(buf.Bytes())
}

type routerFields struct {
	PbServices []*parse.PbService
}

func (f *routerFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := routerTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return handleSplitLineMark(buf.Bytes())
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
	return handleSplitLineMark(buf.Bytes())
}

type mixRouterFields struct {
	PbServices []*parse.PbService
}

func (f *mixRouterFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := mixRouterTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return handleSplitLineMark(buf.Bytes())
}

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
