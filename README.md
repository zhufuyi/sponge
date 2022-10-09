## sponge

<p align="center">
<img align="center" width="300px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/logo.jpg">
</p>

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/sponge)](https://goreportcard.com/report/github.com/zhufuyi/sponge)
[![codecov](https://codecov.io/gh/zhufuyi/sponge/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/sponge)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/sponge.svg)](https://pkg.go.dev/github.com/zhufuyi/sponge)
[![License: MIT](https://img.shields.io/github/license/zhufuyi/sponge)](https://img.shields.io/github/license/zhufuyi/sponge)

sponge 是一个微服务框架，支持自动生成http服务代码、grpc服务代码和CICD完整流程脚本，生成独立子模块config、ecode、model、dao、handler、router、http、proto、service、grpc，低耦合的模块代码经过组合成完整服务(类似打散的海绵细胞会自动重组新的海绵)，只需填充业务代码。

功能:

- web框架 [gin](https://github.com/gin-gonic/gin)
- rpc框架 [grpc](https://github.com/grpc/grpc-go)
- 配置文件解析 [viper](https://github.com/spf13/viper)
- 日志 [zap](go.uber.org/zap)
- 数据库组件 [gorm](gorm.io/gorm)
- 缓存组件 [go-redis](github.com/go-redis/redis)
- 文档 [swagger](github.com/swaggo/swag)
- 校验器 [validator](github.com/go-playground/validator)
- 链路跟踪 [opentelemetry](go.opentelemetry.io/otel)
- 指标采集 [prometheus](github.com/prometheus/client_golang/prometheus)
- 限流 [ratelimiter](golang.org/x/time/rate)
- 熔断 [hystrix](github.com/afex/hystrix-go)
- 配置中心 [nacos](https://github.com/alibaba/nacos)
- 服务注册与发现 [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/nacos)
- 包管理工具 [go modules](https://github.com/golang/go/wiki/Modules)
- 性能分析 [go profile](https://go.dev/blog/pprof)
- 代码检测 [golangci-lint](https://github.com/golangci/golangci-lint)
- 持续集成CI [jenkins](https://github.com/jenkinsci/jenkins)
- 持续部署CD [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

<br>

目录结构遵循[golang-standards/project-layout](https://github.com/golang-standards/project-layout)。

```
.
├── api                 # grpc的proto文件和对应代码
├── assets              # 与存储库一起使用的其他资产(图像、徽标等)
├── build               # 打包和持续集成
├── cmd                 # 应用程序的目录
├── configs             # 配置文件目录
├── deployments         # IaaS、PaaS、系统和容器编排部署配置和模板
├── docs                # 设计文档和接口文档
├── internal            # 私有应用程序和库代码
│   ├── cache           # 基于业务封装的cache
│   ├── config          # 配置文件映射的go struct
│   ├── dao             # 数据访问
│   ├── ecode           # 自定义业务错误码
│   ├── handler         # http的业务功能实现
│   ├── model           # 数据库 model
│   ├── routers         # http 路由
│   ├── server          # 服务入口，包括http和grpc服务
│   └── service         # grpc的业务功能实现
├── pkg                 # 外部应用程序可以使用的库代码
├── scripts             # 存放用于执行各种构建，安装，分析等操作的脚本
├── test                # 额外的外部测试应用程序和测试数据
└── third_party         # 外部辅助工具，分叉代码和其他第三方工具
```

<br>

开发规范遵循 [Uber Go 语言编码规范](https://github.com/uber-go/guide/blob/master/style.md) 。

<br>

## 快速开始

### 快速创建http项目

#### 创建新http服务

根据module名称、服务名称、项目名称、仓库地址和mysql表生成一个完整的http服务代码，自动实现增删改查数据功能，支持缓存、链路跟踪、指标采集、限流、性能分析等服务治理，支持构建、部署、CICD等，执行代码步骤：

**(1) 生成http服务代码**

> sponge http --module-name=account --server-name=account --project-name=account --repo-addr=zhufuyi --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=use

命令参数说明：

- --module-name: go.mod文件的module名
- --server-name: http服务名称
- --project-name: http服务所属项目名称
- --db-dsn: 连接msyql地址
- --db-table: 表名
- --embed: 是否嵌入gorm.Model(id和时间作为内嵌字段)，可选参数，默认为true
- --out: 生成代码到指定目录，例如`$PATH/src/<http服务名称>`，可选参数，默认生成代码在当前目录，注：如果文件存在会直接替换

**(2) 修改配置文件configs/<服务名称>.yml**

- 修改redis配置
- 如果字段`enableTracing`为true，必须设置jaeger地址
- 如果字段`enableRegistryDiscovery`为true，必须设置etcd地址

**(3) 生成swagger文档**

```bash
# 第一次运行，更新项目依赖库
make mod && make fmt

# 生成swagger文档
make docs
```

**(4) 启动服务**

方式一：在本地二进制运行

> make run

复制 `http://localhost:8080/swagger/index.html` 到浏览器，测试增删改查接口。

方式二：在docker运行

```bash
# 构建docker镜像
make docker-image

# 启动服务
make deploy-docker

# 查看服务运行状态，如果为healthy说明启动成功
cd deployments/docker-compose
docker-compose ps
```

方式三：在k8s运行

```bash
# 构建镜像
make image-build REPO_HOST=zhufuyi TAG=latest

# 推送镜像到远程镜像仓库，上传成功后删除本地镜像
make image-push REPO_HOST=zhufuyi TAG=latest

# 部署k8s
kubectl apply -f deployments/kubernetes/
make deploy-k8s

# 查看服务状态
kubectl get -f account-deployment.yml  
```

也可以使用Jenkins自动构建部署到k8s。

<br>

### 快速创建grpc项目

#### 创建新grpc服务

根据module名称、服务名称、项目名称、仓库地址和mysql表生成一个完整的grpc服务代码，自动实现增删改查数据功能，支持缓存、链路跟踪、指标采集、限流、熔断、性能分析等服务治理，支持构建、部署、CICD。生成一个完整的grpc服务代码步骤：

**(1) 生成grpc服务代码**

> sponge grpc --module-name=account --server-name=account --project-name=account --repo-addr=zhufuyi --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

命令参数说明：

- --project-name 项目名称
- --db-dsn 连接msyql
- --db-table 表名
- --embedded 可选参数，使用gorm.Model，id和时间作为内嵌字段
- --server-name 可选参数，服务名称，用在生成proto文件路径上，例如`serverName`，如果为空，默认等于project-name参数值
- --out: 生成代码到指定目录，例如`$PATH/src/<http服务名称>`，可选参数，默认生成代码在当前目录，注：如果文件存在会直接替换

**(2) 修改配置文件configs/<服务名称>.yml**

- 修改redis配置
- 如果字段`enableTracing`为true，必须设置jaeger地址
- 如果字段`enableRegistryDiscovery`为true，必须设置etcd地址

**(4) 生成grpc代码**

```bash
# 第一次运行，更新项目依赖库
make mod && make fmt

# 生成*.pb.go
make proto
```

**(5) 启动服务**

方式一：本地二进制运行

> make run

使用IDE打开`internal/service/表名_client_test.go`，测试增删改查接口和和生成压测报告。

方式二：在docker运行

```bash
# 构建docker镜像
make docker-image

# 启动服务
make deploy-docker

# 查看服务运行状态，如果为healthy说明启动成功
cd deployments/docker-compose
docker-compose ps
```

方式三：在k8s运行

```bash
# 构建镜像
make image-build REPO_HOST=zhufuyi TAG=latest

# 推送镜像到远程镜像仓库，上传成功后删除本地镜像
make image-push REPO_HOST=zhufuyi TAG=latest

# 部署k8s
kubectl apply -f deployments/kubernetes/
make deploy-k8s

# 查看服务状态
kubectl get -f account-deployment.yml  
```
