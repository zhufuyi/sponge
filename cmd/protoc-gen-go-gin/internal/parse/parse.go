// Package parse is parsed proto file to struct
package parse

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// Field message field
type Field struct {
	Name      string // field name
	FieldType string // field type
	Comment   string // field comment
}

// GoTypeZero default zero value for type
func (r Field) GoTypeZero() string {
	switch r.FieldType {
	case "bool":
		return "false"
	case "int32", "uint32", "sint32", "int64", "uint64", "sint64", "sfixed32", "fixed32", "sfixed64", "fixed64":
		return "0"
	case "float", "double":
		return "0.0"
	case "string":
		return `""`
	default:
		return "nil"
	}
}

// ServiceMethod RPCMethod fields
type ServiceMethod struct {
	MethodName    string // Create
	Request       string // CreateRequest
	RequestFields []*Field
	Reply         string // CreateReply
	ReplyFields   []*Field
	Comment       string // e.g. Create a record
	InvokeType    int    // 0:unary, 1: client-side streaming, 2: server-side streaming, 3: bidirectional streaming

	ServiceName         string // Greeter
	LowerServiceName    string // greeter first character to lower
	LowerCutServiceName string // GreeterService --> greeter

	// http_rule
	Path   string // rule
	Method string // HTTP Method
	Body   string

	IsPassGinContext   bool
	IsIgnoreShouldBind bool

	RequestImportPkgName string // e.g. userV1
	ReplyImportPkgName   string // e.g. userV1
	ProtoPkgName         string // e.g. userV1
}

// AddOne counter
func (t *ServiceMethod) AddOne(i int) int {
	return i + 1
}

// PbService service fields
type PbService struct {
	Name      string           // Greeter
	LowerName string           // greeter first character to lower
	Methods   []*ServiceMethod // service methods

	CutServiceName      string // GreeterService --> Greeter
	LowerCutServiceName string // GreeterService --> greeter

	ImportPkgMap map[string]string // e.g. [userV1]:[userV1 "user/api/user/v1"]

	ProtoFileDir string // e.g. api/user/v1
	ProtoPkgName string // e.g. userV1
	ModuleName   string
}

// RandNumber rand number 1~100
func (s *PbService) RandNumber() int {
	return rand.Intn(99) + 1
}

func parsePbService(s *protogen.Service, protoFileDir string, moduleName string) *PbService {
	protoPkgName := convertToPkgName(protoFileDir)
	cutServiceName := getCutServiceName(s.GoName)
	importPkgMap := map[string]string{}

	var methods []*ServiceMethod
	for _, m := range s.Methods {
		rpcMethod := &RPCMethod{} //nolint
		rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule != nil && ok {
			rpcMethod = buildHTTPRule(m, rule, protoPkgName)
		} /*else {
			// if the http method and path is not set, set default value.
			//rpcMethod = defaultMethod(m)
		}*/

		requestImportPkgName := convertToPkgName(m.Input.GoIdent.GoImportPath.String())
		replyImportPkgName := convertToPkgName(m.Output.GoIdent.GoImportPath.String())
		if requestImportPkgName != "" {
			importPkgMap[requestImportPkgName] = requestImportPkgName + " " + m.Input.GoIdent.GoImportPath.String()
		}
		if replyImportPkgName != "" {
			importPkgMap[replyImportPkgName] = replyImportPkgName + " " + m.Output.GoIdent.GoImportPath.String()
		}

		methods = append(methods, &ServiceMethod{
			MethodName:    m.GoName,
			Request:       m.Input.GoIdent.GoName,
			RequestFields: getFields(m.Input),
			Reply:         m.Output.GoIdent.GoName,
			ReplyFields:   getFields(m.Output),
			Comment:       getMethodComment(m),
			InvokeType:    getInvokeType(m.Desc.IsStreamingClient(), m.Desc.IsStreamingServer()),

			ServiceName:         s.GoName,
			LowerServiceName:    strings.ToLower(s.GoName[:1]) + s.GoName[1:],
			LowerCutServiceName: strings.ToLower(cutServiceName[:1]) + cutServiceName[1:],

			Path:   rpcMethod.Path,
			Method: rpcMethod.Method,
			Body:   rpcMethod.Body,

			IsPassGinContext:   rpcMethod.IsPassGinContext,
			IsIgnoreShouldBind: rpcMethod.IsIgnoreShouldBind,

			RequestImportPkgName: requestImportPkgName,
			ReplyImportPkgName:   replyImportPkgName,
			ProtoPkgName:         protoPkgName,
		})
	}

	return &PbService{
		Name:                s.GoName,
		LowerName:           strings.ToLower(s.GoName[:1]) + s.GoName[1:],
		Methods:             methods,
		CutServiceName:      cutServiceName,
		LowerCutServiceName: strings.ToLower(cutServiceName[:1]) + cutServiceName[1:],
		ImportPkgMap:        importPkgMap,
		ProtoFileDir:        protoFileDir,
		ProtoPkgName:        protoPkgName,
		ModuleName:          moduleName,
	}
}

