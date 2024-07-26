package parse

import (
	"fmt"
	"net/http"
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

	return methods
}

func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule) *RPCMethod {
	var (
		path       string
		method     string
		customKind string
		selector   = rule.Selector
	)

	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = http.MethodGet
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = http.MethodPut
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = http.MethodPost
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = http.MethodDelete
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = http.MethodPatch
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		customKind = strings.ToLower(pattern.Custom.Kind)
		method = http.MethodPost // default
	}
	md := buildMethodDesc(m, method, path, customKind, selector)
	return md
}

func buildMethodDesc(m *protogen.Method, httpMethod, path string, customKind string, selector string) *RPCMethod {
	defer func() {
		methodSets[m.GoName]++
	}()
	md := &RPCMethod{
		Name:       m.GoName,
		Num:        methodSets[m.GoName],
		Request:    m.Input.GoIdent.GoName,
		Reply:      m.Output.GoIdent.GoName,
		Path:       path,
		Method:     httpMethod,
		Selector:   selector,
		CustomKind: customKind,
		InvokeType: getInvokeType(m.Desc.IsStreamingClient(), m.Desc.IsStreamingServer()),
	}
	md.checkCustomKind()
	md.checkSelector()
	md.InitPathParams()
	return md
}

// RPCMethod describes a rpc method
type RPCMethod struct {
	Name       string // SayHello
	Num        int    // one rpc RPCMethod can correspond to multiple http requests
	Request    string // SayHelloReq
	Reply      string // SayHelloResp
	InvokeType int    // 0:unary, 1: client-side streaming, 2: server-side streaming, 3: bidirectional streaming

	// http_rule
	Path         string // rule
	Method       string // HTTP Method
	Body         string
	ResponseBody string

	CustomKind string
	Selector   string
	// if Selector is [ctx], and IsPassGinContext is true
	// if true, pass gin.Context to the rpc method
	IsPassGinContext bool
	// if Selector is [no_bind], IsPassGinContext and IsPassGinContext are both true
	// if true, ignore c.ShouldBindXXX for this method, you must use c.ShouldBindXXX() in rpc method
	IsIgnoreShouldBind bool
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

// parse selector and set custom control variables
func parseVariable(str string) (prefixStr string, isPassGinContext bool, isIgnoreShouldBind bool) {
	str = strings.ReplaceAll(str, " ", "")
	startIdx := strings.Index(str, "[")
	endIdx := strings.LastIndex(str, "]")
	if startIdx != -1 && endIdx != -1 {
		options := str[startIdx+1 : endIdx]
		ss := strings.Split(options, ",")
		for _, s := range ss {
			if s == "ctx" {
				isPassGinContext = true
			}
			if s == "no_bind" {
				isIgnoreShouldBind = true
				isPassGinContext = true // pass gin.Context
			}
		}
		prefixStr = str[:startIdx]
	} else {
		prefixStr = str
	}

	return prefixStr, isPassGinContext, isIgnoreShouldBind
}

func (m *RPCMethod) checkCustomKind() {
	if m.CustomKind == "" {
		return
	}

	customKindStr, isPassGinContext, isIgnoreShouldBind := parseVariable(m.CustomKind)
	m.IsPassGinContext = isPassGinContext
	m.IsIgnoreShouldBind = isIgnoreShouldBind

	switch customKindStr {
	case "get":
		m.Method = http.MethodGet
	case "post":
		m.Method = http.MethodPost
	case "put":
		m.Method = http.MethodPut
	case "delete":
		m.Method = http.MethodDelete
	case "patch":
		m.Method = http.MethodPatch
	case "options":
		m.Method = http.MethodOptions
	case "head":
		m.Method = http.MethodHead
	case "trace":
		m.Method = http.MethodTrace
	case "connect":
		m.Method = http.MethodConnect
	default:
		m.Method = http.MethodPost
	}
}

func (m *RPCMethod) checkSelector() {
	_, isPassGinContext, isIgnoreShouldBind := parseVariable(m.Selector)
	m.IsPassGinContext = isPassGinContext
	m.IsIgnoreShouldBind = isIgnoreShouldBind
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
