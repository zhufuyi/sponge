[sponge](https://github.com/zhufuyi/sponge) 是一个强大的自动生成web和微服务代码工具，也是一个基于gin和grpc封装的微服务框架。sponge拥有丰富的生成代码命令，生成不同的功能代码可以组合成完整的服务(类似人为打散的海绵细胞可以自动重组成一个新的海绵)。服务代码功能包括日志、服务注册与发现、注册中心、限流、熔断、链路跟踪、指标监控、pprof性能分析、统计、缓存、CICD等功能。生成代码统一在UI界面上操作，很容易构建出一个完整的项目工程代码，让开发人员聚焦在业务逻辑代码的实现，无需花费时间和精力在项目的配置和集成上。

<br>

### 生成代码框架

sponge主要基于**SQL**和**Protobuf**两种方式生成代码，每种方式拥有生成不同功能代码，生成代码的框架图如下所示：

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>

<br>

生成代码的UI界面(同时支持命令方式生成代码)：

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-ui.png">
</p>

<br>

### 生成服务代码的组成结构

sponge生成web服务代码过程中剥离了业务逻辑与非业务逻辑两大部分代码。例如把一个完整web服务代码看作一个鸡蛋，蛋壳表示web服务框架代码，蛋白和蛋黄都表示业务逻辑代码，蛋黄是业务逻辑的核心(需要人工编写的代码)，例如定义mysql表、定义api接口、编写具体逻辑代码都属于蛋黄部分。蛋白是业务逻辑核心代码与web框架代码连接的桥梁(自动生成，不需要人工编写)，例如根据proto文件生成的注册路由代码、handler方法函数代码、参数校验代码、错误码、swagger文档等都属于蛋白部分。

web服务代码的鸡蛋模型剖析图如下图所示：

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/web-http-pb-anatomy.png">
</p>

<br>

gRPC服务代码的鸡蛋模型剖析图如下图所示：

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/micro-rpc-pb-anatomy.png">
</p>

<br>

rpc网关服务鸡蛋模型剖析图如下图所示：

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/micro-rpc-gw-pb-anatomy.png">
</p>

<br>

### 微服务框架

sponge生成的微服务代码框架如下图所示，这是典型的微服务分层结构，具有高性能，高扩展性，包含了常用的服务治理功能。

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>

<br>

### 主要功能

生成的服务代码包含的功能(按需使用)：

- Web 框架 [gin](https://github.com/gin-gonic/gin)
- RPC 框架 [grpc](https://github.com/grpc/grpc-go)
- 配置解析 [viper](https://github.com/spf13/viper)
- 配置中心 [nacos](https://github.com/alibaba/nacos)
- 日志 [zap](https://go.uber.org/zap)
- 数据库组件 [gorm](https://gorm.io/gorm)
- 缓存组件 [go-redis](https://github.com/go-redis/redis), [ristretto](https://github.com/dgraph-io/ristretto)
- 文档 [swagger](https://github.com/swaggo/swag)
- 鉴权 [jwt](https://github.com/golang-jwt/jwt)
- 校验 [validator](https://github.com/go-playground/validator)
- 限流 [ratelimit](https://github.com/zhufuyi/sponge/tree/main/pkg/shield/ratelimit)
- 熔断 [circuitbreaker](https://github.com/zhufuyi/sponge/tree/main/pkg/shield/circuitbreaker)
- 链路跟踪 [opentelemetry](https://go.opentelemetry.io/otel)
- 监控 [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana)
- 服务注册与发现 [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/)
- 性能分析 [go profile](https://go.dev/blog/pprof)
- 资源统计 [gopsutil](https://github.com/shirou/gopsutil)
- 代码规范检查 [golangci-lint](https://github.com/golangci/golangci-lint)
- 持续集成部署 CICD [jenkins](https://github.com/jenkinsci/jenkins), [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

<br>

### 目录结构

生成的服务代码目录结构遵循 [project-layout](https://github.com/golang-standards/project-layout)，代码目录结构如下所示：

```bash
.
├── api            # proto文件和生成的*pb.go目录
├── assets         # 其他与资源库一起使用的资产(图片、logo等)目录
├── build          # 打包和持续集成目录
├── cmd            # 程序入口目录
├── configs        # 配置文件的目录
├── deployments    # IaaS、PaaS、系统和容器协调部署的配置和模板目录
├── docs           # 设计文档和界面文档目录
├── internal       # 私有应用程序和库的代码目录
│    ├── cache        # 基于业务包装的缓存目录
│    ├── config       # Go结构的配置文件目录
│    ├── dao          # 数据访问目录
│    ├── ecode        # 自定义业务错误代码目录
│    ├── handler      # http的业务功能实现目录
│    ├── model        # 数据库模型目录
│    ├── routers      # http路由目录
│    ├── rpcclient    # 连接rpc服务的客户端目录
│    ├── server       # 服务入口，包括http、rpc等
│    ├── service      # rpc的业务功能实现目录
│    └── types        # http的请求和响应类型目录
├── pkg            # 外部应用程序可以使用的库目录
├── scripts        # 用于执行各种构建、安装、分析等操作的脚本目录
├── test           # 额外的外部测试程序和测试数据
└── third_party    # 外部帮助程序、分叉代码和其他第三方工具
```

<br>

### 快速安装sponge

- [在linux或macOS安装sponge](https://github.com/zhufuyi/sponge/blob/main/assets/install-cn.md#%E5%9C%A8linux%E6%88%96macos%E4%B8%8A%E5%AE%89%E8%A3%85sponge)
- [在windows安装sponge](https://github.com/zhufuyi/sponge/blob/main/assets/install-cn.md#%E5%9C%A8windows%E4%B8%8A%E5%AE%89%E8%A3%85sponge)

<br>

### 快速开始

安装完成sponge后，启动UI服务：

```bash
sponge run
```

在浏览器访问 `http://localhost:24631`，在UI页面上操作生成代码。

<br>

### 使用示例

#### 基础服务示例

- [1_web-gin-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/1_web-gin-CRUD)
- [2_web-gin-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/2_web-gin-protobuf)
- [3_micro-grpc-CRUD](https://github.com/zhufuyi/sponge_examples/tree/main/3_micro-grpc-CRUD)
- [4_micro-grpc-protobuf](https://github.com/zhufuyi/sponge_examples/tree/main/4_micro-grpc-protobuf)
- [5_micro-gin-rpc-gateway](https://github.com/zhufuyi/sponge_examples/tree/main/5_micro-gin-rpc-gateway)
- [6_micro-cluster](https://github.com/zhufuyi/sponge_examples/tree/main/6_micro-cluster)

#### 完整项目示例

- [7_community-single](https://github.com/zhufuyi/sponge_examples/tree/main/7_community-single)
- [8_community-cluster](https://github.com/zhufuyi/sponge_examples/tree/main/8_community-cluster)

<br>

### 文档

[sponge 使用文档](https://go-sponge.com/zh-cn/)

<br>

### 视频介绍

- [01 sponge的形成过程](https://www.bilibili.com/video/BV1s14y1F7Fz/)
- [02 sponge的框架介绍](https://www.bilibili.com/video/BV13u4y1F7EU/)
- [03 一键生成web服务完整项目代码](https://www.bilibili.com/video/BV1RY411k7SE/)
- [04 批量生成CRUD接口代码到web服务](https://www.bilibili.com/video/BV1AY411C7J7/)
- [05 一键生成通用的web服务项目代码](https://www.bilibili.com/video/BV1CX4y1D7xj/)
- [06 批量生成任意API接口代码到web服务](https://www.bilibili.com/video/BV1P54y1g7J9/)
- [07 一键生成rpc服务完整项目代码](https://www.bilibili.com/video/BV1Tg4y1b79U/)
- [08 批量生成CRUD代码到rpc服务](https://www.bilibili.com/video/BV1TY411z7rY/)
- [09 一键生成通用的rpc服务完整项目代码](https://www.bilibili.com/video/BV1WY4y1X7zH/)
- [10 批量生成rpc方法代码到rpc服务](https://www.bilibili.com/video/BV1Yo4y1q76o/)
- [11 rpc测试神器，简单便捷](https://www.bilibili.com/video/BV1VT411z7oj/)
- [12 一键生成rpc网关服务完整项目代码](https://www.bilibili.com/video/BV1mV4y1D7k9/)
- [13 十分钟搭建一个微服务集群示例](https://www.bilibili.com/video/BV1YM4y127YK/)
- [14 sponge实战：用chatGPT打造你的专属面试题库](https://www.bilibili.com/video/BV1V24y1w7wG/)

<br>

如果对您有帮助给个star⭐，欢迎加入**go sponge微信群交流**，加微信进群。

<img width="300px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/wechat-group.jpg">

<br>

### 如何贡献

非常欢迎您的加入，提 Issue 或 Pull Request。

Pull Request说明:

1. Fork 代码
2. 创建自己的分支: git checkout -b feat/xxxx
3. 提交你的修改: git commit -am 'feat: add xxxxx'
4. 推送您的分支: git push origin feat/xxxx
5. 提交pull request

<br><br>
