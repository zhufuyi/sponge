// Package parse is parsed proto file to struct
package parse

import (
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// ServiceMethod method fields
type ServiceMethod struct {
	MethodName    string   // e.g. Create
	Request       string   // e.g. CreateRequest
	RequestFields []*Field // request fields
	Reply         string   // e.g. CreateReply
	ReplyFields   []*Field
	Comment       string // e.g. Create a record
	InvokeType    int    // 0:unary, 1: client-side streaming, 2: server-side streaming, 3: bidirectional streaming

	ServiceName      string // Greeter
	LowerServiceName string // greeter first character to lower

	RequestImportPkgName string // e.g. userV1
	ReplyImportPkgName   string // e.g. userV1
	ProtoPkgName         string // e.g. userV1
}

// Field request fields
type Field struct {
	Name      string
	FieldType string
	Comment   string
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

// AddOne counter
func (t *ServiceMethod) AddOne(i int) int {
	return i + 1
}

// PbService service fields
type PbService struct {
	Name      string           // Greeter
	LowerName string           // greeter first character to lower
	ProtoName string           // proto file name greeter.proto
	Methods   []*ServiceMethod // service methods

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
	importPkgMap := map[string]string{}

	var methods []*ServiceMethod
	for _, m := range s.Methods {
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

			ServiceName:      s.GoName,
			LowerServiceName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],

			RequestImportPkgName: requestImportPkgName,
			ReplyImportPkgName:   replyImportPkgName,
			ProtoPkgName:         protoPkgName,
		})
	}

	return &PbService{
		Name:         s.GoName,
		LowerName:    strings.ToLower(s.GoName[:1]) + s.GoName[1:],
		Methods:      methods,
		ImportPkgMap: importPkgMap,
		ProtoFileDir: protoFileDir,
		ProtoPkgName: protoPkgName,
		ModuleName:   moduleName,
	}
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
			Comment:   getComment(f.Comments),
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

func getComment(commentSet protogen.CommentSet) string {
	comment1 := getCommentStr(commentSet.Leading.String())
	comment2 := getCommentStr(commentSet.Trailing.String())
	if comment1 == "" {
		return comment2
	}
	return comment1 + " " + comment2
}

func getCommentStr(comment string) string {
	if len(comment) > 2 && comment[len(comment)-1] == '\n' {
		return comment[:len(comment)-1]
	}
	return comment
}

// GetServices parse protobuf services
func GetServices(file *protogen.File, moduleName string) []*PbService {
	protoFileDir := getProtoFileDir(file.GeneratedFilenamePrefix)
	protoName := getProtoFilename(file.GeneratedFilenamePrefix)

	var pss []*PbService
	for _, s := range file.Services {
		ps := parsePbService(s, protoFileDir, moduleName)
		ps.ProtoName = protoName
		pss = append(pss, ps)
	}
	return pss
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

func getProtoFilename(filenamePrefix string) string {
	filenamePrefix = strings.ReplaceAll(filenamePrefix, ".proto", "")
	filenamePrefix = strings.ReplaceAll(filenamePrefix, getPathDelimiter(), "/")
	ss := strings.Split(filenamePrefix, "/")

	if len(ss) == 0 {
		return ""
	} else if len(ss) == 1 {
		return ss[0] + ".proto"
	}

	return ss[len(ss)-1] + ".proto"
}

func getPathDelimiter() string {
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}

	return delimiter
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

	pkgName := convertToPkgName(protoFileDir)
	if _, ok := pkgMap[pkgName]; !ok {
		pkgMap[pkgName] = fmt.Sprintf(`%s "%s"`, pkgName, moduleName+"/"+protoFileDir)
	}

	var importPkg []string
	for _, v := range pkgMap {
		importPkg = append(importPkg, v)
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
		}
	}

	pkgName := convertToPkgName(protoFileDir)
	if v, ok := pkgMap[pkgName]; ok {
		return []byte(v)
	}

	return []byte(fmt.Sprintf(`%s "%s"`, pkgName, moduleName+"/"+protoFileDir))
}
