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
[![Join the chat at https://gitter.im/zhufuyi/sponge](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/zhufuyi/sponge?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

</div>

[sponge](https://github.com/zhufuyi/sponge) is a powerful tool for generating code for web and microservice projects, and a microservice framework based on gin and grpc packaging. sponge has a rich set of code generation commands, generating a total of 12 different functional codes that can be These can be combined into a complete service (similar to a sponge cell that can be automatically reassembled into a new sponge). Microservice code features include logging, service registration and discovery, registry, flow limiting, fusing, link tracking, metrics monitoring, pprof performance analysis, statistics, caching, CICD and more. The code decoupled modular design makes it easy to build complete project code from development to deployment, making it more convenient, easy and efficient to develop projects using the go language.

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

### Online code generation demo

Online code generation demo: [https://go-sponge.com/ui](https://go-sponge.com/ui)

💡 Warning.

> Some of the generated code requires mysql account and password, do not try here to avoid the risk of being exposed.
> Server resources are limited and may not be available later. It is recommended to download the sponge binary and run the UI service to get the same generated code interface.

<br>

### Installation

- [Install sponge in linux or macOS](https://github.com/zhufuyi/sponge/blob/main/assets/install-en.md#install-sponge-in-linux-or-macos)
- [Install sponge in windows](https://github.com/zhufuyi/sponge/blob/main/assets/install-en.md#install-sponge-in-windows)

<br>

### Quick start

Once you have installed sponge and the dependencies, you are ready to go, start the ui service from the command line:

```bash
sponge run
```

Visit `http://localhost:24631` in your browser, 12 types of codes can be generated.

<br>

Example of sql based web project code generation.

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/web-http-en.gif">
</p>

After downloading the web project code, execute the command to start the service.

```bash
# update swagger
make docs

# build and run
make run
```

<br>

💡 If you don't want to use the UI interface, you can use the sponge command line to generate code. There is a wealth of examples in the command line helper information, and some of the code generation commands are more convenient than using the UI interface.

<br>

### Documentation

[sponge usage documentation](https://go-sponge.com/)

<br>

**If it's help to you, give it a star ⭐.**

<br>

## License

See the [LICENSE](LICENSE) file for licensing information.

<br>
