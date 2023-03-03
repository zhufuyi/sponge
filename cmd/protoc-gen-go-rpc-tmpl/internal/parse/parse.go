package parse

import (
	"math/rand"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// ServiceMethod method fields
type ServiceMethod struct {
	MethodName    string         // e.g. Create
	Request       string         // e.g. CreateRequest
	RequestFields []RequestField // request fields
	Reply         string         // e.g. CreateReply

	ServiceName      string // Greeter
	LowerServiceName string // greeter first character to lower
}

// RequestField request fields
type RequestField struct {
	FieldName string
	FieldType string
	Comment   string
}

// GoTypeZero default zero value for type
func (r RequestField) GoTypeZero() string {
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
			RequestFields: getRequestFields(m.Input.Fields),
			Reply:         m.Output.GoIdent.GoName,

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

func getRequestFields(fields []*protogen.Field) []RequestField {
	var reqFields []RequestField
	for _, field := range fields {
		reqFields = append(reqFields, RequestField{
			FieldName: field.GoName,
			FieldType: field.Desc.Kind().String(),
			Comment:   strings.ReplaceAll(field.Comments.Trailing.String(), "\n", ""),
		})
	}
	return reqFields
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
