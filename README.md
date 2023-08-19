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

[sponge](https://github.com/zhufuyi/sponge) is a powerful golang productivity tool that integrates `automatic code generation`, `microservices framework`, `general basic development framework`. sponge has a wealth of code generation commands to generate different functionality of the code can be combined into a complete service (similar to artificially scattered sponge cells that can automatically recombine into a new sponge). Code decoupling modular design , it is easy to build a complete project from development to deployment, so that you develop web or microservices project easily, golang can also be "low-code development".

<br>

### sponge generates the code framework

sponge is mainly based on **SQL** and **Protobuf** two ways to generate code, each way has to generate code for different functions.

**Generate code framework:**

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>
<p align="center">Figure 1-1 sponge generation code framework diagram</p>

<br>

**Generate code framework corresponding UI interface:**

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/en_sponge-ui.png">
</p>

<br>

### Services framework

sponge generated microservice code framework is shown in the figure below, which is a typical microservice hierarchical structure, with high performance, high scalability, contains commonly used service governance features, you can easily replace or add their own service governance features.

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>
<p align="center">Figure 1-2 Microservices framework diagram</p>

<br>

### Egg model for complete service code

The sponge separates the two major parts of code during the process of generating web service code. It isolates the business logic from the non-business logic. For example, consider the entire web service code as an egg. The eggshell represents the web service framework code, while both the albumen and yolk represent the business logic code. The yolk is the core of the business logic (manually written code). It includes defining MySQL tables, defining API interfaces, and writing specific logic code.On the other hand, the albumen acts as a bridge connecting the core business logic code to the web framework code (automatically generated, no manual writing needed). This includes the registration of route codes generated from proto files, handler method function codes, parameter validation codes, error codes, Swagger documentation, and more.

Egg model profiling diagram for `⓷Web services created based on protobuf`:

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/en_web-http-pb-anatomy.png">
</p>

This is the egg model for web service code, and there are egg models for microservice (gRPC) code, and rpc gateway service code described in [sponge documentation](https://go-sponge.com/).

<br>

### Quick start

**Installation sponge:**

sponge can be installed on Windows, macOS, and Linux environments. Click here to view [Installation Instructions](https://github.com/zhufuyi/sponge/blob/main/assets/install-en.md).

After installing the sponge, start the UI service:

```bash
sponge run
```

Visit `http://localhost:24631` in your browser, generate code by manipulating it on the page.

<br>

### Documentation

[sponge documentation](https://go-sponge.com/)

<br>

### Examples of use

#### Simple examples (excluding business logic code)

- [1_web-gin-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/1_web-gin-CRUD)
- [2_web-gin-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/2_web-gin-protobuf)
- [3_micro-grpc-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/3_micro-grpc-CRUD)
- [4_micro-grpc-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [5_micro-gin-rpc-gateway](https://github.com/zhufuyi/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [6_micro-cluster](https://github.com/zhufuyi/sponge_examples/tree/main/6_micro-cluster)

#### Full project examples (including business logic code)

- [7_community-single](https://github.com/zhufuyi/sponge_examples/tree/main/7_community-single)
- [8_community-cluster](https://github.com/zhufuyi/sponge_examples/tree/main/8_community-cluster)

<br>

**If it's help to you, give it a star ⭐.**

<br>

### License

See the [LICENSE](LICENSE) file for licensing information.

<br>

### How to contribute

You are more than welcome to join us, raise an Issue or Pull Request.

Pull Request instructions.

1. Fork the code
2. Create your own branch: `git checkout -b feat/xxxx`
3. Commit your changes: `git commit -am 'feat: add xxxxx'`
4. Push your branch: `git push origin feat/xxxx`
5. Commit your pull request
