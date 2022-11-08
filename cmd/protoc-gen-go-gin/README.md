## protoc-gen-go-gin

According to protobuf to generate gin registration route code, you can use their own log and response to replace the default value, in addition to generating registration route code, but also support the generation of http handler template code and call rpc server template code.

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
service GreeterService {
  rpc Create(CreateDemoRequest) returns (CreateDemoReply) {
    option (google.api.http) = {
      post: "/api/v1/demo"
      body: "*"
    };
  }

  rpc GetByID(GetDemoByIDRequest) returns (GetDemoByIDReply) {
    option (google.api.http) = {
      get: "/api/v1/demo/{id}"
    };
  }
}
```

<br>

#### Generate code

(1) Generate only *_router.pb.go

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative  \
  api/v1/*.proto
```

<br>

(2) Generate *_router.pb.go and handler template file *_logic.go

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=handler \
  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/handler \
  api/v1/*.proto
```

<br>

(3) Generate *_router.pb.go and call the rpc template file *_logic.go

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=service \
  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/service \
  api/v1/*.proto
```
