## protoc-gen-go-gin

According to protobuf to generate gin registration route codes, you can use their own log and response to replace the default value, in addition to generating registration route code, but also support the generation of http handler template code, call rpc server template code, http or rpc error code codes.

<br>

### Installation

#### Installation of dependency tools

```bash
# install protoc in linux
mkdir -p protocDir \
  && curl -L -o protocDir/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.20.1/protoc-3.20.1-linux-x86_64.zip \
  && unzip protocDir/protoc.zip -d protocDir\
  && mv protocDir/bin/protoc protocDir/include/ $GOROOT/bin/ \
  && rm -rf protocDir

# install protoc-gen-go, protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

#### Install protoc-gen-go-gin

> go install github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@latest

<br>

### Usage

#### protobuf documentation conventions

By default the rpc method is named `method+resource`, using camel naming, which is mapped when generating code, as shown below:

- `"GET", "FIND", "QUERY", "LIST", "SEARCH"`  --> `GET`
- `"POST", "CREATE"`  --> `POST`
- `"PUT", "UPDATE"`  --> `PUT`
- `"DELETE"`  --> `DELETE`

Specify the route using the `google.api.http` option

```protobuf
service Greeter {
  rpc Create(CreateRequest) returns (CreateReply) {
    option (google.api.http) = {
      post: "/api/v1/greeter"
      body: "*"
    };
  }

  rpc GetByID(GetByIDRequest) returns (GetByIDReply) {
    option (google.api.http) = {
      get: "/api/v1/greeter/{id}"
    };
  }
}
```

<br>

#### Generate code

(1) Generate only *_router.go

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative  \
  api/v1/*.proto
```

<br>

(2) Generate codes with plugin handler

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=handler \
  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName \
  api/v1/*.proto
```

A total of 4 files are generated: the registration route file _*router.pb.go, the injection route file *_router.go (default save path in internal/routers), the logic code template file *.go (default save path in internal/handler), the error code file *_http.go (default save path in internal/ecode).

<br>

(3) Generate codes with plugin service

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=service \
  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName \
  api/v1/*.proto
```

A total of 4 files are generated: the registration route file *_router.pb.go, the injection route file *_router.go (default save path in internal/routers), and the logic code template file *.go ( default save path in internal/service), the error code file *_rpc.go (default save path in internal/ecode).