// GetServices parse protobuf services
func GetServices(file *protogen.File, moduleName string) []*PbService {
	protoFileDir := getProtoFileDir(file.GeneratedFilenamePrefix)
	var pss []*PbService
	for _, s := range file.Services {
		pss = append(pss, parsePbService(s, protoFileDir, moduleName))
	}
	return pss
}

func getCutServiceName(name string) string {
	service := "Service"
	if len(name) < len(service) {
		return name
	}
	l := len(name) - len(service)
	if name[l:] == service {
		if name[:l] == "" {
			return name
		}
		return name[:l]
	}
	return name
}

func getFields(m *protogen.Message) []*Field {
	var fields []*Field
	for _, f := range m.Fields {
		fieldType := f.Desc.Kind().String()
		if f.Desc.Cardinality().String() == "repeated" {
			fieldType = "[]" + fieldType
		}
		fields = append(fields, &Field{
			Name:      f.GoName,
			FieldType: fieldType,
			Comment:   getFieldComment(f.Comments),
		})
	}
	return fields
}

func getMethodComment(m *protogen.Method) string {
	symbol := "// "
	symbolLen := len(symbol)
	commentPrefix := symbol + m.GoName + " "
	comment := m.Comments.Leading.String()

	if len(comment) >= symbolLen {
		if comment[:symbolLen] == symbol {
			if comment[len(comment)-1] == '\n' {
				comment = comment[:len(comment)-1]
			}
			if len(comment) >= symbolLen {
				if len(comment[symbolLen:]) > len(m.GoName) {
					commentPrefixLower := strings.ToLower(comment[symbolLen : len(m.GoName)+symbolLen+1])
					if commentPrefixLower == strings.ToLower(m.GoName+" ") {
						return commentPrefix + comment[symbolLen+len(m.GoName)+1:]
					}
				}
				return commentPrefix + comment[symbolLen:]
			}
		}
	}

	return commentPrefix + "......"
}

func getFieldComment(commentSet protogen.CommentSet) string {
	comment1 := getFieldCommentStr(commentSet.Leading.String())
	comment2 := getFieldCommentStr(commentSet.Trailing.String())
	if comment1 == "" {
		return comment2
	}
	return comment1 + " " + comment2
}

func getFieldCommentStr(comment string) string {
	if len(comment) > 2 && comment[len(comment)-1] == '\n' {
		return comment[:len(comment)-1]
	}
	return comment
}

func getInvokeType(isStreamingClient bool, isStreamingServer bool) int {
	if isStreamingClient {
		if isStreamingServer {
			return 3 // bidirectional streaming
		}
		return 1 // client-side streaming
	}

	if isStreamingServer {
		return 2 // server-side streaming
	}

	return 0 // unary
}

func getProtoFileDir(protoPath string) string {
	ss := strings.Split(protoPath, "/")
	if len(ss) > 1 {
		return strings.Join(ss[:len(ss)-1], "/")
	}
	return ""
}

func convertToPkgName(importPath string) string {
	importPath = strings.ReplaceAll(importPath, `"`, "")
	ss := strings.Split(importPath, "/")
	l := len(ss)
	if l > 1 {
		pkgName := strings.ToLower(ss[l-1])
		if isVersionNum(pkgName) || pkgName == "pb" || len(pkgName) < 2 {
			return removeMiddleLine(ss[l-2]) + strings.ToUpper(pkgName[:1]) + pkgName[1:]
		}
		return removeMiddleLine(ss[l-1])
	}
	return ""
}

