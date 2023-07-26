## sponge [中文](assets/readme-cn.md)

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

[sponge](https://github.com/zhufuyi/sponge) is a powerful tool for generating web and microservice code, a microservice framework based on gin and grpc encapsulation, and an open source framework for rapid application development. Sponge has a wealth of code generation commands, and different functional codes can be combined to form a complete service (similar to artificially scattered sponge cells that can automatically recombine into a new sponge). Microservice code functions include logging, service registration and discovery, registration center, flow control, fuse, link tracking, metric monitoring, pprof performance analysis, statistics, cache, CICD and other functions. Generate code unified in the UI interface operation, it is easy to build a complete project engineering code, allowing developers to focus on the implementation of the business logic code, without spending time and energy on the project configuration and integration.

<br>

### sponge generates the code framework

The generated code is based on three approaches **Yaml**, **SQL** and **Protobuf**, each possessing different functional code generation, and the framework diagram of the generated code is shown in Figure 1-1.

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>
<p align="center">Figure 1-1 sponge generation code framework diagram</p>

<br>

UI interface for generating code:

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-ui.png">
</p>

<br>

### Components of the generated service code

The sponge separates the two major parts of code during the process of generating web service code. It isolates the business logic from the non-business logic. For example, consider the entire web service code as an egg. The eggshell represents the web service framework code, while both the albumen and yolk represent the business logic code. The yolk is the core of the business logic (manually written code). It includes defining MySQL tables, defining API interfaces, and writing specific logic code.On the other hand, the albumen acts as a bridge connecting the core business logic code to the web framework code (automatically generated, no manual writing needed). This includes the registration of route codes generated from proto files, handler method function codes, parameter validation codes, error codes, Swagger documentation, and more.

The web service egg model dissection diagram is shown in the following figure:

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/en_web-http-pb-anatomy.png">
</p>

<br>

The gRPC service egg model dissection is shown in the following figure:

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/en_micro-rpc-pb-anatomy.png">
</p>

<br>

The rpc gateway service egg model dissection is shown in the following figure:

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/en_micro-rpc-gw-pb-anatomy.png">
</p>

<br>

### Services framework

The microservice code framework created by sponge is shown in Figure 1-2, this is a typical microservice hierarchy with high performance, high scalability, and includes common service governance features.

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>
<p align="center">Figure 1-2 Microservices framework diagram</p>

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

#### Basic Services examples

- [1_web-gin-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/1_web-gin-CRUD)
- [2_web-gin-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/2_web-gin-protobuf)
- [3_micro-grpc-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/3_micro-grpc-CRUD)
- [4_micro-grpc-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [5_micro-gin-rpc-gateway](https://github.com/zhufuyi/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [6_micro-cluster](https://github.com/zhufuyi/sponge_examples/tree/main/6_micro-cluster)

#### Examples of a complete project

- [7_community-single](https://github.com/zhufuyi/sponge_examples/tree/main/7_community-single)
- [8_community-cluster](https://github.com/zhufuyi/sponge_examples/tree/main/8_community-cluster)

<br>

### Documentation

[sponge usage documentation](https://go-sponge.com/)

<br>

**If it's help to you, give it a star ⭐.**

<br>

## License

See the [LICENSE](LICENSE) file for licensing information.

<br>
