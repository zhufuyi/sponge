## protoc-gen-go-gin

根据protobuf生成gin调用rpc的handler方法。

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

> go install github.com/zhufuyi/gotool/tools/protoc-gen-go-gin@latest

注：go版本必须大于1.16

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

### 生成*_gin.pb.go代码

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-gin_out=. --go-gin_opt=paths=source_relative  \
  api/v1/*.proto
```

