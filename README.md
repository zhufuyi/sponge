## sponge

sponge 是一个微服务框架，支持http和grpc及服务治理，结合[goctl](https://github.com/zhufuyi/goctl)工具自动生成框架代码。

功能:

- web框架 [gin](https://github.com/gin-gonic/gin)
- rpc框架 [grpc](https://github.com/grpc/grpc-go)
- 配置文件解析 [viper](https://github.com/spf13/viper)
- 日志 [zap](go.uber.org/zap)
- 数据库组件 [gorm](gorm.io/gorm)
- 缓存组件 [go-redis](github.com/go-redis/redis)
- 生成文档 [swagger](github.com/swaggo/swag)
- 校验器 [validator](github.com/go-playground/validator)
- 链路跟踪 [opentelemetry](go.opentelemetry.io/otel)
- 指标采集 [prometheus](github.com/prometheus/client_golang/prometheus)
- 限流 [ratelimiter](golang.org/x/time/rate)
- 熔断 [hystrix](github.com/afex/hystrix-go)
- 包管理工具 [go modules](https://github.com/golang/go/wiki/Modules)
- 性能分析 [go profile](https://go.dev/blog/pprof)
- 代码检测 [golangci-lint](https://github.com/golangci/golangci-lint)

<br>

### 目录结构

目录结构遵循[golang-standards/project-layout](https://github.com/golang-standards/project-layout)。

```
├── cmd                 # 应用程序的目录
├── config              # 配置文件目录
├── docs                # 设计和用户文档
├── internal            # 私有应用程序和库代码
│   ├── cache           # 基于业务封装的cache
│   ├── dao             # 数据访问
│   ├── ecode           # 自定义业务错误码
│   ├── handler         # http的业务功能实现
│   ├── model           # 数据库 model
│   ├── routers         # http 路由
│   ├── server          # 服务入口，包括http和grpc服务
│   └── service         # grpc的业务功能实现
├── pkg                 # 外部应用程序可以使用的库代码
├── scripts             # 存放用于执行各种构建，安装，分析等操作的脚本
├── third_party         # 外部辅助工具，分叉代码和其他第三方工具
├── test                # 额外的外部测试应用程序和测试数据
├── build               # 打包和持续集成
└── deployments         # IaaS、PaaS、系统和容器编排部署配置和模板
```

<br>

### 运行

> make run

