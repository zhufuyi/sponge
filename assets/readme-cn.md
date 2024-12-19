## [English](../README.md) | 简体中文

<br>

**sponge** 是一个强大的 `Go` 开发框架，其核心理念是通过解析 `JSON`、`SQL` 或 `Protobuf` 文件逆向生成模块化的代码，这些代码可以灵活、无缝地组合成多种类型的完整后端服务(类似海绵细胞的特性，打散的海绵细胞能自动重新组合成新的海绵)。sponge 提供一站式项目开发解决方案，涵盖代码生成、开发、测试、API 文档生成和部署，大幅提升开发效率，降低开发难度，实现以"低代码"方式构建高质量项目。

<br>

### 适用场景

sponge 适用于快速开发各种高性能后端服务，包括但不限于：
- `Web` 服务；
- `gRPC` 服务；
- `HTTP+gRPC` 混合服务；
- `gRPC Gateway API` 服务。

此外，开发者还可以通过自定义模板，生成满足业务需求的各类代码。

<br>

### 核心亮点

1. **一键生成完整后端服务代码**  
   对于仅需 `CRUD API` 的 `Web` 或 `gRPC` 服务，无需编写任何 `Go` 代码。只需连接数据库(如 `MySQL`、`MongoDB`、`PostgreSQL`、`SQLite`)，即可一键生成完整后端服务代码，并轻松部署到 Linux 服务器、Docker 或 Kubernetes 上。

2. **高效开发通用服务**  
   开发通用的 `Web`、`gRPC`、`HTTP+gRPC` 或 `gRPC Gateway` 服务，只需专注于以下三部分：
    - 数据库表的定义；
    - 在 Protobuf 文件中定义 API 描述信息；
    - 在生成的模板中编写业务逻辑代码。  

   服务的框架代码和 CRUD API 代码均由 sponge 自动生成。

3. **支持自定义模板，灵活扩展**  
   sponge 支持通过自定义模板生成项目所需的多种代码类型，不局限于 `Go` 语言。例如：
    - 后端代码；
    - 前端代码；
    - 配置文件；
    - 测试代码；
    - 构建和部署脚本等。

<br>

### 快速开始

