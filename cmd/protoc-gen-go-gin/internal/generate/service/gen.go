package service

import (
	"bytes"

	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin/internal/parse"

	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateFiles generate service logic, router, error code files.
func GenerateFiles(file *protogen.File) ([]byte, []byte, []byte) {
	if len(file.Services) == 0 {
		return nil, nil, nil
	}

	pss := parse.GetServices(file)
	logicContent := genServiceLogicFile(pss)
	routerFileContent := genRouterFile(pss)
	errCodeFileContent := genErrCodeFile(pss)

	return logicContent, routerFileContent, errCodeFileContent
}

func genServiceLogicFile(fields []*parse.PbService) []byte {
	lf := &serviceLogicFields{PbServices: fields}
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

type serviceLogicFields struct {
	PbServices []*parse.PbService
}

func (f *serviceLogicFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := serviceLogicTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type routerFields struct {
	PbServices []*parse.PbService
}

func (f *routerFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := routerTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type errCodeFields struct {
	PbServices []*parse.PbService
}

func (f *errCodeFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := rpcErrCodeTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return bytes.ReplaceAll(buf.Bytes(), []byte("// --blank line--"), []byte{})
}
