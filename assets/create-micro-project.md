## 4 Quick creation of microservices

### 4.1 Create rpc service from mysql

#### 4.1.1 Generating rpc service code

Using the TEACHER table in **Section 3.1.1** as an example, create the rpc service.

```bash
sponge micro rpc \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --out=./edusys
```

Check the parameter description command `sponge micro rpc -h`, which generates the rpc service code in the current edusys directory with the following directory structure.

```
.
├── api
│    ├── edusys
│    │    └── v1
│    └── types
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── cache
│    ├── config
│    ├── dao
│    ├── ecode
│    ├── model
│    ├── server
│    └── service
├── scripts
└── third_party
```

The Makefile file in the edusys directory integrates commands related to compiling, testing, running, deploying, etc. Switch to the edusys directory and execute the command to run the service:.

```bash
# Generate *pb.go
make proto

# Compile and run services
make run
```

The rpc service includes the CRUD logic code as well as the rpc client test and pressure test code, using **Goland** or **VS Code** to open the `internal/service/teacher_client_test.go` file

- For tests of methods under **Test_teacherService_methods**, fill in the test parameters before testing.
- Execute the method pressure test under **Test_teacherService_benchmark**, fill in the pressure test parameters before testing, generate the pressure test report after execution, and copy the pressure test report file path to the browser to view the statistics, as shown in Figure 4-1.

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/performance-test.jpg)
*Figure 4-1 Performance test reporting interface*

From the service startup log, you see that the default listening on port **8282** (rpc service) and port **8283** (collecting metrics or profiles) is turned on to print resource statistics per minute. In practice, some modifications are made as needed.

- To use redis as a cache, open the configuration file `configs/edusys.yml`, change the **cacheType** field value to redis, and fill in the **redis** configuration address and port.
- By default, the flow limiting, fusion, link tracking, service registration and discovery functions are off, you can open the configuration file `configs/edusys.yml` to turn on the relevant functions, if you turn on the link tracking function, you need to fill in the jaeger configuration information, if you turn on the service registration and discovery function, you need to fill in one of the consul, etcd, nacos configuration information.
- If a configuration field name is added or modified, execute the command `sponge config --server-dir=./edusys` to update the corresponding go struct; you do not need to execute the update command to modify only the field values.
- Modify the error code and error message corresponding to the CRUD method, open `ingernal/ecode/teacher_rpc.go`, modify the variable **teacherNO** value (the value is unique), the return message description is modified according to your needs, the interface error messages for the teacher table operations are added here.

<br>

#### 4.1.2 Generating service code

Two new tables course and teach were added, the structure of the data table is shown in section **3.1.3** and the service code was generated.

```bash
sponge micro service \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=course,teach \
  --out=./edusys
```

View the parameter description command `sponge micro service -h`, parameter `out` specifies the existing rpc service folder edusys, if parameter `out` is empty, you must specify the `module-name` and `server-name` parameters, generate the service code in the current directory, and then manually copy it to the folder edusys, the effect of both ways is The effect is the same.

After executing the command, the course and teach related code is generated in the following directory, and if you add custom methods or new protocol buffers files, the code is also added manually in the following directory.

```
.
├── api
│    └── edusys
│          └── v1
└── internal
      ├── cache
      ├── dao
      ├── ecode
      ├── model
      └── service
```

<br>

Switch to the edusys directory and execute the command to run the service.

```bash
# Update *.pb.go
make proto

# Compile and run services
make run
```

Open the `internal/service/course_client_test.go` and `internal/service/teach_client_test.go` files using **Goland** or **VS Code** to test the CRUD methods, you need to fill in the parameters before testing.

<br>

### 4.2 Creating rpc services from proto files

sponge not only supports creating rpc services based on mysql, but also supports generating rpc services based on proto files.

#### 4.2.1 Custom Methods

The following is a sample proto file teacher.proto Contents.

```protobuf
syntax = "proto3";

package api.edusys.v1;
option go_package = "edusys/api/edusys/v1;v1";

service teacher {
  rpc Register(RegisterRequest) returns (RegisterReply) {}
  rpc Login(LoginRequest) returns (LoginReply) {}
}

message RegisterRequest {
  string name = 1;
  string email = 2;
  string password = 3;
}

message RegisterReply {
  int64 id = 1;
}

message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginReply {
  string token = 1;
}
```

