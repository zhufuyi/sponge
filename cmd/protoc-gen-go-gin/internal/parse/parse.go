package parse

import (
	"math/rand"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// ServiceMethod method fields
type ServiceMethod struct {
	MethodName string // Create
	Request    string // CreateRequest
	Reply      string // CreateReply

	ServiceName      string // Greeter
	LowerServiceName string // greeter first character to lower
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
}

// RandNumber rand number 1~100
func (s *PbService) RandNumber() int {
	return rand.Intn(99) + 1
}

func parsePbService(s *protogen.Service) *PbService {
	var methods []*ServiceMethod
	for _, m := range s.Methods {
		methods = append(methods, &ServiceMethod{
			MethodName: m.GoName,
			Request:    m.Input.GoIdent.GoName,
			Reply:      m.Output.GoIdent.GoName,

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

// GetServices parse protobuf services
func GetServices(file *protogen.File) []*PbService {
	var pss []*PbService
	for _, s := range file.Services {
		pss = append(pss, parsePbService(s))
	}
	return pss
}
