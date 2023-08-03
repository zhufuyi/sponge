package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var methodSets = make(map[string]int)

// GetMethods get rpc method descriptions
func GetMethods(m *protogen.Method) []*RPCMethod {
	var methods []*RPCMethod

	// http rule config
	rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	if rule != nil && ok {
		for _, bind := range rule.AdditionalBindings {
			methods = append(methods, buildHTTPRule(m, bind))
		}
		methods = append(methods, buildHTTPRule(m, rule))
		return methods
	}

	// default http RPCMethod mapping
	methods = append(methods, defaultMethod(m))
	return methods
}

// defaultMethodPath generates a http route based on the function name
// If the first word of the RPCMethod name is not a http RPCMethod mapping, then POST is returned by default
func defaultMethod(m *protogen.Method) *RPCMethod {
	names := strings.Split(toSnakeCase(m.GoName), "_")
	var (
		paths      []string
		httpMethod string
		path       string
	)

	switch strings.ToUpper(names[0]) {
	case http.MethodGet, "FIND", "QUERY", "LIST", "SEARCH":
		httpMethod = http.MethodGet
	case http.MethodPost, "CREATE":
		httpMethod = http.MethodPost
	case http.MethodPut, "UPDATE":
		httpMethod = http.MethodPut
	case http.MethodPatch:
		httpMethod = http.MethodPatch
	case http.MethodDelete:
		httpMethod = http.MethodDelete
	default:
		httpMethod = "POST"
		paths = names
	}

	if len(paths) > 0 {
		path = strings.Join(paths, "/")
	}

	if len(names) > 1 {
		path = strings.Join(names[1:], "/")
	}

	md := buildMethodDesc(m, httpMethod, path)
	md.Body = "*"
	return md
}

func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule) *RPCMethod {
	var (
		path   string
		method string
	)
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = "GET"
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = "PUT"
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = "POST"
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = "DELETE"
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = "PATCH"
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}
	md := buildMethodDesc(m, method, path)
	return md
}

func buildMethodDesc(m *protogen.Method, httpMethod, path string) *RPCMethod {
	defer func() { methodSets[m.GoName]++ }()
	md := &RPCMethod{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: m.Input.GoIdent.GoName,
		Reply:   m.Output.GoIdent.GoName,
		Path:    path,
		Method:  httpMethod,
	}
	md.InitPathParams()
	return md
}

var matchFirstCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(input string) string {
	output := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	output = matchAllCap.ReplaceAllString(output, "${1}_${2}")
	output = strings.ReplaceAll(output, "-", "_")
	return strings.ToLower(output)
}

// RPCMethod describes a rpc method
type RPCMethod struct {
	Name    string // SayHello
	Num     int    // one rpc RPCMethod can correspond to multiple http requests
	Request string // SayHelloReq
	Reply   string // SayHelloResp

	// http_rule
	Path         string // rule
	Method       string // HTTP Method
	Body         string
	ResponseBody string
}

// HandlerName for gin handler name
func (m *RPCMethod) HandlerName() string {
	return fmt.Sprintf("%s_%d", m.Name, m.Num)
}

// HasPathParams whether to include routing parameters
func (m *RPCMethod) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}
	return false
}

// InitPathParams conversion parameter routing {xx} --> :xx
func (m *RPCMethod) InitPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}') {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}
	m.Path = strings.Join(paths, "/")
}