<br>

#### 4.2.2 Generating rpc service code

Open a terminal and execute the command.

```bash
sponge micro rpc-pb \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --protobuf-file=./teacher.proto \
  --out=./edusys
```

Check the parameter description command `sponge micro rpc-pb -h`, which supports \* sign matching (example `--protobuf-file=*.proto`), indicating that code is generated based on a bulk proto file, and that multiple proto files include at least one service, otherwise no code can be generated.

The generated rpc service code directory is shown below, there are some differences with the rpc service code directory generated by `sponge micro rpc`, there are no **cache**, **dao**, **model** subdirectories under the internal directory.

```
.
├── api
│    └── edusys
│          └── v1
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── config
│    ├── ecode
│    ├── server
│    └── service
├── scripts
└── third_party
```

Switch to the edusys directory and execute the command to run the service.

```bash
# Generate *pb.go file, generate service template code
make proto

# Compile and run services
make run
```

After starting the rpc service, use **Goland** or **VS Code** to open the `internal/service/teacher_client_test.go` file and test each method under **Test_teacher_methods**, fill in the test parameters before testing, you will find that the request returns an internal error, because in the template code file ` internal/service/teacher.go` (the filename teacher is the proto filename) inserts the code `panic("implement me")`, which is meant to prompt for filling in the business logic code.

<br>

#### 4.2.3 Adding new methods

Depending on the business requirements, new methods need to be added, operating in two cases.

**(1) Add new method to original proto file**

Open `api/edusys/v1/teacher.proto` and add the **bindPhone** method, for example.

Execution order.

```bash
# Generate *pb.go file, generate service template code
make proto
```

Generate the template code in the `internal/service` and `internal/ecode` directories, then copy the template code to the business logic code area at.

- The template code file with the suffix **.gen.datetime** is generated in the `internal/service` directory (example teacher.go.gen.xxxx225732), because teacher.go already exists and will not overwrite the business logic code originally written, so a new file is generated, open the file teacher.go.gen.xxxx225732, copy the template code that adds the **bindPhone** method to the teacher.go file, and then fill in the business logic code.
- The file with the suffix **teacher_rpc.go.gen.datetime** is generated in the `internal/ecode` directory, and the error code corresponding to the **bindPhone** method is copied into the teacher_rpc.go file.
- Delete all files with the suffix **.gen.datetime**.

<br>

**(2) Add new method to new proto file**

For example, if a new **course.proto** file is added, copy the **course.proto** file to the `api/edusys/v1` directory to complete the newly added interface.

Execution order.

```bash
# Generate *pb.go file, generate service template code
make proto
```

Generate code files with the **course** name prefix in the `internal/service`, `internal/ecode`, and `internal/routers` directories by doing the following two operations.

- Fill in the business code in the `internal/service/course.go` file.
- Modify the custom error code and message description in the `internal/ecode/course_rpc.go` file.

<br>

#### 4.2.4 Refining the rpc service code

The rpc service code generated by the `sponge micro rpc-pb` command does not have code related to `dao`, `cache`, `model` and other manipulation data, users can implement it themselves, if you use mysql database and redis cache, you can use **sponge** tool to generate `dao`, `cache`, `model` code directly.

Generate CRUD operation database code command.

```bash
sponge micro dao \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --include-init-db=true \
  --out=./edusys
```

Check the parameter description command `sponge micro dao -h`, the parameter `-include-init-db` is used only once in a service, remove the parameter `-include-init-db` the next time you generate `dao` code, otherwise it will cause the latest `dao` code to not be generated.

Whether you implement the `dao` code yourself or use the `dao` code generated by sponge, there are a number of operations that need to be done afterwards.

