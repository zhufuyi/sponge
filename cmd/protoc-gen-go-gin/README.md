## protoc-gen-go-gin

根据protobuf生成gin注册路由代码，可以使用自己log和response替换默认值，除了生成注册路由代码，还支持生成http的handler模板代码和调用rpc服务端模板代码。

<br>

## 安装

### 安装依赖工具

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

### 安装 protoc-gen-go-gin

> go install github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@latest

<br>

## 使用说明

### proto 文件约定

默认情况下 rpc method 命名为`方法+资源`，使用驼峰方式命名，生成代码时会进行映射，方法映射方式如下所示:

- `"GET", "FIND", "QUERY", "LIST", "SEARCH"`  --> `GET`
- `"POST", "CREATE"`  --> `POST`
- `"PUT", "UPDATE"`  --> `PUT`
- `"DELETE"`  --> `DELETE`


使用 google.api.http option 指定路由

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

### 生成代码

只生成*_router.pb.go

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative  \
  api/v1/*.proto
```

生成*_router.pb.go和handler模板文件*_handler.go，用在由proto生成http的handler使用

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugins=handler \
  --plugin=./protoc-gen-go-gin* \
  api/v1/*.proto
```

生成*_router.pb.go和调用rpc模板文件*_service.go，用在rpc的gateway上
```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugins=service \
  --plugin=./protoc-gen-go-gin* \
  api/v1/*.proto
```
