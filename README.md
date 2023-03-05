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

[sponge](https://github.com/zhufuyi/sponge) is a microservice framework and a tool for quickly creating web and microservice code. sponge has a rich set of generated code commands, generating a total of 12 different functional codes that can be combined into complete services (similar to a sponge that has been artificially broken up cells can be automatically reassembled into a new sponge). Microservice code functions include logging, service registration and discovery, registry, rate limit, circuit breaker, tracking, monitoring, pprof performance analysis, statistics, caching, CICD and more. The decoupled modular design makes it easy to build complete project code from development to deployment, making development in the go language more convenient, easy and efficient.

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

<br>

💡 If you don't want to use the UI interface, you can use the sponge command line to generate code. There is a wealth of examples in the command line helper information, and some of the code generation commands are more convenient than using the UI interface.

<br>

### Documentation

[sponge usage documentation](https://go-sponge.com/)

<br>

**If it's useful to you, give it a star ⭐.**

<br>

## License

See the [LICENSE](LICENSE) file for licensing information.

<br>
