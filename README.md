<p align="center">
<img width="500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/logo.png">
</p>

<div align=center>

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/sponge)](https://goreportcard.com/report/github.com/zhufuyi/sponge)
[![codecov](https://codecov.io/gh/zhufuyi/sponge/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/sponge)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/sponge.svg)](https://pkg.go.dev/github.com/zhufuyi/sponge)
[![Go](https://github.com/zhufuyi/sponge/workflows/Go/badge.svg?branch=main)](https://github.com/zhufuyi/sponge/actions)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/github/license/zhufuyi/sponge)](https://img.shields.io/github/license/zhufuyi/sponge)

</div>

## 当前版本为魔改原版本，具体用法请看 [官方文档](https://github.com/zhufuyi/sponge) 

### 源码安装 使用说明
```go
    git clone https://github.com/ice-leng/sponge.git
    cd sponge/cmd/sponge
    go run ./main.go init

**Sponge** is a powerful development framework that integrates `automatic code generation`, `Gin and GRPC`. Sponge has a rich set of code generation commands, and the generated different functional codes can be combined into a complete service (similar to how artificially broken sponge cells can automatically reassemble into a new complete sponge). Sponge provides one-stop project development (code generation, development, testing, api documentation, deployment), it greatly improves development efficiency and reduces development difficulty, develop high-quality projects with a "low code approach".

<br>

### Sponge Core Design Philosophy

Sponge's core design concept is to reversely generate modular code through `SQL` or `Protobuf` files. These codes can be flexibly and seamlessly combined into various types of backend services, thus greatly improving development efficiency and simplifying backend service development. Sponge's main goals are as follows:

- If you develop a web or gRPC service with only CRUD API, you don't need to write any go code to compile and deploy it to Linux servers, dockers, k8s, and you just need to connect to the database to automatically generate the complete backend service go code by sponge.
- If you develop general web, gRPC, http+gRPC, gRPC gateway services, you only need to focus on the three core parts of `define tables in the database`, `define API description information in the protobuf file`, and `fill in business logic code in the generated template file`, and other go codes (including CRUD api) are generated by sponge.

<br>

### Sponge Generates the Code Framework

Sponge generation code is mainly based on `SQL` and `Protobuf` files, where `SQL` supports database **mysql** , **mongodb**, **postgresql**, **tidb**, **sqlite**.

#### Generate Code Framework

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>

<br>

#### Generate Code Framework Corresponding UI Interface

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/en_sponge-ui.png">
</p>

<br>

### Microservice framework

Sponge is also a microservices framework, the framework diagram is shown below, which is a typical microservice hierarchical structure, with high performance, high scalability, contains commonly used service governance features, you can easily replace or add their own service governance features.

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>

<br>

Performance testing of http and grpc service code created by the microservices framework: 50 concurrent, 1 million total requests.

![http-server](https://raw.githubusercontent.com/zhufuyi/microservices_framework_benchmark/main/test/assets/http-server.png)

![grpc-server](https://raw.githubusercontent.com/zhufuyi/microservices_framework_benchmark/main/test/assets/grpc-server.png)

Click to view the [**test code**](https://github.com/zhufuyi/microservices_framework_benchmark).

<br>

### Key Features

- Web framework [gin](https://github.com/gin-gonic/gin)
- RPC framework [grpc](https://github.com/grpc/grpc-go)
- Configuration parsing [viper](https://github.com/spf13/viper)
- Configuration center [nacos](https://github.com/alibaba/nacos)
- Logging component [zap](https://github.com/uber-go/zap)
- Database ORM component [gorm](https://github.com/go-gorm/gorm), [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
- Cache component [go-redis](https://github.com/go-redis/redis), [ristretto](https://github.com/dgraph-io/ristretto)
- Automated API documentation [swagger](https://github.com/swaggo/swag), [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2)
- Authentication [jwt](https://github.com/golang-jwt/jwt)
- Websocket [gorilla/websocket](https://github.com/gorilla/websocket)
- Crontab [cron](https://github.com/robfig/cron)
- Message Queue [rabbitmq](https://github.com/rabbitmq/amqp091-go), [kafka](https://github.com/IBM/sarama)
- Distributed Transaction Manager [dtm](https://github.com/dtm-labs/dtm)
- Distributed lock [dlock](https://github.com/zhufuyi/sponge/tree/main/pkg/dlock)
- Parameter validation [validator](https://github.com/go-playground/validator)
- Adaptive rate limiting [ratelimit](https://github.com/zhufuyi/sponge/tree/main/pkg/shield/ratelimit)
- Adaptive circuit breaking [circuitbreaker](https://github.com/zhufuyi/sponge/tree/main/pkg/shield/circuitbreaker)
- Distributed Tracing [opentelemetry](https://github.com/open-telemetry/opentelemetry-go)
- Metrics monitoring [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana)
- Service registration and discovery [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/nacos)
- Adaptive collecting [profile](https://go.dev/blog/pprof)
- Resource statistics [gopsutil](https://github.com/shirou/gopsutil)
- Code quality checking [golangci-lint](https://github.com/golangci/golangci-lint)
- Continuous integration and deployment [jenkins](https://github.com/jenkinsci/jenkins), [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

<br>

### Project Code Directory Structure

The project code directory structure created by sponge follows the [project-layout](https://github.com/golang-standards/project-layout).

Here is the directory structure for the generated `monolithic application single repository (monolith)` or `microservice multi-repository (multi-repo)` code:

```bash
.
├── api            # Protobuf files and generated * pb.go directory
├── assets         # Store various static resources, such as images, markdown files, etc.
├── cmd            # Program entry directory
├── configs        # Directory for configuration files
├── deployments    # Bare metal, docker, k8s deployment script directory.
├── docs           # Directory for API interface Swagger documentation.
├── internal       # Directory for business logic code.
│    ├── cache        # Cache directory wrapped around business logic.
│    ├── config       # Directory for Go structure configuration files.
│    ├── dao          # Data access directory.
│    ├── ecode        # Directory for system error codes and custom business error codes.
│    ├── handler      # Directory for implementing HTTP business functionality (specific to web services).
│    ├── model        # Database model directory.
│    ├── routers      # HTTP routing directory.
│    ├── rpcclient    # Directory for client-side code that connects to grpc services.
│    ├── server       # Directory for creating services, including HTTP and grpc.
│    ├── service      # Directory for implementing grpc business functionality (specific to grpc services).
│    └── types        # Directory for defining request and response parameter structures for HTTP.
├── pkg            # Directory for shared libraries.
├── scripts        # Directory for scripts.
├── test           # Directory for scripts required for testing services  and test SQL.
├── third_party    # Directory for third-party protobuf files or external helper programs.
├── Makefile       # Develop, test, deploy related command sets .
├── go.mod         # Go module dependencies and version control file.
└── go.sum         # Go module dependencies key and checksum file.
```

<br>

Here is the directory structure for the generated `microservice monolithic repository (mono-repo)` code (also known as large repository directory structure):

```bash
.
├── api
│    ├── server1       # Protobuf files and generated *pb.go directory for service 1.
│    ├── server2       # Protobuf files and generated *pb.go directory for service 2.
│    ├── server3       # Protobuf files and generated *pb.go directory for service 3.
│    └── ...
├── server1        # Code directory for Service 1, it has a similar structure to the microservice multi repo directory.
├── server2        # Code directory for Service 2, it has a similar structure to the microservice multi repo directory.
├── server3        # Code directory for Service 3, it has a similar structure to the microservice multi repo directory.
├── ...
├── third_party    # Third-party protobuf files.
├── go.mod         # Go module dependencies and version control file.
└── go.sum         # Go module dependencies' checksums and hash keys.
```

### 主要魔改功能有
- 基于数据库dsn 添加表前缀 
```html
    数据库dsn: root:@(127.0.0.1:3306)/hyperf;prefix=t_
```
- 去掉下载代码功能，替换为，命令行 在那个目录，代码就在这个目录下生成
```go
    mkdir xxx
    cd xxx
    sponge run 
    ... // web 操作 代码下载 
	ls -al
```

Access `http://localhost:24631` in a local browser and manipulate the generated code on the UI page.

> If you want to access it on a cross-host browser, you need to specify the host ip or domain name when starting the UI, example `sponge run -a http://your_host_ip:24631`. It is also possible to start the UI service on docker to support cross-host access, Click for instructions on [starting the sponge UI service in docker](https://github.com/zhufuyi/sponge/blob/main/assets/install-en.md#docker-environment).

<br>

### Sponge Development Documentation

Detailed step-by-step, configuration, deployment instructions for developing projects using sponge, Click here to view the [sponge development documentation](https://go-sponge.com/)

<br>

### Examples of use

#### Examples of create services

- [Create **web** service based on **sql** (including CRUD)](https://github.com/zhufuyi/sponge_examples/tree/main/1_web-gin-CRUD)
- [Create **grpc** service based on **sql** (including CRUD)](https://github.com/zhufuyi/sponge_examples/tree/main/2_micro-grpc-CRUD)
- [Create **web** service based on **protobuf**](https://github.com/zhufuyi/sponge_examples/tree/main/3_web-gin-protobuf)
- [Create **grpc** service based on **protobuf** ](https://github.com/zhufuyi/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [Create **grpc gateway** service based on **protobuf**](https://github.com/zhufuyi/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [Create **grpc+http** service based on **protobuf**](https://github.com/zhufuyi/sponge_examples/tree/main/_10_micro-grpc-http-protobuf)

#### Examples of develop complete project

- [Simple community web backend service](https://github.com/zhufuyi/sponge_examples/tree/main/7_community-single)
- [Simple community web service broken down into microservice](https://github.com/zhufuyi/sponge_examples/tree/main/8_community-cluster)

#### Distributed transaction examples

- [Simple distributed order system](https://github.com/zhufuyi/sponge_examples/tree/main/9_order-grpc-distributed-transaction)
- [Flash sale](https://github.com/zhufuyi/sponge_examples/tree/main/_12_sponge-dtm-flashSale)
- [E-Commerce system](https://github.com/zhufuyi/sponge_examples/tree/main/_14_eshop)

<br>
<br>

**If it's help to you, give it a star ⭐.**

<br>