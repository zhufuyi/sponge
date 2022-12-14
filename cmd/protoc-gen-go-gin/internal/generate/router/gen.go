package router

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	contextPkg         = protogen.GoImportPath("context")
	errcodePkg         = protogen.GoImportPath("github.com/zhufuyi/sponge/pkg/errcode")
	middlewarePkg      = protogen.GoImportPath("github.com/zhufuyi/sponge/pkg/gin/middleware")
	zapPkg             = protogen.GoImportPath("go.uber.org/zap")
	ginPkg             = protogen.GoImportPath("github.com/gin-gonic/gin")
	deprecationComment = "// Deprecated: Do not use."
)

var methodSets = make(map[string]int)

// GenerateFile generates a *_router.pb.go file.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + "_router.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-go-gin. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	g.P("// import packages: ", contextPkg.Ident(" "), errcodePkg.Ident(" "),
		middlewarePkg.Ident(" "), zapPkg.Ident(" "), ginPkg.Ident(" "))
	g.P()

	for _, s := range file.Services {
		genService(file, g, s)
	}
	return g
}

func genService(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}

	// HTTP Server.
	sd := &tmplField{
		Name:      s.GoName,
		LowerName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],
		FullName:  string(s.Desc.FullName()),
		FilePath:  file.Desc.Path(),
	}

	for _, m := range s.Methods {
		sd.Methods = append(sd.Methods, genMethod(m)...)
	}

	g.P(sd.execute())
}

func genMethod(m *protogen.Method) []*method {
	var methods []*method

	// http rule config
	rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	if rule != nil && ok {
		for _, bind := range rule.AdditionalBindings {
			methods = append(methods, buildHTTPRule(m, bind))
		}
		methods = append(methods, buildHTTPRule(m, rule))
		return methods
	}

	// default http method mapping
	methods = append(methods, defaultMethod(m))
	return methods
}

// defaultMethodPath generates an http route based on the function name
// If the first word of the method name is not an http method mapping, then POST is returned by default
func defaultMethod(m *protogen.Method) *method {
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

func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule) *method {
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

func buildMethodDesc(m *protogen.Method, httpMethod, path string) *method {
	defer func() { methodSets[m.GoName]++ }()
	md := &method{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: m.Input.GoIdent.GoName,
		Reply:   m.Output.GoIdent.GoName,
		Path:    path,
		Method:  httpMethod,
	}
	md.initPathParams()
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

// ------------------------------------------------------------------------------------------

type tmplField struct {
	Name      string // Greeter
	LowerName string // greeter
	FullName  string // v1.Greeter
	FilePath  string // api/v1/demo.proto

	Methods   []*method
	MethodSet map[string]*method
}

func (s *tmplField) execute() string {
	if s.MethodSet == nil {
		s.MethodSet = map[string]*method{}
		for _, m := range s.Methods {
			m := m
			s.MethodSet[m.Name] = m
		}
	}

	buf := new(bytes.Buffer)
	if err := handlerTmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return buf.String()
}

type method struct {
	Name    string // SayHello
	Num     int    // one rpc method can correspond to multiple http requests
	Request string // SayHelloReq
	Reply   string // SayHelloResp
	// http_rule
	Path         string // rule
	Method       string // HTTP Method
	Body         string
	ResponseBody string
}

// HandlerName for gin handler name
func (m *method) HandlerName() string {
	return fmt.Sprintf("%s_%d", m.Name, m.Num)
}

// HasPathParams whether to include routing parameters
func (m *method) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}
	return false
}

// initPathParams conversion parameter routing {xx} --> :xx
func (m *method) initPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}
	m.Path = strings.Join(paths, "/")
}
