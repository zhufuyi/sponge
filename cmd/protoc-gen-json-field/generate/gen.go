// Package generate is generate json field code.
package generate

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/go-dev-frame/sponge/cmd/protoc-gen-json-field/parser"
)

// ProtoInfo is the info for parsing proto file
type ProtoInfo struct {
	FileName            string    // proto file name, example: foo_bar.proto
	FileNamePrefix      string    // proto file name, not include suffix, example: foo_bar
	FileNamePrefixCamel string    // proto file name, not include suffix, example: FooBar
	Services            []Service // services in proto file

	FileDir string // proto file directory, example: api/user/v1
	Package string // proto package name, example: api.user.v1

	// go related fields
	GoPkgName         string            // go package name, example: userV1
	GoPackage         string            // go package name, example: "moduleName/api/user/v1"
	ImportPkgMap      map[string]string // rpc params import packages, example: userV1 -> userV1 "moduleName/api/user/v1"
	FieldImportPkgMap map[string]string // message import packages, example: userV1 -> userV1 "moduleName/api/user/v1"
}

type Service struct {
	ServiceName               string // service name, example: foobar or Foobar
	ServiceNameCamel          string // service name camel case, example: FooBar
	ServiceNameCamelFCL       string // service name camel case first character to lower, example: fooBar
	ServiceNamePluralCamel    string // service name plural, camel case, example: FooBars
	ServiceNamePluralCamelFCL string // service name plural, camel case and first character lower, example: fooBars

	GoPkgName string // go package name, example: userV1

	Methods []RPCMethod // rpc methods
}

type RPCMethod struct {
	MethodName string // method name, example: Create
	Comment    string // method comment, example: // Create a record
	InvokeType string // rpc invoke type: unary_call, client_side_streaming, server_side_streaming, bidirectional_streaming

	RequestName          string  // method request message, example: CreateRequest
	RequestFields        []Field // request fields
	RequestImportPkgName string  // request import package name, example: emptypb

	ReplyName          string  // method reply message, example: CreateReply
	ReplyFields        []Field // reply fields
	ReplyImportPkgName string  // reply import package name, example: emptypb

	// google.api.http options fields
	HTTPRouter        string // http router, example: /api/user/v1/create
	HTTPRequestMethod string // http request method, example: POST
	HTTPRequestBody   string // http request body, example: CreateRequest

	// google.api.http selector custom options, only for gin
	IsPassGinContext bool
	IsIgnoreGinBind  bool
}

type Field struct {
	Name           string // field name
	GoType         string // field go type
	GoTypeCrossPkg string // field go type cross package
	Comment        string // field comment
	FieldType      string // field type, if field type is message, it will be used as import package name
	ImportPkgName  string // import package name, example: anypb
	ImportPkgPath  string // import path e.g. google.golang.org/protobuf/types/known/anypb
}

// GenerateFiles generate service logic, router, error code files.
func GenerateFiles(file *protogen.File) ([]byte, error) {
	pss := parser.GetServices(file)
	goPackage := file.GoDescriptorIdent.GoImportPath.String()
	goPkgName := parser.GetProtoPkgName(goPackage)

	v := newProtoInfo(pss, goPkgName)

	v.Package = string(file.Desc.Package())
	v.FileDir = strings.ReplaceAll(v.Package, ".", "/")
	v.GoPackage = goPackage
	v.GoPkgName = goPkgName
	_, v.FileNamePrefix = filepath.Split(file.GeneratedFilenamePrefix)
	v.FileNamePrefixCamel = xstrings.ToCamelCase(v.FileNamePrefix)
	v.FileName = v.FileNamePrefix + ".proto"

	return json.MarshalIndent(v, "", "  ")
}

func newProtoInfo(pss []*parser.PbService, goPkgName string) *ProtoInfo {
	var (
		services          = []Service{}
		importPkgMap      = make(map[string]string)
		fieldImportPkgMap = make(map[string]string)
	)
	for _, ps := range pss {
		var methods []RPCMethod
		for _, m := range ps.Methods {
			methods = append(methods, RPCMethod{
				MethodName:           m.MethodName,
				Comment:              m.Comment,
				InvokeType:           getInvokeType(m.InvokeType),
				RequestName:          m.Request,
				RequestFields:        getMessageFields(m.RequestFields),
				RequestImportPkgName: m.RequestImportPkgName,
				ReplyName:            m.Reply,
				ReplyFields:          getMessageFields(m.ReplyFields),
				ReplyImportPkgName:   m.ReplyImportPkgName,
				HTTPRouter:           m.Path,
				HTTPRequestMethod:    m.Method,
				HTTPRequestBody:      m.Body,
				IsPassGinContext:     m.IsPassGinContext,
				IsIgnoreGinBind:      m.IsIgnoreShouldBind,
			})
		}
		importPkgMap = mergeMap(importPkgMap, ps.ImportPkgMap)
		fieldImportPkgMap = mergeMap(fieldImportPkgMap, ps.FieldImportPkgMap)
		serviceNameCamel := xstrings.ToCamelCase(ps.Name)
		pluralName := inflection.Plural(ps.Name)

		services = append(services, Service{
			ServiceName:               ps.Name,
			ServiceNameCamel:          serviceNameCamel,
			ServiceNameCamelFCL:       firstLetterToLower(serviceNameCamel),
			ServiceNamePluralCamel:    customEndOfLetterToLower(serviceNameCamel, pluralName),
			ServiceNamePluralCamelFCL: firstLetterToLower(customEndOfLetterToLower(serviceNameCamel, pluralName)),
			GoPkgName:                 goPkgName,
			Methods:                   methods,
		})
	}

	return &ProtoInfo{
		Services:          services,
		ImportPkgMap:      importPkgMap,
		FieldImportPkgMap: fieldImportPkgMap,
	}
}

func firstLetterToLower(str string) string {
	if len(str) == 0 {
		return str
	}

	if (str[0] >= 'A' && str[0] <= 'Z') || (str[0] >= 'a' && str[0] <= 'z') {
		return strings.ToLower(str[:1]) + str[1:]
	}

	return str
}

func customEndOfLetterToLower(srcStr string, str string) string {
	l := len(str) - len(srcStr)
	if l == 1 {
		if str[len(str)-1] == 'S' {
			return str[:len(str)-1] + "s"
		}
	} else if l == 2 {
		if str[len(str)-2:] == "ES" {
			return str[:len(str)-2] + "es"
		}
	}

	return str
}

func getInvokeType(t int) string {
	switch t {
	case 0:
		return "unary_call"
	case 1:
		return "client_side_streaming"
	case 2:
		return "server_side_streaming"
	case 3:
		return "bidirectional_streaming"
	default:
		return ""
	}
}

func getMessageFields(fs []*parser.Field) []Field {
	fields := make([]Field, 0, len(fs))
	for _, f := range fs {
		fields = append(fields, Field{
			Name:           f.Name,
			GoType:         f.GoType,
			GoTypeCrossPkg: f.GoTypeCrossPkg,
			Comment:        f.Comment,
			FieldType:      f.FieldType,
			ImportPkgName:  f.ImportPkgName,
			ImportPkgPath:  f.ImportPkgPath,
		})
	}
	return fields
}

func mergeMap(m1, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