- Add mysql and redis to the initialization and release resource code of the service, open the `cmd/edusys/initial/initApp.go` file, backcomment out the call to mysql and redis initialization code, open the `cmd/edusys/initial/registerClose.go` file , backcomment out the call to mysql and redis release resource code, the initial code is a one-time change.
- The generated `dao` code, and custom methods **register** and **login** can not correspond exactly, you need to manually in the file `internal/dao/teacher.go` to supplement the code (file name teacher is the name of the table), and then in the `internal/handler/teacher. go` to fill in the business logic code (filename teacher is the name of the proto file), the business code returns the error using the error code defined in the `internal/ecode` directory, if the error message is returned directly, the requesting side will receive an UNKNOWN error message, that is, an undefined error message.
- The default uses local memory for caching, change it to use redis as cache, change the field **cacheType** value to redis in the configuration file `configs/edusys.yml`, and fill in the redis address and port.

Switching to the edusys directory to run the service again.

```bash
# Compile and run services
make run
```

After starting the rpc service, use **Goland** or **VS Code** to open the `internal/service/teacher_client_test.go` file to test each method.

<br>

### 4.3 Create rpc gateway service from proto file

Microservices usually provide fine-grained APIs, and the actual APIs provided to the client are coarse-grained APIs that require data from different microservices to be aggregated together to form an API that meets the actual requirements, which is the role of the rpc gateway. rpc gateway itself is also an http service, as shown in Figure 4-2.

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/rpc-gateway.png)
*Figure 4-2 rpc gateway framework diagram*

<br>

#### 4.3.1 Defining protocol buffers

Take e-commerce microservices as an example, the product detail page has information such as product, inventory, and product evaluation, which are stored in different microservices, and generally rarely request each microservice to get the data, direct requests to microservices will cause network pressure to multiply, and the usual practice is to aggregate multiple microservice data to return at once.

The following four folders have a simple proto file under each folder.

- **comment**: proto directory of comment services
- **inventory**: proto directory for inventory services
- **product**: proto directory of products and services
- **shopgw**: proto directory for rpc gateway services

```
.
├── comment
│    └── v1
│          └──comment.proto
├── inventory
│    └── v1
│          └── inventory.proto
├── product
│    └── v1
│          └── product.proto
└── shopgw
      └── v1
            └── shopgw.proto
```

The **comment.proto** file reads as follows.

```protobuf
syntax = "proto3";

package api.comment.v1;

option go_package = "shopgw/api/comment/v1;v1";

service Comment {
  rpc ListByProductID(ListByProductIDRequest) returns (ListByProductIDReply) {}
}

message ListByProductIDRequest {
  int64 productID = 1;
}

message CommentDetail {
  int64 id=1;
  string username = 2;
  string content = 3;
}

message ListByProductIDReply {
  int32 total = 1;
  int64 productID = 2;
  repeated CommentDetail commentDetails = 3;
}
```

<br>

The **inventory.proto** file reads as follows.

```protobuf
syntax = "proto3";

package api.inventory.v1;

option go_package = "shopgw/api/inventory/v1;v1";

service Inventory {
  rpc GetByID(GetByIDRequest) returns (GetByIDReply) {}
}

message GetByIDRequest {
  int64 id = 1;
}

message InventoryDetail {
  int64 id = 1;
  float num = 4;
  int32 soldNum =3;
}

message GetByIDReply {
  InventoryDetail inventoryDetail = 1;
}
```

<br>

The **product.proto** file reads as follows.

```protobuf
syntax = "proto3";

package api.product.v1;

option go_package = "shopgw/api/product/v1;v1";

service Product {
  rpc GetByID(GetByIDRequest) returns (GetByIDReply) {}
}

message GetByIDRequest {
  int64 id = 1;
}

message ProductDetail {
  int64 id = 1;
  string name = 2;
  float price = 3;
  string description = 4;
}

message GetByIDReply {
  ProductDetail productDetail = 1;
  int64 inventoryID = 2;
}
```

<br>

The contents of the **shopgw.proto** file are as follows. The proto for the rpc gateway service is a little different from the proto for other microservices in that you need to specify the method's routing and swagger description information.

