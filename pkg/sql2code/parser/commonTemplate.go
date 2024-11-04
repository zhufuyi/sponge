package parser

import (
	"sync"
	"text/template"

	"github.com/pkg/errors"
)

// nolint
var (
	handlerCreateStructCommonTmpl    *template.Template
	handlerCreateStructCommonTmplRaw = `
// Create{{.TableName}}Request request params
type Create{{.TableName}}Request struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.JSONName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerUpdateStructCommonTmpl    *template.Template
	handlerUpdateStructCommonTmplRaw = `
// Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request request params
type Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.JSONName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerDetailStructCommonTmpl    *template.Template
	handlerDetailStructCommonTmplRaw = `
// {{.TableName}}ObjDetail detail
type {{.TableName}}ObjDetail struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.JSONName}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}`

	protoFileCommonTmpl    *template.Template
	protoFileCommonTmplRaw = `syntax = "proto3";

package api.serverNameExample.v1;

import "api/types/types.proto";
import "validate/validate.proto";

option go_package = "github.com/zhufuyi/sponge/api/serverNameExample/v1;v1";

service {{.TName}} {
  // create {{.TName}}
  rpc Create(Create{{.TableName}}Request) returns (Create{{.TableName}}Reply) {}

  // delete {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc DeleteBy{{.CrudInfo.ColumnNameCamel}}(Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {}

  // update {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc UpdateBy{{.CrudInfo.ColumnNameCamel}}(Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {}

  // get {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc GetBy{{.CrudInfo.ColumnNameCamel}}(Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {}

  // list of {{.TName}} by query parameters
  rpc List(List{{.TableName}}Request) returns (List{{.TableName}}Reply) {}

  // delete {{.TName}} by batch {{.CrudInfo.ColumnNameCamelFCL}}
  rpc DeleteBy{{.CrudInfo.ColumnNamePluralCamel}}(Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request) returns (Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply) {}

  // get {{.TName}} by condition
  rpc GetByCondition(Get{{.TableName}}ByConditionRequest) returns (Get{{.TableName}}ByConditionReply) {}

  // list of {{.TName}} by batch {{.CrudInfo.ColumnNameCamelFCL}}
  rpc ListBy{{.CrudInfo.ColumnNamePluralCamel}}(List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request) returns (List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply) {}

  // list {{.TName}} by last {{.CrudInfo.ColumnNameCamelFCL}}
  rpc ListByLast{{.CrudInfo.ColumnNameCamel}}(List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Request) returns (List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Reply) {}
}


/*
Notes for defining message fields:
    1. Suggest using camel case style naming for message field names, such as firstName, lastName, etc.
    2. If the message field name ending in 'id', it is recommended to use xxxID naming format, such as userID, orderID, etc.
    3. Add validate rules https://github.com/envoyproxy/protoc-gen-validate#constraint-rules, such as:
        uint64 id = 1 [(validate.rules).uint64.gte  = 1];
*/


// protoMessageCreateCode

message Create{{.TableName}}Reply {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1;
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetGRPCProtoValidation}};
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageUpdateCode

message Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageDetailCode

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetGRPCProtoValidation}};
}

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}Request {
  api.types.Params params = 1;
}

message List{{.TableName}}Reply {
  int64 total = 1;
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 2;
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request {
  repeated {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNamePluralCamelFCL}} = 1 [(validate.rules).repeated.min_items = 1];
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply {

}

message Get{{.TableName}}ByConditionRequest {
  types.Conditions conditions = 1;
}

message Get{{.TableName}}ByConditionReply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request {
  repeated {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNamePluralCamelFCL}} = 1 [(validate.rules).repeated.min_items = 1];
}

message List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply {
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 1;
}

message List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} last{{.CrudInfo.ColumnNameCamel}} = 1;
  uint32 limit = 2 [(validate.rules).uint32.gt = 0]; // limit size per page
  string sort = 3; // sort by column name of table, default is -{{.CrudInfo.ColumnName}}, the - sign indicates descending order.
}

message List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Reply {
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 1;
}
`

	protoFileSimpleCommonTmpl    *template.Template
	protoFileSimpleCommonTmplRaw = `syntax = "proto3";

package api.serverNameExample.v1;

import "api/types/types.proto";
import "validate/validate.proto";

option go_package = "github.com/zhufuyi/sponge/api/serverNameExample/v1;v1";

service {{.TName}} {
  // create {{.TName}}
  rpc Create(Create{{.TableName}}Request) returns (Create{{.TableName}}Reply) {}

  // delete {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc DeleteBy{{.CrudInfo.ColumnNameCamel}}(Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {}

  // update {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc UpdateBy{{.CrudInfo.ColumnNameCamel}}(Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {}

  // get {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc GetBy{{.CrudInfo.ColumnNameCamel}}(Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {}

  // list of {{.TName}} by query parameters
  rpc List(List{{.TableName}}Request) returns (List{{.TableName}}Reply) {}
}


/*
Notes for defining message fields:
    1. Suggest using camel case style naming for message field names, such as firstName, lastName, etc.
    2. If the message field name ending in 'id', it is recommended to use xxxID naming format, such as userID, orderID, etc.
    3. Add validate rules https://github.com/envoyproxy/protoc-gen-validate#constraint-rules, such as:
        uint64 id = 1 [(validate.rules).uint64.gte  = 1];
*/


// protoMessageCreateCode

message Create{{.TableName}}Reply {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1;
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetGRPCProtoValidation}};
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageUpdateCode

message Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageDetailCode

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetGRPCProtoValidation}};
}

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}Request {
  api.types.Params params = 1;
}

message List{{.TableName}}Reply {
  int64 total = 1;
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 2;
}
`

	protoFileForWebCommonTmpl    *template.Template
	protoFileForWebCommonTmplRaw = `syntax = "proto3";

package api.serverNameExample.v1;

import "api/types/types.proto";
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "tagger/tagger.proto";
import "validate/validate.proto";

option go_package = "github.com/zhufuyi/sponge/api/serverNameExample/v1;v1";

/*
Reference https://github.com/grpc-ecosystem/grpc-gateway/blob/db7fbefff7c04877cdb32e16d4a248a024428207/examples/internal/proto/examplepb/a_bit_of_everything.proto
Default settings for generating swagger documents
NOTE: because json does not support 64 bits, the int64 and uint64 types under *.swagger.json are automatically converted to string types
Tips: add swagger option to rpc method, example:
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "get user by id",
      description: "get user by id",
      security: {
        security_requirement: {
          key: "BearerAuth";
          value: {}
        }
      }
    };
*/
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  host: "localhost:8080"
  base_path: ""
  info: {
    title: "serverNameExample api docs";
    version: "2.0";
  }
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
  security_definitions: {
    security: {
      key: "BearerAuth";
      value: {
        type: TYPE_API_KEY;
        in: IN_HEADER;
        name: "Authorization";
        description: "Type Bearer your-jwt-token to Value";
      }
    }
  }
};

service {{.TName}} {
  // create {{.TName}}
  rpc Create(Create{{.TableName}}Request) returns (Create{{.TableName}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}"
      body: "*"
    };
  }

  // delete {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc DeleteBy{{.CrudInfo.ColumnNameCamel}}(Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {
    option (google.api.http) = {
      delete: "/api/v1/{{.TName}}/left_curly_bracket{{.CrudInfo.ColumnNameCamelFCL}}right_curly_bracket"
    };
  }

  // update {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc UpdateBy{{.CrudInfo.ColumnNameCamel}}(Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {
    option (google.api.http) = {
      put: "/api/v1/{{.TName}}/left_curly_bracket{{.CrudInfo.ColumnNameCamelFCL}}right_curly_bracket"
      body: "*"
    };
  }

  // get {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc GetBy{{.CrudInfo.ColumnNameCamel}}(Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {
    option (google.api.http) = {
      get: "/api/v1/{{.TName}}/left_curly_bracket{{.CrudInfo.ColumnNameCamelFCL}}right_curly_bracket"
    };
  }

  // list of {{.TName}} by query parameters
  rpc List(List{{.TableName}}Request) returns (List{{.TableName}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}/list"
      body: "*"
    };
  }

  // delete {{.TName}} by batch {{.CrudInfo.ColumnNameCamelFCL}}
  rpc DeleteBy{{.CrudInfo.ColumnNamePluralCamel}}(Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request) returns (Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}/delete/ids"
      body: "*"
    };
  }

  // get {{.TName}} by condition
  rpc GetByCondition(Get{{.TableName}}ByConditionRequest) returns (Get{{.TableName}}ByConditionReply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}/condition"
      body: "*"
    };
  }

  // list of {{.TName}} by batch {{.CrudInfo.ColumnNameCamelFCL}}
  rpc ListBy{{.CrudInfo.ColumnNamePluralCamel}}(List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request) returns (List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}/list/ids"
      body: "*"
    };
  }

  // list {{.TName}} by last {{.CrudInfo.ColumnNameCamelFCL}}
  rpc ListByLast{{.CrudInfo.ColumnNameCamel}}(List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Request) returns (List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Reply) {
    option (google.api.http) = {
      get: "/api/v1/{{.TName}}/list"
    };
  }
}


/*
Notes for defining message fields:
    1. Suggest using camel case style naming for message field names, such as firstName, lastName, etc.
    2. If the message field name ending in 'id', it is recommended to use xxxID naming format, such as userID, orderID, etc.
    3. Add validate rules https://github.com/envoyproxy/protoc-gen-validate#constraint-rules, such as:
        uint64 id = 1 [(validate.rules).uint64.gte  = 1];

If used to generate code that supports the HTTP protocol, notes for defining message fields:
    1. If the route contains the path parameter, such as /api/v1/userExample/{id}, the defined
        message must contain the name of the path parameter and the name should be added
        with a new tag, such as int64 id = 1 [(tagger.tags) = "uri:\"id\""];
    2. If the request url is followed by a query parameter, such as /api/v1/getUserExample?name=Tom,
        a form tag must be added when defining the query parameter in the message, such as:
        string name = 1 [(tagger.tags) = "form:\"name\""].
    3. If the message field name contain underscores(such as 'field_name'), it will cause a problem
        where the JSON field names of the Swagger request parameters are different from those of the
        GRPC JSON tag names. There are two solutions: Solution 1, remove the underline from the
         message field name. Option 2, use the tool 'protoc-go-inject-tag' to modify the JSON tag name,
         such as: string first_name = 1 ; // @gotags: json:"firstName"
*/


// protoMessageCreateCode

message Create{{.TableName}}Reply {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1;
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetWebProtoValidation}};
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageUpdateCode

message Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageDetailCode

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetWebProtoValidation}};
}

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}Request {
  api.types.Params params = 1;
}

message List{{.TableName}}Reply {
  int64 total = 1;
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 2;
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request {
  repeated {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNamePluralCamelFCL}} = 1 [(validate.rules).repeated.min_items = 1];
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply {

}

message Get{{.TableName}}ByConditionRequest {
  types.Conditions conditions = 1;
}

message Get{{.TableName}}ByConditionReply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Request {
  repeated {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNamePluralCamelFCL}} = 1 [(validate.rules).repeated.min_items = 1];
}

message List{{.TableName}}By{{.CrudInfo.ColumnNamePluralCamel}}Reply {
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 1;
}

message List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} last{{.CrudInfo.ColumnNameCamel}} = 1 [(tagger.tags) = "form:\"last{{.CrudInfo.ColumnNameCamel}}\""];
  uint32 limit = 2 [(validate.rules).uint32.gt = 0, (tagger.tags) = "form:\"limit\""]; // limit size per page
  string sort = 3 [(tagger.tags) = "form:\"sort\""]; // sort by column name of table, default is -{{.CrudInfo.ColumnName}}, the - sign indicates descending order.
}

message List{{.TableName}}ByLast{{.CrudInfo.ColumnNameCamel}}Reply {
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 1;
}
`

	protoFileForSimpleWebCommonTmpl    *template.Template
	protoFileForSimpleWebCommonTmplRaw = `syntax = "proto3";

package api.serverNameExample.v1;

import "api/types/types.proto";
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "tagger/tagger.proto";
import "validate/validate.proto";

option go_package = "github.com/zhufuyi/sponge/api/serverNameExample/v1;v1";

/*
Reference https://github.com/grpc-ecosystem/grpc-gateway/blob/db7fbefff7c04877cdb32e16d4a248a024428207/examples/internal/proto/examplepb/a_bit_of_everything.proto
Default settings for generating swagger documents
NOTE: because json does not support 64 bits, the int64 and uint64 types under *.swagger.json are automatically converted to string types
Tips: add swagger option to rpc method, example:
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "get user by id",
      description: "get user by id",
      security: {
        security_requirement: {
          key: "BearerAuth";
          value: {}
        }
      }
    };
*/
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  host: "localhost:8080"
  base_path: ""
  info: {
    title: "serverNameExample api docs";
    version: "2.0";
  }
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
  security_definitions: {
    security: {
      key: "BearerAuth";
      value: {
        type: TYPE_API_KEY;
        in: IN_HEADER;
        name: "Authorization";
        description: "Type Bearer your-jwt-token to Value";
      }
    }
  }
};

service {{.TName}} {
  // create {{.TName}}
  rpc Create(Create{{.TableName}}Request) returns (Create{{.TableName}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}"
      body: "*"
    };
  }

  // delete {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc DeleteBy{{.CrudInfo.ColumnNameCamel}}(Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {
    option (google.api.http) = {
      delete: "/api/v1/{{.TName}}/left_curly_bracket{{.CrudInfo.ColumnNameCamelFCL}}right_curly_bracket"
    };
  }

  // update {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc UpdateBy{{.CrudInfo.ColumnNameCamel}}(Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {
    option (google.api.http) = {
      put: "/api/v1/{{.TName}}/left_curly_bracket{{.CrudInfo.ColumnNameCamelFCL}}right_curly_bracket"
      body: "*"
    };
  }

  // get {{.TName}} by {{.CrudInfo.ColumnNameCamelFCL}}
  rpc GetBy{{.CrudInfo.ColumnNameCamel}}(Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request) returns (Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply) {
    option (google.api.http) = {
      get: "/api/v1/{{.TName}}/left_curly_bracket{{.CrudInfo.ColumnNameCamelFCL}}right_curly_bracket"
    };
  }

  // list of {{.TName}} by query parameters
  rpc List(List{{.TableName}}Request) returns (List{{.TableName}}Reply) {
    option (google.api.http) = {
      post: "/api/v1/{{.TName}}/list"
      body: "*"
    };
  }
}


/*
Notes for defining message fields:
    1. Suggest using camel case style naming for message field names, such as firstName, lastName, etc.
    2. If the message field name ending in 'id', it is recommended to use xxxID naming format, such as userID, orderID, etc.
    3. Add validate rules https://github.com/envoyproxy/protoc-gen-validate#constraint-rules, such as:
        uint64 id = 1 [(validate.rules).uint64.gte  = 1];

If used to generate code that supports the HTTP protocol, notes for defining message fields:
    1. If the route contains the path parameter, such as /api/v1/userExample/{id}, the defined
        message must contain the name of the path parameter and the name should be added
        with a new tag, such as int64 id = 1 [(tagger.tags) = "uri:\"id\""];
    2. If the request url is followed by a query parameter, such as /api/v1/getUserExample?name=Tom,
        a form tag must be added when defining the query parameter in the message, such as:
        string name = 1 [(tagger.tags) = "form:\"name\""].
    3. If the message field name contain underscores(such as 'field_name'), it will cause a problem
        where the JSON field names of the Swagger request parameters are different from those of the
        GRPC JSON tag names. There are two solutions: Solution 1, remove the underline from the
         message field name. Option 2, use the tool 'protoc-go-inject-tag' to modify the JSON tag name,
         such as: string first_name = 1 ; // @gotags: json:"firstName"
*/


// protoMessageCreateCode

message Create{{.TableName}}Reply {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1;
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetWebProtoValidation}};
}

message Delete{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageUpdateCode

message Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {

}

// protoMessageDetailCode

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
  {{.CrudInfo.ProtoType}} {{.CrudInfo.ColumnNameCamelFCL}} = 1 {{.CrudInfo.GetWebProtoValidation}};
}

message Get{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Reply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}Request {
  api.types.Params params = 1;
}

message List{{.TableName}}Reply {
  int64 total = 1;
  repeated {{.TableName}} {{.CrudInfo.TableNamePluralCamelFCL}} = 2;
}
`

	protoMessageCreateCommonTmpl    *template.Template
	protoMessageCreateCommonTmplRaw = `message Create{{.TableName}}Request {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.JSONName}} = {{$v.AddOne $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	protoMessageUpdateCommonTmpl    *template.Template
	protoMessageUpdateCommonTmplRaw = `message Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.JSONName}} = {{$v.AddOneWithTag2 $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	protoMessageDetailCommonTmpl    *template.Template
	protoMessageDetailCommonTmplRaw = `message {{.TableName}} {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.JSONName}} = {{$v.AddOne $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	serviceStructCommonTmpl    *template.Template
	serviceStructCommonTmplRaw = `
		{
			name: "Create",
			fn: func() (interface{}, error) {
				// todo enter parameters before testing
// serviceCreateStructCode
			},
			wantErr: false,
		},

		{
			name: "UpdateBy{{.CrudInfo.ColumnNameCamel}}",
			fn: func() (interface{}, error) {
				// todo enter parameters before testing
// serviceUpdateStructCode
			},
			wantErr: false,
		},
`

	serviceCreateStructCommonTmpl    *template.Template
	serviceCreateStructCommonTmplRaw = `				req := &serverNameExampleV1.Create{{.TableName}}Request{
					{{- range .Fields}}
						{{.Name}}:  {{.GoTypeZero}}, {{if .Comment}} // {{.Comment}}{{end}}
					{{- end}}
				}
				return cli.Create(ctx, req)`

	serviceUpdateStructCommonTmpl    *template.Template
	serviceUpdateStructCommonTmplRaw = `				req := &serverNameExampleV1.Update{{.TableName}}By{{.CrudInfo.ColumnNameCamel}}Request{
					{{- range .Fields}}
						{{.Name}}:  {{.GoTypeZero}}, {{if .Comment}} // {{.Comment}}{{end}}
					{{- end}}
				}
				return cli.UpdateBy{{.CrudInfo.ColumnNameCamel}}(ctx, req)`

	commonTmplParseOnce sync.Once
)

func initCommonTemplate() {
	commonTmplParseOnce.Do(func() {
		var err, errSum error

		handlerCreateStructCommonTmpl, err = template.New("goPostStruct").Parse(handlerCreateStructCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "handlerCreateStructCommonTmplRaw:"+err.Error())
		}
		handlerUpdateStructCommonTmpl, err = template.New("goPutStruct").Parse(handlerUpdateStructCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "handlerUpdateStructCommonTmplRaw:"+err.Error())
		}
		handlerDetailStructCommonTmpl, err = template.New("goGetStruct").Parse(handlerDetailStructCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "handlerDetailStructCommonTmplRaw:"+err.Error())
		}
		protoFileCommonTmpl, err = template.New("protoFile").Parse(protoFileCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoFileCommonTmplRaw:"+err.Error())
		}
		protoFileSimpleCommonTmpl, err = template.New("protoFileSimple").Parse(protoFileSimpleCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoFileSimpleCommonTmplRaw:"+err.Error())
		}
		protoFileForWebCommonTmpl, err = template.New("protoFileForWeb").Parse(protoFileForWebCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoFileForWebCommonTmplRaw:"+err.Error())
		}
		protoFileForSimpleWebCommonTmpl, err = template.New("protoFileForSimpleWeb").Parse(protoFileForSimpleWebCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoFileForSimpleWebCommonTmplRaw:"+err.Error())
		}
		protoMessageCreateCommonTmpl, err = template.New("protoMessageCreate").Parse(protoMessageCreateCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoMessageCreateCommonTmplRaw:"+err.Error())
		}
		protoMessageUpdateCommonTmpl, err = template.New("protoMessageUpdate").Parse(protoMessageUpdateCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoMessageUpdateCommonTmplRaw:"+err.Error())
		}
		protoMessageDetailCommonTmpl, err = template.New("protoMessageDetail").Parse(protoMessageDetailCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "protoMessageDetailCommonTmplRaw:"+err.Error())
		}
		serviceCreateStructCommonTmpl, err = template.New("serviceCreateStruct").Parse(serviceCreateStructCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "serviceCreateStructCommonTmplRaw:"+err.Error())
		}
		serviceUpdateStructCommonTmpl, err = template.New("serviceUpdateStruct").Parse(serviceUpdateStructCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "serviceUpdateStructCommonTmplRaw:"+err.Error())
		}
		serviceStructCommonTmpl, err = template.New("serviceStruct").Parse(serviceStructCommonTmplRaw)
		if err != nil {
			errSum = errors.Wrap(errSum, "serviceStructCommonTmplRaw:"+err.Error())
		}

		if errSum != nil {
			panic(errSum)
		}
	})
}
