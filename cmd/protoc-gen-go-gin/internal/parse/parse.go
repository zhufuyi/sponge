package parse

import (
	"math/rand"
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

	ServiceName         string // Greeter
	LowerServiceName    string // greeter first character to lower
	LowerCutServiceName string // GreeterService --> greeter

	// http_rule
	Path   string // rule
	Method string // HTTP Method
	Body   string
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
}

// RandNumber rand number 1~100
func (s *PbService) RandNumber() int {
	return rand.Intn(99) + 1
}

func parsePbService(s *protogen.Service) *PbService {
	cutServiceName := getCutServiceName(s.GoName)

	var methods []*ServiceMethod
	for _, m := range s.Methods {
		rpcMethod := &RPCMethod{} //nolint
		rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule != nil && ok {
			rpcMethod = buildHTTPRule(m, rule)
		} else {
			rpcMethod = defaultMethod(m)
		}

		methods = append(methods, &ServiceMethod{
			MethodName:    m.GoName,
			Request:       m.Input.GoIdent.GoName,
			RequestFields: getFields(m.Input),
			Reply:         m.Output.GoIdent.GoName,
			ReplyFields:   getFields(m.Output),
			Comment:       getMethodComment(m),

			ServiceName:         s.GoName,
			LowerServiceName:    strings.ToLower(s.GoName[:1]) + s.GoName[1:],
			LowerCutServiceName: strings.ToLower(cutServiceName[:1]) + cutServiceName[1:],

			Path:   rpcMethod.Path,
			Method: rpcMethod.Method,
			Body:   rpcMethod.Body,
		})
	}

	return &PbService{
		Name:                s.GoName,
		LowerName:           strings.ToLower(s.GoName[:1]) + s.GoName[1:],
		Methods:             methods,
		CutServiceName:      cutServiceName,
		LowerCutServiceName: strings.ToLower(cutServiceName[:1]) + cutServiceName[1:],
	}
}

// GetServices parse protobuf services
func GetServices(file *protogen.File) []*PbService {
	var pss []*PbService
	for _, s := range file.Services {
		pss = append(pss, parsePbService(s))
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

//func getRequestFields(m *protogen.Method) []*Field {
//	var fields []*Field
//	for _, f := range m.Input.Fields {
//		fieldType := f.Desc.Kind().String()
//		if f.Desc.Cardinality().String() == "repeated" {
//			fieldType = "[]" + fieldType
//		}
//		fields = append(fields, &Field{
//			Name:      f.GoName,
//			FieldType: fieldType,
//			Comment:   getFieldComment(f.Comments),
//		})
//	}
//	return fields
//}
//
//func getReplyFields(m *protogen.Method) []*Field {
//	var fields []*Field
//	for _, f := range m.Output.Fields {
//		fieldType := f.Desc.Kind().String()
//		if f.Desc.Cardinality().String() == "repeated" {
//			fieldType = "[]" + fieldType
//		}
//		fields = append(fields, &Field{
//			Name:      f.GoName,
//			FieldType: fieldType,
//			Comment:   getFieldComment(f.Comments),
//		})
//	}
//	return fields
//}

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
