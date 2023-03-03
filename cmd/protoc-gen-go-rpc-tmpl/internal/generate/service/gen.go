package service

import (
	"bytes"
	"strings"

	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl/internal/parse"
	"github.com/zhufuyi/sponge/pkg/gofile"

	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateFiles generate service template code and error codes
func GenerateFiles(filenamePrefix string, file *protogen.File) ([]byte, []byte, []byte) {
	if len(file.Services) == 0 {
		return nil, nil, nil
	}

	protoName := getProtoFilename(filenamePrefix)
	pss := parse.GetServices(protoName, file)
	serviceTmplContent := genServiceTmplFile(pss)
	serviceTestTmplContent := genServiceTestTmplFile(pss)
	errCodeFileContent := genErrCodeFile(pss)

	return serviceTmplContent, serviceTestTmplContent, errCodeFileContent
}

func genServiceTmplFile(fields []*parse.PbService) []byte {
	lf := &serviceTmplFields{PbServices: fields}
	return lf.execute()
}

func genServiceTestTmplFile(pbs []*parse.PbService) []byte {
	lf := &serviceTestTmplFields{PbServices: pbs}
	return lf.execute()
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
	return buf.Bytes()
}

type serviceTestTmplFields struct {
	PbServices []*parse.PbService
}

func (f *serviceTestTmplFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := serviceLogicTestTmpl.Execute(buf, f); err != nil {
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

func getProtoFilename(filenamePrefix string) string {
	filenamePrefix = strings.ReplaceAll(filenamePrefix, ".proto", "")
	filenamePrefix = strings.ReplaceAll(filenamePrefix, gofile.GetPathDelimiter(), "/")
	ss := strings.Split(filenamePrefix, "/")

	if len(ss) == 0 {
		return ""
	} else if len(ss) == 1 {
		return ss[0] + ".proto"
	}

	return ss[len(ss)-1] + ".proto"
}
