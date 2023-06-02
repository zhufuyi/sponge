package parse

import (
	"math/rand"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// ServiceMethod RPCMethod fields
type ServiceMethod struct {
	MethodName string // Create
	Request    string // CreateRequest
	Reply      string // CreateReply

	ServiceName      string // Greeter
	LowerServiceName string // greeter first character to lower

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
			MethodName: m.GoName,
			Request:    m.Input.GoIdent.GoName,
			Reply:      m.Output.GoIdent.GoName,

			ServiceName:      s.GoName,
			LowerServiceName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],

			Path:   rpcMethod.Path,
			Method: rpcMethod.Method,
			Body:   rpcMethod.Body,
		})
	}

	cutServiceName := getCutServiceName(s.GoName)

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