```protobuf
syntax = "proto3";

package api.shopgw.v1;

import "api/product/v1/product.proto";
import "api/comment/v1/comment.proto";
import "api/inventory/v1/inventory.proto";
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

option go_package = "shopgw/api/shopgw/v1;v1";

// default settings for generating *.swagger.json documents
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  host: "localhost:8080"
  base_path: ""
  info: {
    title: "eshop api docs";
    version: "v0.0.0";
  };
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
};

service ShopGw {
  rpc GetDetailsByProductID(GetDetailsByProductIDRequest) returns (GetDetailsByProductIDReply) {
    option (google.api.http) = {
      get: "/api/v1/detail"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "get detail",
      description: "get detail from product id",
      tags: "shopgw",
    };
  }
}

message GetDetailsByProductIDRequest {
  int64 productID = 1;
}

message GetDetailsByProductIDReply {
  api.product.v1.ProductDetail productDetail = 1;
  api.inventory.v1.InventoryDetail inventoryDetail = 2;
  repeated api.comment.v1.CommentDetail commentDetails = 3;
}
```

<br>

#### 4.3.2 Generating rpc gateway service code

Generate the rpc gateway service code from the **shopgw.proto** file.

```bash
sponge micro rpc-gw-pb \
  --module-name=shopgw \
  --server-name=shopgw \
  --project-name=eshop \
  --repo-addr=zhufuyi \
  --protobuf-file=./shopgw/v1/shopgw.proto \
  --out=./shopgw
```

Viewing the parameter description command `sponge micro rpc-gw-pb -h`, the generated rpc gateway service code is in the current shopgw directory with the following directory structure.

```
.
├── api
│    └── shopgw
│          └── v1
├── build
├── cmd
│    └── shopgw
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── config
│    ├── ecode
│    ├── routers
│    ├── rpcclient
│    └── server
├── scripts
└── third_party
```

Since **product.proto** depends on the files **product.proto**, **inventory.proto**, **comment.proto**, copy the three dependent proto files to the api directory, the api directory structure is as follows.

```
.
├── comment
│    └── v1
│          └── comment.proto
├── inventory
│    └── v1
│          └── inventory.proto
├── product
│    └── v1
│          └── product.proto
└── shopgw
      └── v1
            └── shopgw.proto
```

<br>

Switching to the shopgw directory to run the service.

```bash
# Generate *pb.go files, generate template code, update swagger documentation
make proto

# Compile and run services
make run
```

Copy http://localhost:8080/apis/swagger/index.html to the browser to test the interface, as shown in Figure 4-3. The request returns a 500 error because the template code (internal/service/shopgw_logic.go file) calls `panic("implement me")` directly, which is meant to prompt for business logic code to be filled in.

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/rpc-gw-swag.jpg)
*Figure 4-3 swagger documentation interface for rpc gatewy*

<br>

#### 4.3.3 Refining the rpc gateway service code

**(1) Generate code to connect to the rpc server**

The service does not yet have a connection to the rpc service code, and the following is the command that generates the client code to connect to the **product**, **inventory**, and **comment** rpc services.

```bash
sponge micro rpc-cli \
  --rpc-server-name=comment,inventory,product \
  --out=./shopgw
```

View the parameter description command `sponge micro rpc-cli -h`, the parameter `out` specifies the existing service folder shopgw, and the generated code is in the `internal/rpcclent` directory.

<br>

**(2) Initializing and closing rpc connections**

The connection rpc server code includes initialization and shutdown functions, which are filled in according to the calling template code at

- Initialize when starting the service, under the code segment `// initializing the rpc server connection` in the `cmd/shopgw/initial/initApp.go` file, calling the initialization function based on the template.
- To release resources when closing the service, the release function is called according to the template in the `cmd/shopgw/initial/registerClose.go` file under the code segment `// close the rpc client connection`.

<br>

**(3) Modification of configuration**

Connection **product**, **inventory**, **comment** three rpc service code has been available, but the rpc service address has not been configured, you need to add in the configuration file `configs/shopgw.yml` under the field `grpcClient` to connect product, inventory , comment three microservice configuration information.

