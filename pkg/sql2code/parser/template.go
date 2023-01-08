package parser

import (
	"github.com/pkg/errors"
	"sync"
	"text/template"
)

var (
	modelStructTmpl    *template.Template
	modelStructTmplRaw = `
{{- if .Comment -}}
// {{.TableName}} {{.Comment}}
{{end -}}
type {{.TableName}} struct {
{{- range .Fields}}
	{{.Name}} {{.GoType}} {{if .Tag}}` + "`{{.Tag}}`" + `{{end}}{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
{{if .NameFunc}}
// TableName table name
func (m *{{.TableName}}) TableName() string {
	return "{{.RawTableName}}"
}
{{end}}
`

	modelTmpl    *template.Template
	modelTmplRaw = `package {{.Package}}
{{if .ImportPath}}
import (
	{{- range .ImportPath}}
	"{{.}}"
	{{- end}}
)
{{- end}}
{{range .StructCode}}
{{.}}
{{end}}`

	updateFieldTmpl    *template.Template
	updateFieldTmplRaw = `
{{- range .Fields}}
	if table.{{.Name}} {{.ConditionZero}} {
		update["{{.ColName}}"] = table.{{.Name}}
	}
{{- end}}`

	handlerCreateStructTmpl    *template.Template
	handlerCreateStructTmplRaw = `
// Create{{.TableName}}Request create params
// todo fill in the binding rules https://github.com/go-playground/validator
type Create{{.TableName}}Request struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerUpdateStructTmpl    *template.Template
	handlerUpdateStructTmplRaw = `
// Update{{.TableName}}ByIDRequest update params
type Update{{.TableName}}ByIDRequest struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerDetailStructTmpl    *template.Template
	handlerDetailStructTmplRaw = `
// Get{{.TableName}}ByIDRespond respond detail
type Get{{.TableName}}ByIDRespond struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}`

	modelJSONTmpl    *template.Template
	modelJSONTmplRaw = `{
{{- range .Fields}}
	"{{.ColName}}" {{.GoZero}}
{{- end}}
}
`

	protoFileTmpl    *template.Template
	protoFileTmplRaw = `syntax = "proto3";

package api.serverNameExample.v1;

import "api/types/types.proto";
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "tagger/tagger.proto";
//import "validate/validate.proto";

option go_package = "github.com/zhufuyi/sponge/api/serverNameExample/v1;v1";

// Default settings for generating swagger documents
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  host: "localhost:8080"
  base_path: ""
  info: {
    title: "serverNameExample api docs";
    version: "v0.0.0";
  };
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
};

service {{.TName}}Service {
  rpc Create(Create{{.TableName}}Request) returns (Create{{.TableName}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}"
      body: "*"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "create a new {{.TName}}",
      description: "submit information to create a new {{.TName}}",
      tags: "{{.TName}}",
    };
  }

  rpc DeleteByID(Delete{{.TableName}}ByIDRequest) returns (Delete{{.TableName}}ByIDReply) {
    option (google.api.http) = {
      delete: "/api/v1/{{.TName}}/{id}"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "delete {{.TName}}",
      description: "delete {{.TName}} by id",
      tags: "{{.TName}}",
    };
  }

  rpc UpdateByID(Update{{.TableName}}ByIDRequest) returns (Update{{.TableName}}ByIDReply) {
    option (google.api.http) = {
      put: "/api/v1/{{.TName}}/{id}"
      body: "*"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "update {{.TName}} info",
      description: "update {{.TName}} info by id",
      tags: "{{.TName}}",
    };
  }

  rpc GetByID(Get{{.TableName}}ByIDRequest) returns (Get{{.TableName}}ByIDReply) {
    option (google.api.http) = {
      get: "/api/v1/{{.TName}}/{id}"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "get {{.TName}} details",
      description: "get {{.TName}} details by id",
      tags: "{{.TName}}",
    };
  }

  rpc ListByIDs(List{{.TableName}}ByIDsRequest) returns (List{{.TableName}}ByIDsReply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}s/ids"
      body: "*"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "get a list of {{.TName}} based on multiple ids",
      description: "get a list of {{.TName}} based on multiple ids",
      tags: "{{.TName}}",
    };
  }

  rpc List(List{{.TableName}}Request) returns (List{{.TableName}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}s"
      body: "*"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "get a list of {{.TName}} based on query parameters",
      description: "get a list of {{.TName}} based on query parameters",
      tags: "{{.TName}}",
    };
  }
}

// todo fill in the validate rules https://github.com/envoyproxy/protoc-gen-validate#constraint-rules

// protoMessageCreateCode

message Create{{.TableName}}Reply {
  uint64   id =1;
}

message Delete{{.TableName}}ByIDRequest {
  uint64   id =1 [(tagger.tags) = "uri:\"id\"" ];
}

message Delete{{.TableName}}ByIDReply {

}

// protoMessageUpdateCode

message Update{{.TableName}}ByIDReply {

}

// protoMessageDetailCode

message Get{{.TableName}}ByIDRequest {
  uint64   id =1 [(tagger.tags) = "uri:\"id\"" ];
}

message Get{{.TableName}}ByIDReply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}ByIDsRequest {
  repeated uint64 ids = 1;
}

message List{{.TableName}}ByIDsReply {
  repeated {{.TableName}} {{.TName}}s = 1;
}

message List{{.TableName}}Request {
  types.Params params = 1;
}

message List{{.TableName}}Reply {
  int64 total =1;
  repeated {{.TableName}} {{.TName}}s = 2;
}
`

	protoMessageCreateTmpl    *template.Template
	protoMessageCreateTmplRaw = `message Create{{.TableName}}Request {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.ColName}} = {{$v.AddOne $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	protoMessageUpdateTmpl    *template.Template
	protoMessageUpdateTmplRaw = `message Update{{.TableName}}ByIDRequest {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.ColName}} = {{$v.AddOneWithTag $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	protoMessageDetailTmpl    *template.Template
	protoMessageDetailTmplRaw = `message {{.TableName}} {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.ColName}} = {{$v.AddOne $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	serviceStructTmpl    *template.Template
	serviceStructTmplRaw = `
		{
			name: "Create",
			fn: func() (interface{}, error) {
				// todo enter parameters before testing
// serviceCreateStructCode
			},
			wantErr: false,
		},

		{
			name: "UpdateByID",
			fn: func() (interface{}, error) {
				// todo enter parameters before testing
// serviceUpdateStructCode
			},
			wantErr: false,
		},
`

	serviceCreateStructTmpl    *template.Template
	serviceCreateStructTmplRaw = `				return cli.Create(ctx, &serverNameExampleV1.Create{{.TableName}}Request{
					{{- range .Fields}}
						{{.Name}}:  {{.GoTypeZero}}, {{if .Comment}} // {{.Comment}}{{end}}
					{{- end}}
				})`

	serviceUpdateStructTmpl    *template.Template
	serviceUpdateStructTmplRaw = `				return cli.UpdateByID(ctx, &serverNameExampleV1.Update{{.TableName}}ByIDRequest{
					{{- range .Fields}}
						{{.Name}}:  {{.GoTypeZero}}, {{if .Comment}} // {{.Comment}}{{end}}
					{{- end}}
				})`

	tmplParseOnce sync.Once
)

func initTemplate() {
	tmplParseOnce.Do(func() {
		var err, errSum error

		modelStructTmpl, err = template.New("goStruct").Parse(modelStructTmplRaw)
		if err != nil {
			errSum = errors.Wrap(err, "modelStructTmplRaw")
		}
		modelTmpl, err = template.New("goFile").Parse(modelTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "modelTmplRaw:"+err.Error())
		}
		updateFieldTmpl, err = template.New("goUpdateField").Parse(updateFieldTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "updateFieldTmplRaw:"+err.Error())
		}
		handlerCreateStructTmpl, err = template.New("goPostStruct").Parse(handlerCreateStructTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "handlerCreateStructTmplRaw:"+err.Error())
		}
		handlerUpdateStructTmpl, err = template.New("goPutStruct").Parse(handlerUpdateStructTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "handlerUpdateStructTmplRaw:"+err.Error())
		}
		handlerDetailStructTmpl, err = template.New("goGetStruct").Parse(handlerDetailStructTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "handlerDetailStructTmplRaw:"+err.Error())
		}
		modelJSONTmpl, err = template.New("modelJSON").Parse(modelJSONTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "modelJSONTmplRaw:"+err.Error())
		}
		protoFileTmpl, err = template.New("protoFile").Parse(protoFileTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoFileTmplRaw:"+err.Error())
		}
		protoMessageCreateTmpl, err = template.New("protoMessageCreate").Parse(protoMessageCreateTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoMessageCreateTmplRaw:"+err.Error())
		}
		protoMessageUpdateTmpl, err = template.New("protoMessageUpdate").Parse(protoMessageUpdateTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoMessageUpdateTmplRaw:"+err.Error())
		}
		protoMessageDetailTmpl, err = template.New("protoMessageDetail").Parse(protoMessageDetailTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoMessageDetailTmplRaw:"+err.Error())
		}
		serviceCreateStructTmpl, err = template.New("serviceCreateStruct").Parse(serviceCreateStructTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "serviceCreateStructTmplRaw:"+err.Error())
		}
		serviceUpdateStructTmpl, err = template.New("serviceUpdateStruct").Parse(serviceUpdateStructTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "serviceUpdateStructTmplRaw:"+err.Error())
		}
		serviceStructTmpl, err = template.New("serviceStruct").Parse(serviceStructTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "serviceStructTmplRaw:"+err.Error())
		}

		if errSum != nil {
			panic(errSum)
		}
	})
}