func isVersionNum(pkgName string) bool {
	pattern := `^v\d+$`
	matched, err := regexp.MatchString(pattern, pkgName)
	if err != nil {
		return false
	}
	return matched
}

func removeMiddleLine(str string) string {
	return strings.ReplaceAll(str, "-", "")
}

// GetImportPkg get import package
func GetImportPkg(services []*PbService) []byte {
	pkgMap := make(map[string]string)
	protoFileDir := ""
	moduleName := ""

	for _, service := range services {
		for key, val := range service.ImportPkgMap {
			pkgMap[key] = val
			protoFileDir = service.ProtoFileDir
			moduleName = service.ModuleName
		}
	}

	//pkgName := convertToPkgName(protoFileDir)
	//if _, ok := pkgMap[pkgName]; !ok {
	//	pkgMap[pkgName] = fmt.Sprintf(`%s "%s"`, pkgName, moduleName+"/"+protoFileDir)
	//}
	pkgName := convertToPkgName(protoFileDir)
	selfPkgPath := fmt.Sprintf(`%s "%s"`, pkgName, moduleName+"/"+protoFileDir)
	if _, ok := pkgMap[pkgName]; ok {
		pkgMap[pkgName] = selfPkgPath // real package path priority
	}

	var importPkg []string
	for _, v := range pkgMap {
		importPkg = append(importPkg, v)
	}
	if len(importPkg) == 0 {
		return []byte("")
	}

	return []byte(strings.Join(importPkg, "\n\t"))
}

// GetSourceImportPkg get source import package
func GetSourceImportPkg(services []*PbService) []byte {
	pkgMap := make(map[string]string)
	protoFileDir := ""
	moduleName := ""

	for _, service := range services {
		for key, val := range service.ImportPkgMap {
			pkgMap[key] = val
			protoFileDir = service.ProtoFileDir
			moduleName = service.ModuleName
			break
		}
	}

	pkgName := convertToPkgName(protoFileDir)
	return []byte(fmt.Sprintf(`%s "%s"`, pkgName, moduleName+"/"+protoFileDir))
}

// -------------------------------------------------------------------------------------------

// HTTPPbService http service fields
type HTTPPbService struct {
	Name      string // Greeter
	LowerName string // greeter first character to lower

	Methods       []*RPCMethod // service methods
	UniqueMethods []*RPCMethod

	ImportPkgMap map[string]string // [userV1]:[userV1 "user/api/user/v1"]
}

type HTTPPbServices []*HTTPPbService

// ParseHTTPPbServices parse protobuf services
func ParseHTTPPbServices(file *protogen.File) []*HTTPPbService {
	goImportPath := file.GoImportPath.String()

	var pss []*HTTPPbService
	for _, s := range file.Services {
		importPkgMap := map[string]string{}
		var methods []*RPCMethod
		for _, m := range s.Methods {
			ms := GetMethods(m, goImportPath)
			for _, method := range ms {
				for pkgPath := range method.ImportPkgPaths {
					pkgName := convertToPkgName(pkgPath)
					importPkgMap[pkgName] = pkgName + " " + pkgPath
				}
			}
			methods = append(methods, ms...)
		}

		pss = append(pss, &HTTPPbService{
			Name:          s.GoName,
			LowerName:     strings.ToLower(s.GoName[:1]) + s.GoName[1:],
			Methods:       methods,
			UniqueMethods: removeDuplicates(methods),
			ImportPkgMap:  importPkgMap,
		})
	}

	return pss
}

// MergeImportPkgPath merge import package path
func (services HTTPPbServices) MergeImportPkgPath() string {
	pkgMap := make(map[string]string)

	for _, service := range services {
		for key, val := range service.ImportPkgMap {
			pkgMap[key] = val
		}
	}

	var importPkg []string
	for _, v := range pkgMap {
		importPkg = append(importPkg, v)
	}

	if len(importPkg) == 0 {
		return ""
	}

	return strings.Join(importPkg, "\n\t")
}

func removeDuplicates(methods []*RPCMethod) []*RPCMethod {
	var uniqueMethods []*RPCMethod
	methodMap := make(map[string]struct{})
	for _, method := range methods {
		if _, ok := methodMap[method.Name]; !ok {
			methodMap[method.Name] = struct{}{}
			uniqueMethods = append(uniqueMethods, method)
		}
	}
	return uniqueMethods
}
