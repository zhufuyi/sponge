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
[![Join the chat at https://gitter.im/zhufuyi/sponge](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/zhufuyi/sponge?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

</div>

[sponge](https://github.com/zhufuyi/sponge) is a microservice framework and a tool for quickly creating web and microservice code. sponge has a rich set of generated code commands, generating a total of 12 different functional codes that can be combined into complete services (similar to a sponge that has been artificially broken up cells can be automatically reassembled into a new sponge). Microservice code functions include logging, service registration and discovery, registry, flow limiting, fusing, link tracking, metrics monitoring, pprof performance analysis, statistics, caching, CICD and more. The code is decoupled and modular in design, with commonly used code and scripts generated automatically, requiring only business logic code to be written according to the generated templates. The web and rpc services created include complete engineering code and scripts from development to deployment, making development in the go language easier and more efficient.

<br>

### sponge generates the code framework

The generated code is based on three approaches **Yaml**, **SQL DDL** and **Protobuf**, each possessing different functional code generation, and the framework diagram of the generated code is shown in Figure 1-1.

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

Online code generation demo: [https://go-sponge.com/ui](https://go-sponge.com/ui/micro-rpc-gw-pb)

💡 Warning.

> Some of the generated code requires mysql account and password, do not try here to avoid the risk of being exposed.
> Server resources are limited and may not be available later, it is recommended to build the same generation code platform locally to use.

<br>

### Installation

**(1) install sponge**

```bash
go install github.com/zhufuyi/sponge/cmd/sponge@latest
```

**(2) install dependency plugins and tools**

```bash
sponge init
```

If there is a dependency tool installation error(except `protoc`), execute the command to retry `sponge tools --install`

<br>

**(3) install protoc**

Download it from `https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3` and add the directory where the **protoc** file is located under system path.

<br>

💡 NOTICE:

> If you are using `windows` environment, you need to install some additional dependency tools, see [windows dependency tools](https://go-sponge.com/sponge-install?id=window-environment-installation-dependencies) for installation steps.

<br>

### Quick start

Once you have installed sponge and the dependencies, you are ready to go, start the ui service from the command line:

```bash
sponge run
```

Visit `http://localhost:24631` in your browser, 12 types of codes can be generated.

💡 NOTICE:

> Do not execute the `sponge run` in the directory where the sponge file is located.

<br>

### Documentation

[sponge usage documentation](https://go-sponge.com/)

<br>

**If it's useful to you, give it a star ⭐.**

<br>

## License

See the [LICENSE](LICENSE) file for licensing information.

<br>