```yaml
grpcClient:
  - name: "product"
    host: "127.0.0.1"
    port: 8201
    registryDiscoveryType: ""
  - name: "inventory"
    host: "127.0.0.1"
    port: 8202
    registryDiscoveryType: ""  
  - name: "comment"
    host: "127.0.0.1"
    port: 8203
    registryDiscoveryType: ""
```

If the rpc service uses registration and discovery, the field `registryDiscoveryType` fills in the service registration and discovery type, which supports consul, etcd, and nacos.

Generate the corresponding go struct code.

```bash
sponge config --server-dir=./shopgw
```

<br>

**(4) Fill in the operational code**

The following is a sample business logic code filled in the template file `internal/service/shopgw_logic.go` to fetch data from **product**, **inventory**, **comment** three rpc services respectively aggregated together and returned.

```go
package service

import (
	"context"

	commentV1 "shopgw/api/comment/v1"
	inventoryV1 "shopgw/api/inventory/v1"
	productV1 "shopgw/api/product/v1"
	shopgwV1 "shopgw/api/shopgw/v1"
	"shopgw/internal/rpcclient"
)

var _ shopgwV1.ShopGwLogicer = (*shopGwClient)(nil)

type shopGwClient struct {
	productCli productV1.ProductClient
	inventoryCli inventoryV1.InventoryClient
	commentCli commentV1.CommentClient
}

// NewShopGwClient creating rpc clients
func NewShopGwClient() shopgwV1.ShopGwLogicer {
	return &shopGwClient{
		productCli: productV1.NewProductClient(rpcclient.GetProductRPCConn()),
		inventoryCli: inventoryV1.NewInventoryClient(rpcclient.GetInventoryRPCConn()),
		commentCli: commentV1.NewCommentClient(rpcclient.GetCommentRPCConn()),
	}
}

func (c *shopGwClient) GetDetailsByProductID(ctx context.Context, req *shopgwV1.GetDetailsByProductIDRequest) (*shopgwV1.GetDetailsByProductIDReply, error) {
	productRep, err := c.productCli.GetByID(ctx, &productV1.GetByIDRequest{
		Id: req.ProductID,
	})
	if err ! = nil {
		return nil, err
	}

	inventoryRep, err := c.inventoryCli.GetByID(ctx, &inventoryV1.GetByIDRequest{
		Id: productRep.InventoryID,
	})
	if err ! = nil {
		return nil, err
	}

	commentRep, err := c.commentCli.ListByProductID(ctx, &commentV1.ListByProductIDRequest{
		ProductID: req,
	})
	if err ! = nil {
		return nil, err
	}

	return &shopgwV1.GetDetailsByProductIDReply{
		ProductDetail: productRep.ProductDetail,
		InventoryDetail: inventoryRep.InventoryDetail,
		CommentDetails: commentRep,
	}, nil
}
```

Start the service again.

```bash
# Compile and run services
make run
```

When visiting http://localhost:8080/apis/swagger/index.html in a browser, the request returns a 503 error (service unavailable) because none of the three rpc services **product**, **inventory**, and **comment** are running yet.

The code for all three rpc services **product**, **inventory** and **comment** are not available yet, so how to start them properly. The proto files for these three rpc services are already available, and it is easy to generate the code and start the services according to the section **4.2 Creating rpc services from proto files** steps.

<br>

### 4.4 Summary

The generation of rpc service code is based on both mysql and proto files, according to the proto file method in addition to support the generation of rpc service code, also support the generation of rpc gateway service (http) code.

- The rpc service code generated according to mysql includes CRUD method logic code and proto code for each data table, subsequently if you want to add new methods, just define them in the proto file, manually add business logic code can refer to CRUD logic code.
- Generate rpc service code based on proto file does not include operational database code, but you can use `sponge web dao` command to generate operational database code, generate service template code based on proto file, and populate business logic code in the template code.
- Generate rpc gateway service code based on proto file, interface definition in proto file, generate service template code based on proto file, populate business logic code in template code, use in combination with `sponge micro rpc-cli` command.

According to the actual scenario choose to generate the corresponding service code, if the main is to add, delete and check the data table, according to mysql generate rpc service can write less code; if more custom methods, according to the proto generate rpc service is more appropriate; rpc to http use rpc gateway service.

<br><br>