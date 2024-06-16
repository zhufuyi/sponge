## [English](../README.md) | 简体中文

<br>

[sponge](https://github.com/zhufuyi/sponge) 是一个集成了 `自动生成代码`、`Gin和GRPC` 的强大的开发框架。sponge拥有丰富的生成代码命令，生成不同的功能代码可以组合成完整的服务(类似人为打散的海绵细胞可以自动重组成一个新的海绵)。从生成代码、开发、测试、api文档到部署一站式项目开发，大幅提高了开发效率和降低了开发难度，实现"低代码方式"进行开发项目。

<br>

如果开发只有CRUD api的web或微服务，不需要编写任何go代码就可以编译并部署到linux服务器、docker、k8s上，只需要连接到数据库(mysql、mongodb、postgresql、tidb、sqlite)就可以一键自动生成完整的后端服务go代码。

如果开发通用的web或微服务，只需聚焦`在数据库定义表`、`在proto文件定义api描述信息`、`在生成的模板文件填写业务逻辑代码`三个核心部分，其他go代码都由sponge自动生成。

<br>

### 生成代码框架

sponge主要基于`SQL`和`Protobuf`两种方式生成代码，每种方式生成不同用途的代码。其中`SQL`支持数据库**mysql**、**mongodb**、**postgresql**、**tidb**、**sqlite**。

#### 生成代码的框架图

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>

<br>

#### 生成代码框架对应的UI界面

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-ui.png">
</p>

<br>

### 微服务框架

sponge也是一个微服务框架，框架图如下图所示，这是典型的微服务分层结构，具有高性能，高扩展性，包含了常用的服务治理功能，可以很方便替换或添加自己的服务治理功能。

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
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
- 配置中心 [nacos](https://github.com/alibaba/nacos)
- 日志 [zap](https://github.com/uber-go/zap)
- 数据库组件 [gorm](https://github.com/go-gorm/gorm), [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
- 缓存组件 [go-redis](https://github.com/go-redis/redis), [ristretto](https://github.com/dgraph-io/ristretto)
- 自动化api文档 [swagger](https://github.com/swaggo/swag), [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2)
- 鉴权 [jwt](https://github.com/golang-jwt/jwt)
- 校验 [validator](https://github.com/go-playground/validator)
- Websocket [gorilla/websocket](https://github.com/gorilla/websocket)
- 消息队列组件 [rabbitmq](https://github.com/rabbitmq/amqp091-go)
- 分布式事务管理器 [dtm](https://github.com/dtm-labs/dtm)
- 自适应限流 [ratelimit](https://github.com/zhufuyi/sponge/tree/main/pkg/shield/ratelimit)
- 自适应熔断 [circuitbreaker](https://github.com/zhufuyi/sponge/tree/main/pkg/shield/circuitbreaker)
- 链路跟踪 [opentelemetry](https://github.com/open-telemetry/opentelemetry-go)
- 监控 [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana)
- 服务注册与发现 [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/nacos)
- 自适应采集 [profile](https://go.dev/blog/pprof)
- 资源统计 [gopsutil](https://github.com/shirou/gopsutil)
- 代码质量检查 [golangci-lint](https://github.com/golangci/golangci-lint)
- 持续集成部署 CICD [jenkins](https://github.com/jenkinsci/jenkins), [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

<br>

### 目录结构

生成的服务代码目录结构遵循 [project-layout](https://github.com/golang-standards/project-layout)，代码目录结构如下所示。支持的仓库类型有`单体应用单体仓库(monolith)`、`微服务多仓库(multi-repo)`、`微服务单体仓库(mono-repo)`。

```bash
.
├── api            # proto文件和生成的*pb.go目录
├── assets         # 其他与资源库一起使用的资产(图片、logo等)目录
├── cmd            # 程序入口目录
├── configs        # 配置文件的目录
├── deployments    # IaaS、PaaS、系统和容器协调部署的配置和模板目录
├── docs           # 设计文档和界面文档目录
├── i(I)nternal       # 业务逻辑代码目录，如果首字母是小写(internal)，表示私有代码，如果首字母大写(Internal)表示可以被其他代码复用。
│    ├── cache        # 基于业务包装的缓存目录
│    ├── config       # Go结构的配置文件目录
│    ├── dao          # 数据访问目录
│    ├── ecode        # 自定义业务错误代码目录
│    ├── handler      # http的业务功能实现目录
│    ├── model        # 数据库模型目录
│    ├── routers      # http路由目录
│    ├── rpcclient    # 连接grpc服务的客户端目录
│    ├── server       # 服务入口，包括http、grpc等
│    ├── service      # grpc的业务功能实现目录
│    └── types        # http的请求和响应类型目录
├── pkg            # 外部应用程序可以使用的库目录
├── scripts        # 用于执行各种构建、安装、分析等操作的脚本目录
├── test           # 额外的外部测试程序和测试数据
└── third_party    # 外部帮助程序、分叉代码和其他第三方工具
```

<br>

### 快速开始

#### 安装sponge

支持在windows、mac、linux环境下安装sponge，点击查看[安装sponge说明](https://github.com/zhufuyi/sponge/blob/main/assets/install-cn.md)。

#### 启动UI服务

安装完成后，启动sponge UI服务：

```bash
sponge run
```

在本地浏览器访问 `http://localhost:24631`，在UI页面上操作生成代码。

> 如果想要在跨主机的浏览器上访问，启动UI时需要指定宿主机ip或域名，示例 `sponge run -a http://your_host_ip:24631`。 也可以在docker上启动UI服务来支持跨主机访问，点击查看[docker启动sponge UI服务说明](https://github.com/zhufuyi/sponge/blob/main/assets/install-cn.md#Docker%E7%8E%AF%E5%A2%83)。

<br>

### sponge开发文档

使用sponge开发项目的详细的操作、配置、部署说明，点击查看[sponge开发文档](https://go-sponge.com/zh-cn/)。

<br>

### 使用示例

#### 使用sponge创建服务示例

- [基于sql创建web服务(包括CRUD)](https://github.com/zhufuyi/sponge_examples/tree/main/1_web-gin-CRUD)
- [基于sql创建grpc服务(包括CRUD)](https://github.com/zhufuyi/sponge_examples/tree/main/2_micro-grpc-CRUD)
- [基于protobuf创建web服务](https://github.com/zhufuyi/sponge_examples/tree/main/3_web-gin-protobuf)
- [基于protobuf创建grpc服务](https://github.com/zhufuyi/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [基于protobuf创建grpc网关服务](https://github.com/zhufuyi/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [基于protobuf创建grpc+http服务](https://github.com/zhufuyi/sponge_examples/tree/main/a_micro-grpc-http-protobuf)

#### 使用sponge开发完整项目示例

- [简单的社区web后端服务](https://github.com/zhufuyi/sponge_examples/tree/main/7_community-single)
- [简单的社区web后端服务拆分为微服务](https://github.com/zhufuyi/sponge_examples/tree/main/8_community-cluster)

#### 分布式事务示例

- [简单的分布式订单系统](https://github.com/zhufuyi/sponge_examples/tree/main/9_order-grpc-distributed-transaction)

<br>

### 视频介绍

> 视频教程演示使用sponge v1.3.12版本，后续的版本增加了一些自动化功能、调整了一些UI界面和菜单，建议结合[文档教程](https://go-sponge.com/zh-cn/)使用。

- [01 sponge的形成过程](https://www.bilibili.com/video/BV1s14y1F7Fz/)
- [02 sponge的框架介绍](https://www.bilibili.com/video/BV13u4y1F7EU/)
- [03 一键生成完整的web服务代码](https://www.bilibili.com/video/BV1RY411k7SE/)
- [04 批量生成CRUD api代码到web服务](https://www.bilibili.com/video/BV1AY411C7J7/)
- [05 一键生成通用的web服务代码](https://www.bilibili.com/video/BV1CX4y1D7xj/)
- [06 批量生成任意api模板代码到web服务](https://www.bilibili.com/video/BV1P54y1g7J9/)
- [07 一键生成完整的grpc服务代码](https://www.bilibili.com/video/BV1Tg4y1b79U/)
- [08 批量生成CRUD api代码到grpc服务](https://www.bilibili.com/video/BV1TY411z7rY/)
- [09 一键生成通用的grpc微服务代码](https://www.bilibili.com/video/BV1WY4y1X7zH/)
- [10 批量生成grpc api代码到grpc服务](https://www.bilibili.com/video/BV1Yo4y1q76o/)
- [11 grpc测试神器，简单便捷](https://www.bilibili.com/video/BV1VT411z7oj/)
- [12 一键生成grpc网关服务代码](https://www.bilibili.com/video/BV1mV4y1D7k9/)
- [13 十分钟搭建一个微服务集群示例](https://www.bilibili.com/video/BV1YM4y127YK/)

<br>
<br>

如果对您有帮助给个star⭐，欢迎加入**go sponge微信群交流**，加微信(备注`sponge`)进群。

<img width="300px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/wechat-group.jpg">
