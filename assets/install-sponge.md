### 2.2 Installing sponge

**(1) Installation go**

Download at https://go.dev/dl/ or https://golang.google.cn/dl/ Select version (>=1.16) to install, add `$GOROOT/bin` to the system path.

Note: If you don't have scientific internet access, it is recommended to set up a domestic proxy `go env -w GOPROXY=https://goproxy.cn,direct`

<br>

**(2) Installation of protoc**

Download it from https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3 and add the directory where the **protoc** file is located under systempath.

<br>

**(3) Installation of sponge**

> go install github.com/zhufuyi/sponge/cmd/sponge@latest

Note: The directory where the sponge binary is located must be under systempath.

<br>

**(4) Installation of dependency plug-ins and tools**

> sponge init

Dependency plugins and tools are automatically installed after executing the command: [protoc-gen-go](https://google.golang.org/protobuf/cmd/protoc-gen-go), [protoc-gen-go-grpc](https://google.golang.org/grpc/cmd/protoc-gen-go-grpc), [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate), [protoc-gen-gotag](https://github.com/srikrsna/protoc-gen-gotag), [protoc-gen-go-gin](https://github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin), [protoc-gen-go-rpc-tmpl](https://github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl), [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2), [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc), [golangci-lint](https://github.com/golangci/golangci-lint/cmd/golangci-lint), [swag](https://github.com/swaggo/swag/cmd/swag), [go-callvis](https://github.com/ofabry/go-callvis).

To view the installation of dependency tools.

```bash
# linux environment
sponge tools

# windows environment, need to specify bash.exe location
sponge tools --executor="D:\Program Files\cmder\vendor\git-for-windows\bin\bash.exe"
```

<br>

The help information for the **sponge** command has detailed usage examples, add `-h` to the end of the command to see, for example `sponge web model -h`, which is the help information returned by generating the model code for gorm based on the mysql table.

<br>
