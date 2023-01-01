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

[sponge](https://github.com/zhufuyi/sponge) is a microservice framework, a tool to quickly generate web and microservice code. sponge has a rich generating code commands, a total of 12 different functional code, these functional code can be combined into a complete service (similar to artificially broken sponge cells can be automatically reorganized into a new sponge). Microservice code features include logging, service registration and discovery, registry, rate limiter, circuit breaker, trace, metrics monitoring, pprof performance analysis, statistics, caching, CICD. Code decoupling modular design, including the complete project from development to deployment, common code and scripts are automatically generated, only in accordance with the code template to write business logic code, making the development efficiency improved a lot.

<br>

### sponge generates the code framework

The generated code is based on three approaches **Yaml**, **SQL DDL** and **Protocol buffers**, each possessing different functional code generation, and the framework diagram of the generated code is shown in Figure 1-1.

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>
<p align="center">Figure 1-1 sponge generation code framework diagram</p>

span

### Microservices framework

The microservice code framework created by sponge is shown in Figure 1-2, this is a typical microservice hierarchy with high performance, high scalability, and includes common service governance features.

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>
<p align="center">Figure 1-2 Microservices framework diagram</p>

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

> If you are using `windows` environment, you need to install some additional dependency tools, see [windows dependency tools](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/readme-en.md#21-window-environment-installation-dependencies) for installation steps.

<br>

### Quick start

Once you have installed sponge and the dependencies, you are ready to go, start the ui service from the command line:

```bash
sponge run
```

Visit `http://localhost:24631` in your browser.

💡 NOTICE:

> Do not execute the "sponge run" in the directory where the sponge file is located, as this will result in an error:
>
>> exec: "sponge": cannot run executable found relative to current directory
>>

<br>

Generate web project code from sql example.

![web-http](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/web-http.gif)

After generating the web project code, execute the command to start the service.

```bash
# Update swagger documentation
make docs

# Compile and run the service
make run
```

<br>

Generate other types of code you can try yourself.

<br>

### Documentation

[sponge usage documentation](assets/sponge-doc-en.md)

<br>

## License

See the [LICENSE](LICENSE) file for licensing information.

<br>