1. **安装 sponge**

   支持在 windows、mac、linux 环境下安装 sponge，点击查看 [**安装 sponge 说明**](https://github.com/go-dev-frame/sponge/blob/main/assets/install-cn.md)。

2. **打开生成代码 UI 页面**

   安装完成后，执行命令打开 sponge UI 页面：

   ```bash
   sponge run
   ```

   在本地浏览器访问 `http://localhost:24631`，在页面上操作生成代码，如下图所示：

   <p align="center">
   <img width="1500px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/sponge-ui.png">
   </p>

   > 如果想要在跨主机的浏览器上访问，启动UI时需要指定宿主机ip或域名，示例 `sponge run -a http://your_host_ip:24631`。 也可以在 docker 上启动UI服务来支持跨主机访问，点击查看 [在 docker 运行 sponge UI 服务说明](https://github.com/go-dev-frame/sponge/blob/main/assets/install-cn.md#Docker%E7%8E%AF%E5%A2%83)。

<br>

### 生成代码的框架图

sponge 支持基于自带模板和自定义模板两种方式生成你的项目所需的代码，下面是两种生成代码的框架图。

1. sponge 基于自带模板生成代码框架如下图所示，共有 sql 和 protobuf 两种方式生成代码。

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/sponge-framework.png">
</p>

<br>

2. sponge 基于自定义模板生成代码框架如下图所示，共有 json、sql、protobuf 三种方式生成代码。

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/template-framework.png">
</p>

<br>

### 微服务框架

sponge 支持创建 6 种类型的后端服务，均为微服务架构。下图展示了典型的微服务分层结构，具备高性能、高扩展性，并内置常用的服务治理功能。

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/microservices-framework.png">
</p>

<br>

创建的http和grpc服务代码的性能测试： 50个并发，总共100万个请求。

![http-server](https://raw.githubusercontent.com/zhufuyi/microservices_framework_benchmark/main/test/assets/http-server.png)

![grpc-server](https://raw.githubusercontent.com/zhufuyi/microservices_framework_benchmark/main/test/assets/grpc-server.png)

点击查看[**测试代码**](https://github.com/zhufuyi/microservices_framework_benchmark)。

<br>

### 主要功能

sponge包含丰富的组件(按需使用)：

- Web 框架 [gin](https://github.com/gin-gonic/gin)
- RPC 框架 [grpc](https://github.com/grpc/grpc-go)
- 配置解析 [viper](https://github.com/spf13/viper)
- 日志 [zap](https://github.com/uber-go/zap)
- 数据库组件 [gorm](https://github.com/go-gorm/gorm), [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
- 缓存组件 [go-redis](https://github.com/go-redis/redis), [ristretto](https://github.com/dgraph-io/ristretto)
- 自动化api文档 [swagger](https://github.com/swaggo/swag), [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2)
- 鉴权 [jwt](https://github.com/golang-jwt/jwt)
- 校验 [validator](https://github.com/go-playground/validator)
- Websocket [gorilla/websocket](https://github.com/gorilla/websocket)
- 定时任务 [cron](https://github.com/robfig/cron)
- 消息队列组件 [rabbitmq](https://github.com/rabbitmq/amqp091-go), [kafka](https://github.com/IBM/sarama)
- 分布式事务管理器 [dtm](https://github.com/dtm-labs/dtm)
- 分布式锁 [dlock](https://github.com/go-dev-frame/sponge/tree/main/pkg/dlock)
- 自适应限流 [ratelimit](https://github.com/go-dev-frame/sponge/tree/main/pkg/shield/ratelimit)
- 自适应熔断 [circuitbreaker](https://github.com/go-dev-frame/sponge/tree/main/pkg/shield/circuitbreaker)
- 链路跟踪 [opentelemetry](https://github.com/open-telemetry/opentelemetry-go)
- 监控 [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana)
- 服务注册与发现 [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/nacos)
- 自适应采集 [profile](https://go.dev/blog/pprof)
- 资源统计 [gopsutil](https://github.com/shirou/gopsutil)
- 配置中心 [nacos](https://github.com/alibaba/nacos)
- 代码质量检查 [golangci-lint](https://github.com/golangci/golangci-lint)
- 持续集成部署 CICD [jenkins](https://github.com/jenkinsci/jenkins), [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)
- 生成项目业务架构图 [spograph](https://github.com/go-dev-frame/spograph)
- 自定义模板生成代码 [go template](https://pkg.go.dev/text/template@go1.23.3)

<br>

### 目录结构

生成的服务代码目录结构遵循 [project-layout](https://github.com/golang-standards/project-layout)。

这是生成的`单体应用单体仓库(monolith)`或`微服务多仓库(multi-repo)`代码目录结构：

```bash
.
├── api            # protobuf文件和生成的*pb.go目录
├── assets         # 其他与资源库一起使用的资产(图片、logo等)目录
├── cmd            # 程序入口目录
├── configs        # 配置文件的目录
├── deployments    # 裸机、docker、k8s部署脚本目录
├── docs           # 设计文档和界面文档目录
├── internal       # 业务逻辑代码目录
│    ├── cache        # 基于业务包装的缓存目录
│    ├── config       # Go结构的配置文件目录
│    ├── dao          # 数据访问目录
│    ├── database     # 数据库目录
│    ├── ecode        # 自定义业务错误代码目录
│    ├── handler      # http的业务功能实现目录
│    ├── model        # 数据库模型目录
│    ├── routers      # http路由目录
│    ├── rpcclient    # 连接grpc服务的客户端目录
│    ├── server       # 服务入口，包括http、grpc等
│    ├── service      # grpc的业务功能实现目录
│    └── types        # http的请求和响应类型目录
├── pkg            # 外部应用程序可以使用的库目录
├── scripts        # 执行脚本目录
├── test           # 额外的外部测试程序和测试数据
├── third_party    # 依赖第三方protobuf文件或其他工具的目录
├── Makefile       # 开发、测试、部署相关的命令集合
├── go.mod         # go 模块依赖关系和版本控制文件
└── go.sum         # go 模块依赖项的密钥和校验文件
```

<br>

这是生成的`微服务单体仓库(mono-repo)`代码目录结构(也就是大仓库代码目录结构)：

```bash
.
├── api
│    ├── server1       # 服务1的protobuf文件和生成的*pb.go目录
│    ├── server2       # 服务2的protobuf文件和生成的*pb.go目录
│    ├── server3       # 服务3的protobuf文件和生成的*pb.go目录
│    └── ...
├── server1        # 服务1的代码目录，与微服务多仓库(multi-repo)目录结构基本一样
├── server2        # 服务2的代码目录，与微服务多仓库(multi-repo)目录结构基本一样
├── server3        # 服务3的代码目录，与微服务多仓库(multi-repo)目录结构基本一样
├── ...
├── third_party    # 依赖的第三方protobuf文件
├── go.mod         # go 模块依赖关系和版本控制文件
└── go.sum         # go 模块依赖项的密钥和校验和文件
```

<br>

### sponge 开发文档

点击查看 [sponge 开发项目的详细文档](https://go-sponge.com/zh-cn/)，包括代码生成、开发、配置、部署说明等。

<br>

### 使用示例

#### 使用 sponge 创建服务示例

- [基于sql创建web服务(包括CRUD)](https://github.com/go-dev-frame/sponge_examples/tree/main/1_web-gin-CRUD)
- [基于sql创建grpc服务(包括CRUD)](https://github.com/go-dev-frame/sponge_examples/tree/main/2_micro-grpc-CRUD)
- [基于protobuf创建web服务](https://github.com/go-dev-frame/sponge_examples/tree/main/3_web-gin-protobuf)
- [基于protobuf创建grpc服务](https://github.com/go-dev-frame/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [基于protobuf创建grpc网关服务](https://github.com/go-dev-frame/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [基于protobuf创建grpc+http服务](https://github.com/go-dev-frame/sponge_examples/tree/main/_10_micro-grpc-http-protobuf)

#### 使用 sponge 开发完整项目示例

- [简单的社区web后端服务](https://github.com/go-dev-frame/sponge_examples/tree/main/7_community-single)
- [简单的社区web后端服务拆分为微服务](https://github.com/go-dev-frame/sponge_examples/tree/main/8_community-cluster)

#### 分布式事务示例

- [简单的分布式订单系统](https://github.com/go-dev-frame/sponge_examples/tree/main/9_order-grpc-distributed-transaction)
- [秒杀抢购活动](https://github.com/go-dev-frame/sponge_examples/tree/main/_12_sponge-dtm-flashSale)
- [电商系统](https://github.com/go-dev-frame/sponge_examples/tree/main/_14_eshop)

<br>
<br>

如果对您有帮助给个star⭐，欢迎加入**go sponge微信群交流**，加微信(备注`sponge`)进群。

<img width="300px" src="https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/wechat-group.jpg">
