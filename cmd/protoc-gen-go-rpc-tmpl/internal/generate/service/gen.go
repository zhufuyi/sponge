// Package service is to generate template code, test code, and error code.
package service

import (
	"bytes"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/go-dev-frame/sponge/cmd/protoc-gen-go-rpc-tmpl/internal/parse"
)

// GenerateFiles generate service template code and error codes
func GenerateFiles(file *protogen.File, moduleName string) ([]byte, []byte, []byte) {
	if len(file.Services) == 0 {
		return nil, nil, nil
	}

	pss := parse.GetServices(file, moduleName)
	serviceTmplContent := genServiceTmplFile(pss)
	serviceTestTmplContent := genServiceTestTmplFile(pss)
	errCodeFileContent := genErrCodeFile(pss)

	return serviceTmplContent, serviceTestTmplContent, errCodeFileContent
}

func genServiceTmplFile(fields []*parse.PbService) []byte {
	stf := &serviceTmplFields{PbServices: fields}
	return stf.execute()
}

func genServiceTestTmplFile(pbs []*parse.PbService) []byte {
	sttf := &serviceTestTmplFields{PbServices: pbs}
	return sttf.execute()
}

func genErrCodeFile(fields []*parse.PbService) []byte {
	cf := &errCodeFields{PbServices: fields}
	return cf.execute()
}

type serviceTmplFields struct {
	PbServices []*parse.PbService
}

func (f *serviceTmplFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := serviceLogicTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	content := handleSplitLineMark(buf.Bytes())
	return bytes.ReplaceAll(content, []byte(importPkgPathMark), parse.GetImportPkg(f.PbServices))
}

type serviceTestTmplFields struct {
	PbServices []*parse.PbService
}

func (f *serviceTestTmplFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := serviceLogicTestTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	content := buf.Bytes()
	return bytes.ReplaceAll(content, []byte(importPkgPathMark), parse.GetImportPkg(f.PbServices))
}

type errCodeFields struct {
	PbServices []*parse.PbService
}

func (f *errCodeFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := rpcErrCodeTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	data := bytes.ReplaceAll(buf.Bytes(), []byte("// --blank line--"), []byte{})
	return handleSplitLineMark(data)
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
