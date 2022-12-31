## 2 Install sponge and dependency tools

### 2.1 Window environment installation dependencies

If you use the windows environment, you need to install the relevant dependencies first, and just ignore the other environments.

**(1) Installing mingw64**

mingw64 stands for Minimalist GNUfor Windows, a freely available and freely distributed collection of Windows-specific header files and import libraries using the GNU toolset, download the pre-compiled source generated binaries at

https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z

After downloading and extracting to the `D:\Program Files\mingw64` directory, modify the system environment variable PATH to add `D:\Program Files\mingw64\bin`.

**Install the make command**

Switch to the `D:\Program Files\mingw64\bin` directory, find the `mingw32-make.exe` executable, copy it and rename it to `make.exe`.

<br>

**(2) Installation cmder**

**cmder** is an enhanced command line tool that contains some sponge dependent commands (bash, git, etc.), cmder download at

https://github.com/cmderdev/cmder/releases/download/v1.3.20/cmder.zip

After downloading and extracting to the `D:\Program Files\cmder` directory, modify the system environment variable PATH to add `D:\Program Files\cmder`.

<br>

### 2.2 Install sponge

**(1) Install go**

Download at https://go.dev/dl/ or https://golang.google.cn/dl/ Select version (>=1.16) to install, add `$GOROOT/bin` to the system path.

Note: If you don't have scientific internet access, it is recommended to set up a domestic proxy `go env -w GOPROXY=https://goproxy.cn,direct`

<br>

**(2) Install protoc**

Download it from https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3 and add the directory where the **protoc** file is located under systempath.

<br>

**(3) Install sponge**

> go install github.com/zhufuyi/sponge/cmd/sponge@latest

Note: The directory where the sponge binary is located must be under system path.

<br>

**(4) Install plug-ins and tools**

> sponge init

Dependency plugins and tools are automatically installed after executing the command: [protoc-gen-go](https://google.golang.org/protobuf/cmd/protoc-gen-go), [protoc-gen-go-grpc](https://google.golang.org/grpc/cmd/protoc-gen-go-grpc), [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate), [protoc-gen-gotag](https://github.com/srikrsna/protoc-gen-gotag), [protoc-gen-go-gin](https://github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin), [protoc-gen-go-rpc-tmpl](https://github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl), [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2), [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc), [golangci-lint](https://github.com/golangci/golangci-lint/cmd/golangci-lint), [swag](https://github.com/swaggo/swag/cmd/swag), [go-callvis](https://github.com/ofabry/go-callvis).

If there is a dependency tool installation error, execute the command to retry.

> sponge tools --install

To view the installation of dependency tools.

```bash
# linux environment
sponge tools
```

List of all dependent tools.

```
Installed dependency tools:
    ✔  go
    ✔  protoc
    ✔  protoc-gen-go
    ✔  protoc-gen-go-grpc
    ✔  protoc-gen-validate
    ✔  protoc-gen-gotag
    ✔  protoc-gen-go-gin
    ✔  protoc-gen-go-rpc-tmpl
    ✔  protoc-gen-openapiv2
    ✔  protoc-gen-doc
    ✔  swag
    ✔  golangci-lint
    ✔  go-callvis
```

<br>

The help information for the **sponge** command has detailed usage examples, add `-h` to the end of the command to see, for example `sponge web model -h`, which is the help information returned by generating the model code for gorm based on the mysql table.

<br>

### 2.3 Command UI for sponge

After installing sponge and the dependency tools, you can start using it. sponge provides two ways to generate code, command line and UI interface, in fact, through the UI interface way, in the background is also to execute commands. sponge UI interface supports the memory function of executing commands, it is more convenient to use, start UI service: ``bash

```bash
sponge run
```

Visit `http://localhost:24631` in the browser, open the homepage as shown in Figure 2-1, generate the code according to the actual need, the code is api interface downloaded to the local, so you can also deploy sponge's UI service to the server for long-term operation.

![home](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/home.png)

*Figure2-1 sponge ui diagram*
