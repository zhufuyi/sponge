## protoc-gen-go-rpc-tmpl

According to protobuf to generate rpc server template codes and rpc error code codes.

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

#### Install protoc-gen-go-rpc-tmpl

> go install github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@latest

<br>

### Usage

```bash
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --go-rpc-tmpl_out=. --go-rpc-tmpl_opt=paths=source_relative \
  --go-rpc-tmpl_opt=moduleName=yourModuleName --go-rpc-tmpl_opt=serverName=yourServerName \
  api/v1/*.proto
```

A total of 2 files are generated: the rpc service template file *.go (default save directory is internal/service),  the rpc error code file *_rpc.go (default save directory is internal/ecode).
