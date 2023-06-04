## sponge [中文文档](assets/readme-cn.md)

<p align="center">
<img width="500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/logo.png">
</p>

<div align=center>

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/sponge)](https://goreportcard.com/report/github.com/zhufuyi/sponge)
[![codecov](https://codecov.io/gh/zhufuyi/sponge/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/sponge)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/sponge.svg)](https://pkg.go.dev/github.com/zhufuyi/sponge)
[![Go](https://github.com/zhufuyi/sponge/workflows/Go/badge.svg?branch=main)](https://github.com/zhufuyi/sponge/actions)
[![License: MIT](https://img.shields.io/github/license/zhufuyi/sponge)](https://img.shields.io/github/license/zhufuyi/sponge)

</div>

[sponge](https://github.com/zhufuyi/sponge) is a powerful tool for generating web and microservice code, as well as a microservice framework based on gin and grpc encapsulation. Sponge has a wealth of code generation commands, and different functional codes can be combined to form a complete service (similar to artificially scattered sponge cells that can automatically recombine into a new sponge). Microservice code functions include logging, service registration and discovery, registration center, flow control, fuse, link tracking, metric monitoring, pprof performance analysis, statistics, cache, CICD and other functions. The code is decoupled and modularly designed, making it easy to build complete engineering code from development to deployment, making it more convenient, easy and efficient to develop with Go language.

<br>

### sponge generates the code framework

The generated code is based on three approaches **Yaml**, **SQL** and **Protobuf**, each possessing different functional code generation, and the framework diagram of the generated code is shown in Figure 1-1.

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>
<p align="center">Figure 1-1 sponge generation code framework diagram</p>

<br>

### Microservices framework

The microservice code framework created by sponge is shown in Figure 1-2, this is a typical microservice hierarchy with high performance, high scalability, and includes common service governance features.

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>
<p align="center">Figure 1-2 Microservices framework diagram</p>

<br>

### Supported code types for generation

- [Generate complete web service project code based on MySQL table](https://go-sponge.com/en/ui/web-http)
- [Generate handler code based on MySQL tables, including CRUD code](https://go-sponge.com/en/ui/web-handler)
- [Generate dao code based on MySQL tables, including CRUD code](https://go-sponge.com/en/ui/web-dao)
- [Generate gorm model code based on MySQL tables](https://go-sponge.com/en/ui/web-model)
- [Generate complete rpc service project code based on MySQL table](https://go-sponge.com/en/ui/micro-rpc)
- [Generate service code based on MySQL tables, including CRUD code](https://go-sponge.com/en/ui/micro-service)
- [Generate protobuf code based on MySQL tables](https://go-sponge.com/en/ui/micro-protobuf)
- [Generate generic web service project code based on protobuf](https://go-sponge.com/en/ui/web-http-pb)
- [Generate generic rpc service project code based on protobuf](https://go-sponge.com/en/ui/micro-rpc-pb)
- [Generate rpc gateway service project code based on protobuf](https://go-sponge.com/en/ui/micro-rpc-gw-pb)
- [Generate corresponding go structure code based on yaml](https://go-sponge.com/en/ui/yaml-config)
- [Generate rpc connection code according to parameters](https://go-sponge.com/en/ui/micro-rpc-conn)
- [Generate cache code according to parameters](https://go-sponge.com/en/ui/web-cache)

The generated code can be combined into actual project web or microservice code, and the developer only needs to focus on writing the business logic code, online UI interface demo: [https://go-sponge.com/ui](https://go-sponge.com/ui)

<br>

### Installation

- [Install sponge in linux or macOS](https://github.com/zhufuyi/sponge/blob/main/assets/install-en.md#install-sponge-in-linux-or-macos)
- [Install sponge in windows](https://github.com/zhufuyi/sponge/blob/main/assets/install-en.md#install-sponge-in-windows)

<br>

### Quick start

After installing the sponge, start the UI service:

```bash
sponge run
```

Visit `http://localhost:24631` in your browser, generate code by manipulating it on the page.

<br>

### Examples of use

- [1_web-gin-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/1_web-gin-CRUD)
- [2_web-gin-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/2_web-gin-protobuf)
- [3_micro-grpc-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/3_micro-grpc-CRUD)
- [4_micro-grpc-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [5_micro-gin-rpc-gateway](https://github.com/zhufuyi/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [6_micro-cluster](https://github.com/zhufuyi/sponge_examples/tree/main/6_micro-cluster)

<br>

### Documentation

[sponge usage documentation](https://go-sponge.com/)

<br>

**If it's help to you, give it a star ⭐.**

<br>

## License

See the [LICENSE](LICENSE) file for licensing information.

<br>
