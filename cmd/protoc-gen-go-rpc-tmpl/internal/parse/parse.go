// Package parse is parsed proto file to struct
package parse

import (
	"math/rand"
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

	ServiceName      string // Greeter
	LowerServiceName string // greeter first character to lower
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
}

// RandNumber rand number 1~100
func (s *PbService) RandNumber() int {
	return rand.Intn(99) + 1
}

func parsePbService(s *protogen.Service) *PbService {
	var methods []*ServiceMethod
	for _, m := range s.Methods {
		methods = append(methods, &ServiceMethod{
			MethodName:    m.GoName,
			Request:       m.Input.GoIdent.GoName,
			RequestFields: getFields(m.Input),
			Reply:         m.Output.GoIdent.GoName,
			ReplyFields:   getFields(m.Output),
			Comment:       getMethodComment(m),

			ServiceName:      s.GoName,
			LowerServiceName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],
		})
	}

	return &PbService{
		Name:      s.GoName,
		LowerName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],
		Methods:   methods,
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
func GetServices(protoName string, file *protogen.File) []*PbService {
	var pss []*PbService
	for _, s := range file.Services {
		ps := parsePbService(s)
		ps.ProtoName = protoName
		pss = append(pss, ps)
	}
	return pss
}
