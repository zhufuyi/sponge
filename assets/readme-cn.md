[sponge](https://github.com/zhufuyi/sponge) 是一个微服务框架，一个快速创建微服务代码工具。sponge拥有丰富的生成代码命令，一共生成12种不同功能代码，这些功能代码可以组合成完整的服务(类似人为打散的海绵细胞可以自动重组成一个新的海绵)。微服务代码功能包括日志、服务注册与发现、注册中心、限流、熔断、链路跟踪、指标监控、pprof性能分析、统计、缓存、CICD等功能，开箱即用。代码使用解耦分层结构，很容易的添加或替换功能代码。作为一个提升效率工具，常用的重复代码基本是自动生成，只需要根据生成的模板代码示例填充业务逻辑代码。

<br>

## 1 sponge功能介绍

### 1.1 sponge 命令框架

生成代码基于**Yaml**、**SQL DDL**和**Protocol buffers**三种方式，每种方式拥有生成不同功能代码，生成代码的框架图如图1-1所示：

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png)
*图1-1 spong生成代码框架图*

<br>

- **Yaml**生成配置文件对应go struct代码。
- **SQL DDL**生成的代码包括 http、handler、 dao、 model、 proto、rpc、service，分为web和rpc两大类型，web和rpc类型都拥有自己的子模块，每个模块代码都可以单独生成，模块之间像洋葱一层一层独立解耦。代码包括了标准化的CRUD(增删改查)业务逻辑，可以直接编译就使用。
  - **web类型代码**： 生成http服务代码包括 handler、 dao、model 三个子模块代码，如图所示向内包含，同样原理，生成handler模块代码包含dao、 model两个子模块代码。
  - **rpc类型代码**：生成rpc服务代码包括 service、dao、model、protocol buffers 四个子模块，如图所示向内包含，生成service模块代码包括 dao、model、protocol buffers 三个子模块代码。
- **Protocol buffers**生成的代码包括 http-pb, rpc-pb, rpc-gw-pb，同样分为web和rpc两大类型，其中 http-pb, rpc-pb 通常结合**SQL DDL**生成的dao、model代码使用。
  - **http-pb**：http服务代码包括router、handler模板两个子模块，不包括操作数据库子模块，后续的业务逻辑代码填写到handler模板文件上。
  - **rpc-pb**：rpc服务代码包括service模板一个子模块，不包括操作数据库模块，后续的业务逻辑代码填写到service模板文件上。
  - **rpc-gw-pb**：rpc网关其实是http服务，包括router和service模板两个子模块，这里的service模板代码是调用rpc服务相关的业务逻辑代码。

在同一个文件夹内，如果发现最新生成代码和原来代码冲突，sponge会取消此次生成流程，不会对原来代码有任何修改，因此不必担心写的业务逻辑代码被覆盖问题。

<br>

### 1.2 微服务框架

sponge创建的微服务代码框架如图1-2所示，这是典型的微服务分层结构，包含常用的服务治理功能。

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png)

*图1-2 微服务框架图*

<br>

微服务主要功能：

- Web 框架 [gin](https://github.com/gin-gonic/gin)
- RPC 框架 [grpc](https://github.com/grpc/grpc-go)
- 配置解析 [viper](https://github.com/spf13/viper)
- 配置中心 [nacos](https://github.com/alibaba/nacos)
- 日志 [zap](https://go.uber.org/zap)
- 数据库组件 [gorm](https://gorm.io/gorm)
- 缓存组件 [go-redis](https://github.com/go-redis/redis) [ristretto](github.com/dgraph-io/ristretto)
- 文档 [swagger](https://github.com/swaggo/swag)
- 鉴权 [jwt](https://github.com/golang-jwt/jwt)
- 校验 [validator](https://github.com/go-playground/validator)
- 限流 [ratelimit](pkg/shield/ratelimit)
- 熔断 [circuitbreaker](pkg/shield/circuitbreaker)
- 链路跟踪 [opentelemetry](https://go.opentelemetry.io/otel)
- 监控 [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana)
- 服务注册与发现 [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/)
- 性能分析 [go profile](https://go.dev/blog/pprof)
- 代码规范检查 [golangci-lint](https://github.com/golangci/golangci-lint)
- 持续集成部署 CICD [jenkins](https://github.com/jenkinsci/jenkins) [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

<br>

代码目录结构遵循 [project-layout](https://github.com/golang-standards/project-layout)，代码目录结构如下所示：

```bash
.
├── api            # proto文件和生成的*pb.go目录
├── assets         # 其他与资源库一起使用的资产(图片、logo等)目录
├── build          # 打包和持续集成目录
├── cmd            # 程序入口目录
├── configs        # 配置文件的目录
├── deployments    # IaaS、PaaS、系统和容器协调部署的配置和模板目录
├─ docs            # 设计文档和界面文档目录
├── internal       # 私有应用程序和库的代码目录
│ ├── cache        # 基于业务包装的缓存目录
│ ├── config       # Go结构的配置文件目录
│ ├── dao          # 数据访问目录
│ ├── ecode        # 自定义业务错误代码目录
│ ├── handler      # http的业务功能实现目录
│ ├── model        # 数据库模型目录
│ ├── routers      # http路由目录
│ ├── rpcclient    # 连接rpc服务的客户端目录
│ ├── server       # 服务入口，包括http、rpc等
│ ├── service      # rpc的业务功能实现目录
│ └── types        # http的请求和响应类型目录
├── pkg            # 外部应用程序可以使用的库目录
├── scripts        # 用于执行各种构建、安装、分析等操作的脚本目录
├── test           # 额外的外部测试程序和测试数据
└── third_party    # 外部帮助程序、分叉代码和其他第三方工具
```

web服务和rpc服务目录结构基本一致，其中有一些目录是web服务独有(internal目录下的routers、handler、types)，有一些目录是rpc服务独有(internal目录下的service)。

<br><br>

## 2 安装sponge和依赖工具

### 2.1 window环境安装依赖工具

如果使用windows环境，需要先安装相关依赖工具，其他环境忽略即可。

**(1) 安装mingw64**

mingw64是Minimalist GNUfor Windows的缩写，它是一个可自由使用和自由发布的Windows特定头文件和使用GNU工具集导入库的集合，下载预编译源码生成的二进制文件，下载地址：

https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z

下载后解压到`D:\Program Files\mingw64`目录下，修改系统环境变量PATH，新增`D:\Program Files\mingw64\bin`。

**安装make命令**

切换到`D:\Program Files\mingw64\bin`目录，找到`mingw32-make.exe`可执行文件，复制并改名为`make.exe`。

<br>

**(2) 安装cmder**

**cmder** 是一个增强型命令行工具，包含一些sponge依赖的命令(bash、git等)，cmder下载地址：

https://github.com/cmderdev/cmder/releases/download/v1.3.20/cmder.zip

下载后解压到`D:\Program Files\cmder`目录下，修改系统环境变量PATH，新增`D:\Program Files\cmder`。

<br>

### 2.2 安装 sponge

**(1) 安装 go**

下载地址 https://go.dev/dl/ 或 https://golang.google.cn/dl/ 选择版本(>=1.16)安装，把 `$GOROOT/bin`添加到系统path下。

注：如果没有科学上网，建议设置国内代理 `go env -w GOPROXY=https://goproxy.cn,direct`

<br>

**(2) 安装 protoc**

下载地址 https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3 ，把**protoc**文件所在目录添加系统path下。

<br>

**(3) 安装 sponge**

> go install github.com/zhufuyi/sponge/cmd/sponge@latest

注：sponge二进制文件所在目录必须在系统path下。

<br>

**(4) 安装依赖插件和工具**

> sponge init

执行命令后自动安装了依赖插件和工具：[protoc-gen-go](https://google.golang.org/protobuf/cmd/protoc-gen-go)、 [protoc-gen-go-grpc](https://google.golang.org/grpc/cmd/protoc-gen-go-grpc)、 [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate)、 [protoc-gen-gotag](https://github.com/srikrsna/protoc-gen-gotag)、 [protoc-gen-go-gin](https://github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin)、 [protoc-gen-go-rpc-tmpl](https://github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl)、 [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2)、 [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc)、 [golangci-lint](https://github.com/golangci/golangci-lint/cmd/golangci-lint)、 [swag](https://github.com/swaggo/swag/cmd/swag)、 [go-callvis](https://github.com/ofabry/go-callvis)

查看依赖工具安装情况：

```bash
# linux 环境
sponge tools

# windows环境，需要指定bash.exe位置
sponge tools --executor="D:\Program Files\cmder\vendor\git-for-windows\bin\bash.exe"
```

<br>

**sponge**命令的帮助信息有详细的使用示例，在命令后面添加`-h`查看，例如`sponge web model -h`，这是根据mysql表生成gorm的model代码返回的帮助信息。

<br><br>

## 3 快速创建web项目

### 3.1 根据mysql创建http服务

#### 3.1.1 创建一个表

根据mysql的数据表来生成代码，先准备一个mysql服务([docker安装mysql](https://github.com/zhufuyi/sponge/blob/main/test/server/mysql/docker-compose.yaml))，例如mysql有一个数据库school，数据库下有一个数据表teacher，如下面sql所示：

```sql
CREATE DATABASE IF NOT EXISTS school DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

use school;

create table teacher
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime     null,
    updated_at  datetime     null,
    deleted_at  datetime     null,
    name        varchar(50)  not null comment '用户名',
    password    varchar(100) not null comment '密码',
    email       varchar(50)  not null comment '邮件',
    phone       varchar(30)  not null comment '手机号码',
    avatar      varchar(200) null comment '头像',
    gender      tinyint      not null comment '性别，1:男，2:女，其他值:未知',
    age         tinyint      not null comment '年龄',
    birthday    varchar(30)  not null comment '出生日期',
    school_name varchar(50)  not null comment '学校名称',
    college     varchar(50)  not null comment '学院',
    title       varchar(10)  not null comment '职称',
    profile     text         not null comment '个人简介'
)
    comment '老师';

create index teacher_deleted_at_index
    on teacher (deleted_at);
```

把SQL DDL导入mysql创建一个数据库school，school下面有一个表teacher。

<br>

#### 3.1.2 生成http服务代码

打开终端，执行命令：

```bash
sponge web http \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --out=./edusys  
```

查看参数说明命令`sponge web http -h`，注意参数**repo-addr**是镜像仓库地址，如果使用[docker官方镜像仓库](https://hub.docker.com/)，只需填写注册docker仓库的用户名，如果使用私有仓库地址，需要填写完整仓库地址。

<br>

生成完整的http服务代码是在当前目录edusys下，目录结构如下：

```
.
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── cache
│    ├── config
│    ├── dao
│    ├── ecode
│    ├── handler
│    ├── model
│    ├── routers
│    ├── server
│    └── types
└── scripts
```

在edusys目录下的Makefile文件，集成了编译、测试、运行、部署等相关命令，切换到edusys目录下运行服务：

```bash
# 更新swagger文档
make docs

# 编译和运行服务
make run
```

复制 http://localhost:8080/swagger/index.html 到浏览器测试CRUD接口，如图3-1所示。

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/http-swag.jpg)
*图3-1 http swagger文档界面*

<br>

服务默认只开启了指标采集接口、每分钟的资源统计信息，其他服务治理默认是关闭。在实际应用中，根据需要做一些调整：

- 使用redis作为缓存，打开配置文件`configs/edusys.yml`，把**cacheType**字段值改为redis，并且填写**redis**配置地址和端口。
- 默认限流、熔断、链路跟踪、服务注册与发现功能是关闭的，可以打开配置文件`configs/edusys.yml`开启相关功能，如果开启链路跟踪功能，必须填写jaeger配置信息；如果开启服务注册与发现功能，必须填写consul、etcd、nacos其中一种配置信息。
- 如果增加或修改了配置字段名称，执行命令 `sponge config --server-dir=./edusys`更新对应的go struct，只修改字段值不需要执行更新命令。
- 修改CRUD接口对应的错误码信息，打开`ingernal/ecode/teacher_http.go`，修改变量**teacherNO**值，这是唯一不重复的数值，返回信息说明根据自己需要修改，对teacher表操作的接口自定义错误码都在这里添加。

<br>

#### 3.1.3 生成handler代码

一个服务中，通常不止一个数据表，如果添加了新数据表，生成的handler代码如何自动填充到已存在的服务代码中呢，需要用到`sponge web handler`命令，例如添加了两个新数据表**course**和**teach**：

```sql
create table course
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime    null,
    updated_at datetime    null,
    deleted_at datetime    null,
    code       varchar(10) not null comment '课程代号',
    name       varchar(50) not null comment '课程名称',
    credit     tinyint     not null comment '学分',
    college    varchar(50) not null comment '学院',
    semester   varchar(20) not null comment '学期',
    time       varchar(30) not null comment '上课时间',
    place      varchar(30) not null comment '上课地点'
)
    comment '课程';

create index course_deleted_at_index
    on course (deleted_at);


create table teach
(
    id           bigint unsigned auto_increment
        primary key,
    created_at   datetime    null,
    updated_at   datetime    null,
    deleted_at   datetime    null,
    teacher_id   bigint      not null comment '老师id',
    teacher_name varchar(50) not null comment '老师名称',
    course_id    bigint      not null comment '课程id',
    course_name  varchar(50) not null comment '课程名称',
    score        char(5)     not null comment '学生评价教学质量，5个等级：A,B,C,D,E'
)
    comment '老师课程';

create index teach_course_id_index
    on teach (course_id);

create index teach_deleted_at_index
    on teach (deleted_at);

create index teach_teacher_id_index
    on teach (teacher_id);
```

<br>

生成包含CRUD业务逻辑的handler代码：

```bash
sponge web handler \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=course,teach \
  --out=./edusys
```

查看参数说明命令`sponge web handler -h`，参数`out`是指定已存在的服务文件夹edusys，如果参数`out`为空，必须指定`module-name`参数，在当前目录生成handler子模块代码，然后把handler代码复制到文件夹edusys，两种方式效果都一样。

执行命令后，在`edusys/internal`目录下生成了course和teach相关的代码：

```
.
└── internal
      ├── cache
      ├── dao
      ├── ecode
      ├── handler
      ├── model
      ├── routers
      └── types
```

<br>

切换到edusys目录下执行命令运行服务：

```bash
# 更新swagger文档
make docs

# 编译和运行服务
make run
```

复制 http://localhost:8080/swagger/index.html 到浏览器测试CRUD接口，如图3-2所示。

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/http-swag2.jpg)
*图3-2 http swagger文档界面*

实际使用中需要修改自定义的CRUD接口返回错误码和信息，打开文件`ingernal/ecode/course_http.go`修改变量**courseNO**值，打开文件`ingernal/ecode/teach_http.go`修改变量**teachNO**值。

虽然生成了每个数据表的CRUD接口，不一定适合实际业务逻辑，就需要手动添加业务逻辑代码了，数据库操作代码填写到`internal/dao`目录下，业务逻辑代码填写到`internal/handler`目录下。

<br>

### 3.2 根据proto文件创建http服务

如果不需要标准CRUD接口的http服务代码，可以在proto文件自定义接口，使用spong命令生成http服务和接口模板代码。

#### 3.2.1 自定义接口

protocol buffers 语法规则看[官方文档 ](https://developers.google.com/protocol-buffers/docs/overview#syntax)，下面是一个示例文件 teacher.proto 内容，每个方法定义了路由和swagger文档的描述信息，实际应用中根据需要在 message 添加 tag 和 validate 的描述信息。

```protobuf
syntax = "proto3";

package api.edusys.v1;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

option go_package = "edusys/api/edusys/v1;v1";

// 生成*.swagger.json基本信息
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  host: "localhost:8080"
  base_path: ""
  info: {
    title: "edusys api docs";
    version: "v0.0.0";
  };
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
};

service teacher {
  rpc Register(RegisterRequest) returns (RegisterReply) {
    // 设置路由
    option (google.api.http) = {
      post: "/api/v1/Register"
      body: "*"
    };
    // 设置路由对应的swagger文档
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "注册用户",
      description: "提交信息注册",
      tags: "teacher",
    };
  }

  rpc Login(LoginRequest) returns (LoginReply) {
    option (google.api.http) = {
      post: "/api/v1/login"
      body: "*"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "登录",
      description: "登录",
      tags: "teacher",
    };
  }
}

message RegisterRequest {
  string name = 1;
  string email = 2;
  string password = 3;
}

message RegisterReply {
  int64   id = 1;
}

message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginReply {
  string token = 1;
}
```

<br>

#### 3.2.2 生成http服务代码

打开终端，执行命令：

```bash
sponge web http-pb \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --protobuf-file=./teacher.proto \
  --out=./edusys
```

查看参数说明命令 `sponge web http-pb -h`，支持\*号匹配(示例`--protobuf-file=*.proto`)，表示根据批量proto文件生成代码，多个proto文件中至少包括一个service，否则不允许生成代码。

生成http服务代码的目录如下所示，与`sponge web http`生成的http服务代码目录有一些区别，新增了proto文件相关的**api**和**third_party**目录，internal目录下没有**cache**、**dao**、**model**、**handler**、**types**目录，其中**handler**是存放业务逻辑模板代码目录，通过命令会自动生成。

```
.
├── api
│    └── edusys
│          └──v1
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── config
│    ├── ecode
│    ├── routers
│    └── server
├── scripts
└── third_party
```

切换到edusys目录下执行命令运行服务：

```bash
# 生成*pb.go文件、生成handler模板代码、更新swagger文档
make proto

# 编译和运行服务
make run
```

复制 http://localhost:8080/apis/swagger/index.html 到浏览器测试接口，如图3-3所示，请求会返回500错误，因为模板代码(internal/handler/teacher_logic.go文件)直接调用`panic("implement me")`，这是为了提示要填写业务逻辑代码。

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/http-pb-swag.jpg)
*图3-3 http swagger文档界面*

<br>

#### 3.2.3 添加新的接口

根据业务需求，需要添加新接口，分为两种情况：

**(1) 在原来proto文件添加新接口**

打开 `api/edusys/v1/teacher.proto`，例如添加**bindPhone**方法，并填写路由和swagger文档描述信息，完成添加一个新接口。

执行命令：

```bash
# 生成*pb.go文件、生成handler模板代码、更新swagger文档
make proto
```

在`internal/handler`和`internal/ecode`两个目录下生成新的模板文件，然后复制最新生成的模板代码到业务逻辑代码区：

- 在`internal/handler`目录下生成了后缀为 **.gen.日期时间** 模板代码文件(示例teacher_logic.go.gen.xxxx225619)，因为teacher_logic.go已经存在，不会把写的业务逻辑代码覆盖，所以生成新文件。打开文件 teacher_logic.go.gen.xxxx225619，把添加方法**bindPhone**接口的模板代码复制到teacher_logic.go文件中，然后填写业务逻辑代码。
- 在`internal/ecode`目录下生成了后缀为 **.gen.日期时间** 模板代码文件，把 **bindPhone** 接口错误码复制到 teacher_http.go 文件中。
- 删除所有后缀名为 **.gen.日期时间** 文件。

<br>

**(2) 在新proto文件添加接口**

例如新添加了**course.proto**文件，**course.proto**下的接口必须包括路由和swagger文档描述信息，查看**章节3.2.1**，把**course.proto**文件复制到`api/edusys/v1`目录下，完成新添加接口。

执行命令：

```bash
# 生成*pb.go文件、生成handler模板代码、更新swagger文档
make proto
```

在 `internal/handler`、`internal/ecode`、 `internal/routers` 三个目录下生成**course**名称前缀的代码文件，只需做下面两个操作：

- 在`internal/handler/course.go`文件填写业务代码。
- 在`internal/ecode/course_http.go`文件修改自定义错误码和信息说明。

<br>

#### 3.2.4 完善http服务

`sponge web http-pb`命令生成的http服务代码没有`dao`、`cache`、`model`等操作数据的相关代码，使用者可以自己实现，如果使用mysql数据库和redis缓存，可以使用**sponge**工具直接生成`dao`、`cache`、`model`代码。

生成CRUD操作数据库代码命令：

```bash
sponge web dao \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --include-init-db=true \
  --out=./edusys
```

查看参数说明命令 `sponge web dao -h`，参数`--include-init-db`在一个服务中只使用一次，下一次生成`dao`代码时去掉参数`--include-init-db`，否则会造成无法生成最新的`dao`代码，原因是db初始化代码已经存在。

无论是自己实现`dao`代码还是使用sponge生成的`dao`代码，之后都需要做一些操作：

- 在服务的初始化和释放资源代码中加入mysql和redis，打开`cmd/edusys/initial/initApp.go`文件，把调用mysql和redis初始化代码反注释掉，打开`cmd/edusys/initial/registerClose.go`文件，把调用mysql和redis释放资源代码反注释掉，初始代码是一次性更改。
- 生成的`dao`代码，并不能和自定义方法**register**和**login**完全对应，需要手动在文件`internal/dao/teacher.go`补充代码(文件名teacher是表名称)，然后在`internal/handler/teacher.go`填写业务逻辑代码(文件名teacher是proto文件名称)，业务代码中返回错误使用`internal/ecode`目录下定义的错误码，如果直接返回错误信息，请求端会收到unknown错误信息，也就是未定义错误信息。
- 默认使用了本地内存做缓存，改为使用redis作为缓存，在配置文件`configs/edusys.yml`修改字段**cacheType**值为redis，并填写redis地址和端口。

切换到edusys目录下再次运行服务：

```bash
# 编译和运行服务
make run
```

打开 http://localhost:8080/apis/swagger/index.html 再次请求接口，可以正常返回数据了。

<br>

### 3.3 总结

生成http服务代码有mysql和proto文件两种方式：

- 根据mysql生成的http服务代码包括每个数据表的CRUD接口代码，后续添加新接口，可以参考CRUD方式添加业务逻辑代码，新添加的接口需要手动填写swagger描述信息。
- 根据proto文件生成的http服务虽然不包括操作数据库代码，也没有CRUD接口逻辑代码，根据需要可以使用`sponge web dao`命令生成操作数据库代码。添加了新的接口，除了生成handler模板代码，swagger文档、路由注册代码、接口的错误码会自动生成。

两种方式都可以完成同样的http服务接口，根据实际应用选择其中一种，如果做后台管理服务，使用mysql直接生产CRUD接口代码，可以少写代码。对于多数需要自定义接口服务，使用proto文件方式生成的http服务，这种方式自由度也比较高，写好proto文件之后，除了业务逻辑代码，其他代码都是通过插件生成。

<br><br>

## 4 快速创建微服务

### 4.1 根据mysql创建rpc服务

#### 4.1.1 生成rpc服务代码

以**章节3.1.1**的teacher表为例，创建rpc服务：

```bash
sponge micro rpc \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --out=./edusys  
```

查看参数说明命令 `sponge micro rpc -h`，生成rpc服务代码在当前edusys目录下，目录结构如下：

```
.
├── api
│    ├── edusys
│    │    └── v1
│    └── types
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── cache
│    ├── config
│    ├── dao
│    ├── ecode
│    ├── model
│    ├── server
│    └── service
├── scripts
└── third_party
```

在edusys目录下的Makefile文件，集成了编译、测试、运行、部署等相关命令，切换到edusys目录下执行命令运行服务：

```bash
# 生成*pb.go
make proto

# 编译和运行服务
make run
```

rpc服务包括了CRUD逻辑代码，也包括rpc客户端测试和压测代码，使用**Goland**或**VS Code**打开`internal/service/teacher_client_test.go`文件，

- 对 **Test_teacherService_methods** 下的方法测试，测试前要先填写测试参数。
- 执 **Test_teacherService_benchmark** 下的方法压测，测试前要先填写压测参数，执行结束后生成压测报告，复制压测报告文件路径到浏览器查看统计信息，如图4-1所示。

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/performance-test.jpg)
*图4-1 性能测试报告界面*

从服务启动日志看到默认监听**8282**端口(rpc服务)和**8283**端口(采集metrics或profile)，开启了每分钟的打印资源统计信息。在实际应用中，根据需要做一些修改：

- 使用redis作为缓存，打开配置文件`configs/edusys.yml`，把**cacheType**字段值改为redis，并且填写**redis**配置地址和端口。
- 默认限流、熔断、链路跟踪、服务注册与发现功能是关闭的，可以打开配置文件`configs/edusys.yml`开启相关功能，如果开启链路跟踪功能，需要填写jaeger配置信息，如果开启服务注册与发现功能，需要填写consul、etcd、nacos其中一种配置信息。
- 如果增加或修改了配置字段名称，执行命令 `sponge config --server-dir=./edusys` 更新对应的go struct，只修改字段值不需要执行更新命令。
- 修改CRUD方法对应的错误码和错误信息，打开`ingernal/ecode/teacher_rpc.go`，修改变量**teacherNO**值(数值唯一)，返回信息说明根据自己需要修改，对teacher表操作的接口错误信息都在这里添加。

<br>

#### 4.1.2 生成service代码

添加了两个新表course和teach，数据表的结构看章节**3.1.3**，生成service代码：

```bash
sponge micro service \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=course,teach \
  --out=./edusys
```

查看参数说明命令 `sponge micro service -h`，参数`out`指定已存在的rpc服务文件夹edusys，如果参数`out`为空，必须指定`module-name`和`server-name`两个参数，在当前目录下生成service代码，然后手动复制到文件夹edusys，两种方式效果都一样。

执行命令后，在下面目录下生成了course和teach相关的代码，如果添加自定义方法或新的protocol buffers文件，也是在下面目录手动添加代码。

```
.
├── api
│    └── edusys
│          └── v1
└── internal
      ├── cache
      ├── dao
      ├── ecode
      ├── model
      └── service
```

<br>

切换到edusys目录下执行命令运行服务：

```bash
# 更新*.pb.go
make proto

# 编译和运行服务
make run
```

使用**Goland**或**VS Code**打开`internal/service/course_client_test.go`和`internal/service/teach_client_test.go`文件测试CRUD方法，测试前需要先填写参数。

<br>

### 4.2 根据proto文件创建rpc服务

sponge不仅支持基于mysql创建rpc服务，还支持基于proto文件生成rpc服务。

#### 4.2.1 自定义方法

下面是一个proto示例文件 teacher.proto 内容：

```protobuf
syntax = "proto3";

package api.edusys.v1;
option go_package = "edusys/api/edusys/v1;v1";

service teacher {
  rpc Register(RegisterRequest) returns (RegisterReply) {}
  rpc Login(LoginRequest) returns (LoginReply) {}
}

message RegisterRequest {
  string name = 1;
  string email = 2;
  string password = 3;
}

message RegisterReply {
  int64   id = 1;
}

message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginReply {
  string token = 1;
}
```

<br>

#### 4.2.2 生成rpc服务代码

打开终端，执行命令：

```bash
sponge micro rpc-pb \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --protobuf-file=./teacher.proto \
  --out=./edusys
```

查看参数说明命令`sponge micro rpc-pb -h`，支持\*号匹配(示例`--protobuf-file=*.proto`)，表示根据批量proto文件生成代码，多个proto文件中至少包括一个service，否则无法生成代码。

生成rpc服务代码目录如下所示，与`sponge micro rpc`生成的rpc服务代码目录有一些区别，internal目录下没有**cache**、**dao**、**model**子目录。

```
.
├── api
│    └── edusys
│          └── v1
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── config
│    ├── ecode
│    ├── server
│    └── service
├── scripts
└── third_party
```

切换到edusys目录下执行命令运行服务：

```bash
# 生成*pb.go文件、生成service模板代码
make proto

# 编译和运行服务
make run
```

启动rpc服务之后，使用**Goland**或**VS Code**打开`internal/service/teacher_client_test.go`文件，对 **Test_teacher_methods** 下各个方法进行测试，测试前要先填写测试参，会发现请求返回内部错误，因为在模板代码文件`internal/service/teacher.go`(文件名teacher是proto文件名)插入了代码`panic("implement me")`，这是为了提示要填写业务逻辑代码。

<br>

#### 4.2.3 添加新的方法

根据业务需求，需要添加新方法，分为两种情况操作：

**(1) 在原来proto文件添加新方法**

打开 `api/edusys/v1/teacher.proto`，例如添加**bindPhone**方法。

执行命令：

```bash
# 生成*pb.go文件、生成service模板代码
make proto
```

在 `internal/service`和`internal/ecode` 两个目录下生成模板代码，然后复制模板代码到业务逻辑代码区：

- 在`internal/service`目录下生成了后缀为 **.gen.日期时间** 模板代码文件(示例teacher.go.gen.xxxx225732)，因为teacher.go已经存在，不会把原来写的业务逻辑代码覆盖，所以生成了新的文件，打开文件 teacher.go.gen.xxxx225732，把添加**bindPhone**方法的模板代码复制到teacher.go文件中，然后填写业务逻辑代码。
- 在`internal/ecode`目录下生成了后缀为 **teacher_rpc.go.gen.日期时间** 文件，把**bindPhone**方法对应的错误码复制到teacher_rpc.go文件中。
- 删除所有后缀名为 **.gen.日期时间** 文件。

<br>

**(2) 在新的proto文件添加新方法**

例如新添加了**course.proto**文件，把**course.proto**文件复制到`api/edusys/v1`目录下，完成新添加接口。

执行命令：

```bash
# 生成*pb.go文件、生成service模板代码
make proto
```

在 `internal/service`、`internal/ecode`、 `internal/routers` 三个目录下生成**course**名称前缀的代码文件，只需做下面两个操作：

- 在`internal/service/course.go`文件填写业务代码。
- 在`internal/ecode/course_rpc.go`文件修改自定义错误码和信息说明。

<br>

#### 4.2.4 完善rpc服务代码

`sponge micro rpc-pb`命令生成的rpc服务代码没有`dao`、`cache`、`model`等操作数据的相关代码，使用者可以自己实现，如果使用mysql数据库和redis缓存，可以使用**sponge**工具直接生成`dao`、`cache`、`model`代码。

生成CRUD操作数据库代码命令：

```bash
sponge micro dao \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --include-init-db=true \
  --out=./edusys
```

查看参数说明命令`sponge micro dao -h`，参数`--include-init-db`在一个服务中只使用一次，下一次生成`dao`代码时去掉参数`--include-init-db`，否则会造成无法生成最新的`dao`代码，

无论是自己实现`dao`代码还是使用sponge生成的`dao`代码，之后都需要做一些操作：

- 在服务的初始化和释放资源代码中加入mysql和redis，打开`cmd/edusys/initial/initApp.go`文件，把调用mysql和redis初始化代码反注释掉，打开`cmd/edusys/initial/registerClose.go`文件，把调用mysql和redis释放资源代码反注释掉，初始代码是一次性更改。
- 生成的`dao`代码，并不能和自定义方法**register**和**login**完全对应，需要手动在文件`internal/dao/teacher.go`补充代码(文件名teacher是表名称)，然后在`internal/handler/teacher.go`填写业务逻辑代码(文件名teacher是proto文件名称)，业务代码中返回错误使用`internal/ecode`目录下定义的错误码，如果直接返回错误信息，请求端会收到unknown错误信息，也就是未定义错误信息。
- 默认使用了本地内存做缓存，改为使用redis作为缓存，在配置文件`configs/edusys.yml`修改字段**cacheType**值为redis，并填写redis地址和端口。

切换到edusys目录下再次运行服务：

```bash
# 编译和运行服务
make run
```

启动rpc服务之后，使用**Goland**或**VS Code**打开`internal/service/teacher_client_test.go`文件测试各个方法。

<br>

### 4.3 根据proto文件创建rpc gateway服务

微服务通常提供的是细粒度的API，实际提供给客户端是粗粒度的API，需要从不同微服务获取数据聚合在一起组成符合实际需求的API，这是rpc gateway的作用，rpc gateway本身也是一个http服务，如图4-2所示。

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/rpc-gateway.png)
*图4-2 rpc gateway框架图*

<br>

#### 4.3.1 定义protocol buffers

以电商微服务为例，商品详情页面有商品、库存、商品评价等信息，这些信息保存在不同的微服务中，一般很少请求每个微服务获取数据，直接请求微服务会造成网络压力倍增，通常的做法是聚合多个微服务数据一次性返回。

下面四个文件夹，每个文件夹下都有一个简单的proto文件。

- **comment**: 评论服务的proto目录
- **inventory**: 库存服务的proto目录
- **product**: 产品服务的proto目录
- **shopgw**: rpc网关服务的proto目录

```
.
├── comment
│    └── v1
│          └──comment.proto
├── inventory
│    └── v1
│          └── inventory.proto
├── product
│    └── v1
│          └── product.proto
└── shopgw
      └── v1
            └── shopgw.proto
```

**comment.proto**文件内容如下：

```protobuf
syntax = "proto3";

package api.comment.v1;

option go_package = "shopgw/api/comment/v1;v1";

service Comment {
  rpc ListByProductID(ListByProductIDRequest) returns (ListByProductIDReply) {}
}

message ListByProductIDRequest {
  int64 productID = 1;
}

message CommentDetail {
  int64  id=1;
  string username = 2;
  string content = 3;
}

message ListByProductIDReply {
  int32 total = 1;
  int64 productID = 2;
  repeated CommentDetail commentDetails = 3;
}
```

<br>

**inventory.proto**文件内容如下：

```protobuf
syntax = "proto3";

package api.inventory.v1;

option go_package = "shopgw/api/inventory/v1;v1";

service Inventory {
  rpc GetByID(GetByIDRequest) returns (GetByIDReply) {}
}

message GetByIDRequest {
  int64 id = 1;
}

message InventoryDetail {
  int64 id = 1;
  float num = 4;
  int32 soldNum =3;
}

message GetByIDReply {
  InventoryDetail inventoryDetail = 1;
}
```

<br>

**product.proto**文件内容如下：

```protobuf
syntax = "proto3";

package api.product.v1;

option go_package = "shopgw/api/product/v1;v1";

service Product {
  rpc GetByID(GetByIDRequest) returns (GetByIDReply) {}
}

message GetByIDRequest {
  int64 id = 1;
}

message ProductDetail {
  int64 id = 1;
  string name = 2;
  float price = 3;
  string description = 4;
}

message GetByIDReply {
  ProductDetail productDetail = 1;
  int64 inventoryID = 2;
}
```

<br>

**shopgw.proto**文件内容如下，rpc网关服务的proto和其他微服务的proto有一点区别，需要指定方法的路由和swagger的描述信息。

```protobuf
syntax = "proto3";

package api.shopgw.v1;

import "api/product/v1/product.proto";
import "api/comment/v1/comment.proto";
import "api/inventory/v1/inventory.proto";
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

option go_package = "shopgw/api/shopgw/v1;v1";

// default settings for generating *.swagger.json documents
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  host: "localhost:8080"
  base_path: ""
  info: {
    title: "eshop api docs";
    version: "v0.0.0";
  };
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
};

service ShopGw {
  rpc GetDetailsByProductID(GetDetailsByProductIDRequest) returns (GetDetailsByProductIDReply) {
    option (google.api.http) = {
      get: "/api/v1/detail"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "get detail",
      description: "get detail from product id",
      tags: "shopgw",
    };
  }
}

message GetDetailsByProductIDRequest {
  int64 productID = 1;
}

message GetDetailsByProductIDReply {
  api.product.v1.ProductDetail productDetail = 1;
  api.inventory.v1.InventoryDetail inventoryDetail = 2;
  repeated api.comment.v1.CommentDetail commentDetails = 3;
}
```

<br>

#### 4.3.2 生成rpc gateway服务代码

根据**shopgw.proto**文件生成rpc gateway服务代码：

```bash
sponge micro rpc-gw-pb \
  --module-name=shopgw \
  --server-name=shopgw \
  --project-name=eshop \
  --repo-addr=zhufuyi \
  --protobuf-file=./shopgw/v1/shopgw.proto \
  --out=./shopgw
```

查看参数说明命令 `sponge micro rpc-gw-pb -h`，生成的rpc gateway服务代码在当前shopgw目录下，目录结构如下：

```
.
├── api
│    └── shopgw
│          └── v1
├── build
├── cmd
│    └── shopgw
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── config
│    ├── ecode
│    ├── routers
│    ├── rpcclient
│    └── server
├── scripts
└── third_party
```

因为**product.proto**依赖**product.proto**、**inventory.proto**、**comment.proto**文件，复制三个依赖的proto文件到api目录下，api目录结构如下：

```
.
├── comment
│    └── v1
│          └── comment.proto
├── inventory
│    └── v1
│          └── inventory.proto
├── product
│    └── v1
│          └── product.proto
└── shopgw
      └── v1
            └── shopgw.proto
```

<br>

切换到shopgw目录下运行服务：

```bash
# 生成*pb.go文件、生成模板代码、更新swagger文档
make proto

# 编译和运行服务
make run
```

复制 http://localhost:8080/apis/swagger/index.html 到浏览器测试接口，如图4-3所示。请求会返回500错误，因为模板代码(internal/service/shopgw_logic.go文件)直接调用`panic("implement me")`，这是为了提示要填写业务逻辑代码。

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/rpc-gw-swag.jpg)
*图4-3 rpc gatewy的swagger文档界面*

<br>

#### 4.3.3 完善rpc gateway服务代码

**(1) 生成连接rpc服务端代码**

服务还没有连接rpc服务代码，下面是生成连接**product**、**inventory**、**comment**三个rpc服务的客户端代码命令：

```bash
sponge micro rpc-cli \
  --rpc-server-name=comment,inventory,product \
  --out=./shopgw
```

查看参数说明命令 `sponge micro rpc-cli -h`，参数`out`指定已存在的服务文件夹shopgw，生成的代码在`internal/rpcclent`目录下。

<br>

**(2) 初始化和关闭rpc连接**

连接rpc服务端代码包括了初始化和关闭函数，根据调用模板代码填写：

- 启动服务时候初始化，在`cmd/shopgw/initial/initApp.go`文件的代码段`// initializing the rpc server connection`下，根据模板调用初始化函数。
- 在关闭服务时候释放资源，在`cmd/shopgw/initial/registerClose.go`文件的代码段`// close the rpc client connection`下，根据模板调用释放资源函数。

<br>

**(3) 修改配置**

连接**product**、**inventory**、**comment**三个rpc服务代码已经有了，但rpc服务地址还没配置，需要在配置文件`configs/shopgw.yml`的字段`grpcClient`下添加连接product、inventory、comment三个微服务配置信息：

```yaml
grpcClient:
  - name: "product"
    host: "127.0.0.1"
    port: 8201
    registryDiscoveryType: ""
  - name: "inventory"
    host: "127.0.0.1"
    port: 8202
    registryDiscoveryType: ""   
  - name: "comment"
    host: "127.0.0.1"
    port: 8203
    registryDiscoveryType: ""
```

如果rpc服务使用了注册与发现，字段`registryDiscoveryType`填写服务注册发现类型，支持consul、etcd、nacos三种。

生成对应go struct代码：

```bash
sponge config --server-dir=./shopgw
```

<br>

**(4) 填写业务代码**

下面是在模板文件`internal/service/shopgw_logic.go`填写的业务逻辑代码示例，分别从**product**、**inventory**、**comment**三个rpc服务获取数据聚合在一起返回。

```go
package service

import (
	"context"

	commentV1 "shopgw/api/comment/v1"
	inventoryV1 "shopgw/api/inventory/v1"
	productV1 "shopgw/api/product/v1"
	shopgwV1 "shopgw/api/shopgw/v1"
	"shopgw/internal/rpcclient"
)

var _ shopgwV1.ShopGwLogicer = (*shopGwClient)(nil)

type shopGwClient struct {
	productCli   productV1.ProductClient
	inventoryCli inventoryV1.InventoryClient
	commentCli   commentV1.CommentClient
}

// NewShopGwClient creating rpc clients
func NewShopGwClient() shopgwV1.ShopGwLogicer {
	return &shopGwClient{
		productCli:   productV1.NewProductClient(rpcclient.GetProductRPCConn()),
		inventoryCli: inventoryV1.NewInventoryClient(rpcclient.GetInventoryRPCConn()),
		commentCli:   commentV1.NewCommentClient(rpcclient.GetCommentRPCConn()),
	}
}

func (c *shopGwClient) GetDetailsByProductID(ctx context.Context, req *shopgwV1.GetDetailsByProductIDRequest) (*shopgwV1.GetDetailsByProductIDReply, error) {
	productRep, err := c.productCli.GetByID(ctx, &productV1.GetByIDRequest{
		Id: req.ProductID,
	})
	if err != nil {
		return nil, err
	}

	inventoryRep, err := c.inventoryCli.GetByID(ctx, &inventoryV1.GetByIDRequest{
		Id: productRep.InventoryID,
	})
	if err != nil {
		return nil, err
	}

	commentRep, err := c.commentCli.ListByProductID(ctx, &commentV1.ListByProductIDRequest{
		ProductID: req.ProductID,
	})
	if err != nil {
		return nil, err
	}

	return &shopgwV1.GetDetailsByProductIDReply{
		ProductDetail:   productRep.ProductDetail,
		InventoryDetail: inventoryRep.InventoryDetail,
		CommentDetails:  commentRep.CommentDetails,
	}, nil
}
```

再次启动服务：

```bash
# 编译和运行服务
make run
```

在浏览器访问 http://localhost:8080/apis/swagger/index.html ，请求返回503错误(服务不可用)，原因是**product**、**inventory**、**comment**三个rpc服务都还没运行。

**product**、**inventory**、**comment**三个rpc服务代码都还没有，如何正常启动呢。这三个rpc服务的proto文件已经有了，根据章节 **4.2 根据proto文件创建rpc服务** 步骤生成代码和启动服务就很简单了。

<br>

### 4.4 总结

生成rpc服务代码是基于mysql和proto文件两种方式，根据proto文件方式除了支持生成rpc服务代码，还支持生成rpc gateway服务(http)代码：

- 根据mysql生成的rpc服务代码包括每个数据表的CRUD方法逻辑代码和proto代码，后续如果要添加新方法，只需在proto文件定义，手动添加业务逻辑代码可以参考CRUD逻辑代码。
- 根据proto文件生成rpc服务代码不包括操作数据库代码，但可以使用`sponge web dao`命令生成操作数据库代码，根据proto文件生成service模板代码，在模板代码填充业务逻辑代码。
- 根据proto文件生成rpc gateway服务代码，接口定义在proto文件，根据proto文件生成service模板代码，在模板代码填充业务逻辑代码，结合`sponge micro rpc-cli`命令使用。

根据实际场景选择生成对应服务代码，如果主要是对数据表增删改查，根据mysql生成rpc服务可以少写代码；如果更多的是自定义方法，根据proto生成rpc服务更合适；rpc转http使用rpc gateway服务。

<br><br>

## 5 服务治理

### 5.1 链路跟踪

#### 5.1.1 启动jaeger和elasticsearch服务

链路跟踪使用jaeger，存储使用elasticsearch，在本地使用[docker-compose](https://github.com/docker/compose/releases)启动两个服务。

**(1) elasticsearch服务**

这是 [elasticsearch服务的启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/elasticsearch)，**.env**文件是elasticsearch的启动配置，启动elasticsearch服务：

> docker-compose up -d

<br>

**(2) jaeger服务**

这是 [jaeger服务的启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/jaeger)，**.env**文件是配置jaeger信息，启动jaeger服务：

> docker-compose up -d

在浏览器访问jaeger查询主页 [http://localhost:16686](http://localhost:16686) 。

<br>

#### 5.1.2 单服务链路跟踪示例

以 **章节3.1.2** 创建的http服务代码为例，修改配置文件`configs/edusys.yml`，开启链路跟踪功能(字段enableTrace)，并且填写jaeger配置信息。

如果想跟踪redis，启用redis缓存，把缓存类型字段**cacheType**值改为redis，并配置redis配置，同时在本地使用docker启动redis服务，这是[redis服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/redis)。

启动http服务：

```bash
# 编译和运行服务
make run
```

复制 [http://localhost:8080/swagger/index.html](http://localhost:8080/apis/swagger/index.html) 到浏览器访问swagger主页，以请求get查询为例，连续请求同一个id两次，链路跟踪如图5-1所示。

![one-server-trace](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/one-server-trace.jpg)
*图5-1 单服务链路跟踪页面*

<br>

从图中可以看到第一次请求有4个span，分别是：

- 请求接口 /api/v1/teacher/1
- 查询redis
- 查询mysql
- 设置redis缓存

说明第一次请求从redis查找，没有命中缓存，然后从mysql读取数据，最后置缓存。

第二次请求只有2个span，分别是：

- 请求接口 /api/v1/teacher/1
- 查询redis

说明第二次请求直接命中缓存，比第一次少了查询mysql和设置缓存过程。

这些span是自动生成的，很多时候需要手动添加自定义span，添加span示例：

```go
import "github.com/zhufuyi/sponge/pkg/tracer"

tags := map[string]interface{}{"foo": "bar"}
_, span := tracer.NewSpan(ctx, "spanName", tags)  
defer span.End()
```

<br>

#### 5.1.3 多服务链路跟踪示例

以**章节4.3**生成的rpc gateway服务代码为例，一个共四个服务**shopgw**、**product**、**inventory**、**comment**，分别修改4个服务配置(在configs目录下)，开启链路跟踪功能，并且填写jaeger配置信息。

在 **product**、**inventory**、**comment** 三个服务的**internal/service**目录下找到模板文件，填充代码替代`panic("implement me")`，使得代码可以正常执行，并且手动添加一个**span**，添加随机延时。

启动 **shopgw**、**product**、**inventory**、**comment** 四个服务，在浏览器访问 [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html) ，执行get请求，链路跟踪界面如图5-2所示。

![multi-servers-trace](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/multi-servers-trace.jpg)
*图5-2 多服务链路跟踪页面*

从图中可以看到共有10个span，主要链路：

- 请求接口/api/v1/detail
- shopgw 服务调用product客户端
- product 的rpc服务端
- product 服务中手动添加的mockDAO
- shopgw 服务调用inventory客户端
- inventory 的rpc服务端
- inventory 服务中手动添加的mockDAO
- shopgw 服务调用comment客户端
- comment 的rpc服务端
- comment 服务中手动添加的mockDAO

shopgw服务串行调用了**product**、**inventory**、**comment** 三个服务获取数据，实际中可以改为并行调用会更节省时间，但是要注意控制协程数量。

<br>

### 5.2 监控

#### 5.2.1 启动Prometheus和Grafana服务

采集指标用[Prometheus](https://prometheus.io/docs/introduction/overview)，展示使用[Grafana](https://grafana.com/docs/)，在本地使用docker启动两个服务。

**(1) prometheus服务**

这是 [prometheus服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/monitor/prometheus)，启动prometheus服务：

> docker-compose up -d

在浏览器访问prometheus主页 [http://localhost:9090](http://localhost:9090/) 。

<br>

**(2) grafana服务**

这是 [grafana服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/monitor/grafana)，启动grafana服务：

> docker-compose up -d

在浏览器访问 grafana 主页面 [http://localhost:33000](http://localhost:33000) ，设置prometheus的数据源 `http://192.168.3.37:9090` ，记住prometheus的数据源名称(这里是**Prometheus**)，后面导入监控面板的json的**datasource**值要一致。

<br>

#### 5.2.2 http服务监控

以**章节3.1.2**生成的http服务代码为例，默认提供指标接口 [http://localhost:8080/metrics](http://localhost:8080/metrics) 。

**(1) 在prometheus添加监控目标**

打开prometheus配置文件 prometheus.yml，添加采集目标：

```bash
  - job_name: 'http-edusys'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:8080']
```

注：如果使用vim修改 prometheus.yml 文件，修改前必须将文件 prometheus.yml 权限改为`0777`，否则修改配置文件无法同步到容器中。

执行请求使prometheus配置生效 `curl -X POST http://localhost:9090/-/reload`，稍等一会，然后在浏览器访问 [http://localhost:9090/targets](http://localhost:9090/targets) ， 检查新添加的采集目标是否生效。

<br>

**(2) 在grafana添加监控面板**

把 [http 监控面板](https://github.com/zhufuyi/sponge/blob/main/pkg/gin/middleware/metrics/gin_grafana.json) 导入到grafana，如果监控界面没有数据显示，检查json里的数据源名称与grafana配置prometheus数据源名称是否一致。

<br>

**(3) 压测接口，观察监控数据**

使用[wrk](https://github.com/wg/wrk)工具压测接口

```bash
# 接口1
wrk -t2 -c10 -d10s http://192.168.3.27:8080/api/v1/teacher/1

# 接口2
wrk -t2 -c10 -d10s http://192.168.3.27:8080/api/v1/course/1
```

监控界面如图5-3所示。

![http-grafana](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/http-grafana.jpg)
*图5-3 http 服务监控界面*

<br>

#### 5.2.3 rpc服务监控

以**章节4.1.1**生成的rpc服务代码为例，默认提供指标接口 [http://localhost:8283/metrics](http://localhost:8283/metrics) 。

**(1) 在prometheus添加监控目标**

打开prometheus配置文件 prometheus.yml，添加采集目标：

```bash
  - job_name: 'rpc-server-edusys'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:8283']
```

执行请求使prometheus配置生效 `curl -X POST http://localhost:9090/-/reload`，稍等一会，然后在浏览器访问 [http://localhost:9090/targets](http://localhost:9090/targets)  检查新添加的采集目标是否生效。

<br>

**(2) 在grafana添加监控面板**

把 [rpc server 监控面板](https://github.com/zhufuyi/sponge/blob/main/pkg/grpc/metrics/server_grafana.json) 导入到grafana，如果监控界面没有数据显示，检查json里的数据源名称与grafana配置prometheus数据源名称是否一致。

<br>

**(3) 压测rpc方法，观察监控数据**

使用**Goland**或**VS Code**打开`internal/service/teacher_client_test.go`文件，对**Test_teacherService_methods** 或 **Test_teacherService_benchmark** 下各个方法进行测试。

监控界面如图5-4所示。
![rpc-grafana](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/rpc-grafana.jpg)
*图5-4 rpc server监控界面*

<br>

上面是rpc服务端的监控，rpc的客户端的监控也类似，[rpc client 监控面板](https://github.com/zhufuyi/sponge/blob/main/pkg/grpc/metrics/client_grafana.json) 。

<br>

#### 5.2.4 在prometheus自动添加和移除监控目标

实际使用中服务数量比较多，手动添加监控目标到prometheus比较繁琐，也容易出错。prometheus支持使用consul的服务注册与发现进行动态配置，自动添加和移除监控目标。

在本地启动 consul 服务，这是 [consul 服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/consul)

打开 prometheus 配置 prometheus.yml，添加consul配置：

```yaml
  - job_name: 'consul-micro-exporter'
    consul_sd_configs:
      - server: 'localhost:8500'
        services: []  
    relabel_configs:
      - source_labels: [__meta_consul_tags]
        regex: .*edusys.*
        action: keep
      - regex: __meta_consul_service_metadata_(.+)
        action: labelmap
```

执行请求使prometheus配置生效 `curl -X POST http://localhost:9090/-/reload`。

在prometheus配置好consul服务发现之后，接着把服务的地址信息推送到consul，推送信息 edusys_exporter.json 文件内容如下：

```json
{
  "ID": "edusys-exporter",
  "Name": "edusys",
  "Tags": [
    "edusys-exporter"
  ],
  "Address": "localhost",
  "Port": 8283,
  "Meta": {
    "env": "dev",
    "project": "edusys"
  },
  "EnableTagOverride": false,
  "Check": {
    "HTTP": "http://localhost:8283/metrics",
    "Interval": "10s"
  },
  "Weights": {
    "Passing": 10,
    "Warning": 1
  }
}
```

> curl -XPUT --data @edusys_exporter.json http://localhost:8500/v1/agent/service/register

稍等一会，然后在浏览器打开 [http://localhost:9090/targets](http://localhost:9090/targets)  检查新添加的采集目标是否生效。然后关闭服务，稍等一会，检查是否自动移除采集目标。

<br>

对于自己的服务，通常启动服务时同时提交信息到consul，把 edusys_exporter.json 转为go struct，在程序内部调用http client提交给consul。

<br>

### 5.3 采集go程序profile

通常使用pprof工具来发现和定位程序问题，特别是线上go程序出现问题时可以自动把程序运行现场(profile)保存下来，再使用工具pprof分析定位问题。

sponge生成的服务支持 **http接口** 和 **系统信号通知** 两种方式采集profile，默认开启系统信号通知方式，实际使用一种即可。

<br>

#### 5.3.1 通过http采集profile

通过http接口方式采集profile默认是关闭的，如果需要开启，修改配置里的字段`enableHTTPProfile`为true，通常在开发或测试时使用，如果线上开启会有一点点性能损耗，根据实际情况是否开启使用。

默认路由 `/debug/pprof`，结合**go tool pprof**工具，任意时刻都可以分析当前程序运行状况。

<br>

#### 5.3.2 通过系统信号通知采集profile

使用http接口方式，程序后台一直定时记录profile相关信息等，绝大多数时间都不会去读取这些profile，可以改进一下，只有需要的时候再开始采集profile，采集完后自动关闭，sponge生成的服务支持监听系统信号来开启和停止采集profile，默认使用了 **SIGTRAP**(5) 系统信号(建议改为SIGUSR1，windows环境不支持)，发送信号给服务：

```bash
# 通过名称查看服务pid(第二列)
ps aux | grep 服务名称

# 发送信号给服务
kill -trap pid值

# kill -usr1 pid值
```

服务收到系统信号通知后，开始采集profile并保存到`/tmp/服务名称_profile`目录，默认采集时长为60秒，60秒后自动停止采集profile，如果只想采集30秒，发送第一次信号开始采集，大概30秒后发送第二次信号表示停止采集profile，类似开关。默认采集**cpu**、**memory**、**goroutine**、**block**、**mutex**、**threadcreate**六种类型profile，文件格式`日期时间_pid_服务名称_profile类型.out`，示例：

```
xxx221809_58546_edusys_cpu.out
xxx221809_58546_edusys_mem.out
xxx221809_58546_edusys_goroutine.out
xxx221809_58546_edusys_block.out
xxx221809_58546_edusys_mutex.out
xxx221809_58546_edusys_threadcreate.out
```

因为trace的profile文件相对比较大，因此默认没有采集，根据实际需要可以开启采集trace(服务启动时调用prof.EnableTrace())。

获得离线文件后，使用pprof工具使用交互式或界面方式进行分析：

```bash
# 交互式
go tool pprof [options] source

# 界面
go tool pprof -http=[host]:[port] [options] source
```

<br>

#### 5.3.3 自动采集profile

上面都是手动采集profile，通常都是希望出现问题时自动采集profile。sponge生成的服务默认支持自动采集profile，是结合资源统计的告警功能来实现的，告警条件：

- 记录程序的cpu使用率连续3次(默认每分钟一次)，3次平均使用率超过80%时触发告警。
- 记录程序的使用物理内存连续3次(默认每分钟一次)，3次平均占用系统内存超过80%时触发告警。
- 如果持续超过告警阈值，默认间隔15分钟告警一次。

触发告警时，程序内部调用kill函数发送x系统信号通知采集profile，采集的profile文件保存到`/tmp/服务名_profile`目录，其实就是在**通过系统信号通知采集profile的基础**上把手动触发改为自动触发，即使在半夜程序的cpu或内存过高，第二天也可以通过分析profile来发现程序哪里造成cpu或内存过高。

注：自动采集profile不适合windows环境。

<br>

### 5.4 注册中心

sponge生成的服务默认支持[Nacos](https://nacos.io/zh-cn/docs/v2/what-is-nacos.html)配置中心，配置中心作用是对不同环境、不同服务的配置统一管理，有效的解决地静态配置的缺点。

在本地启动nacos服务，这是[nacos服务启动配置](https://github.com/zhufuyi/sponge/tree/main/test/server/nacos)，启动nacos服务之后，在浏览器打开管理界面 http://localhost:8848/nacos/index.html ，登录账号密码进入主界面。

以 **章节3.1.2** 生成的http服务代码为例使用配置中心nacos，在nacos界面创建一个名称空间`edusys`，然后新建配置，Data ID值为`edusys.yml`，Group值为`dev`，配置内容值`configs/edusys.yml`文件内容，如图5-3所示。

![nacos-config](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/nacos-config.jpg)
*图5-3 nacos添加服务配置*

打开edusys目录下配置中心文件`configs/edusys_cc.yml`，填写nacos配置信息：

```yaml
# Generate the go struct command: sponge config --server-dir=./serverDir

# nacos settings
nacos:
  ipAddr: "192.168.3.37"    # server address
  port: 8848                      # listening port
  scheme: "http"               # http or https
  contextPath: "/nacos"     # path
  namespaceID: "ecfe0595-cae3-43a2-9e47-216dc92207f9" # namespace id
  group: "dev"                    # group name: dev, prod, test
  dataID: "edusys.yml"        # config file id
  format: "yaml"                 # configuration file type: json,yaml,toml
```

编译和启动edusys服务：

```bash
# 切换到main.go位置
cd cmd/edusys

# 编译
go build

# 运行
./edusys -enable-cc -c=../../configs/edusys_cc.yml
```

启动服务参数`-c`表示指定配置文件，参数`-enable-cc`表示从配置中心获取配置。

<br>

### 5.5 限流和熔断

sponge创建的服务支持限流和熔断功能，默认是关闭的，打开服务配置文件，修改字段**enableLimit**值为`true`表示开启限流功能，修改字段**enableCircuitBreaker**改为`true`表示开启熔断功能。

限流和熔断使用第三方库 [aegis](https://github.com/go-kratos/aegis)，根据系统资源和错误率自适应调整，由于不同服务器的处理能力不一样，参数也不好设置，使用自适应参数避免每个服务去手动去设置参数麻烦。

<br><br>

## 6 持续集成部署

sponge创建的服务支持在 [jenkins](https://www.jenkins.io/doc/) 构建和部署，部署目标可以是docker、 [k8s](https://kubernetes.io/docs/home/) ，部署脚本在**deployments**目录下，下面以使用jenkins部署到k8s为示例。

### 6.1 搭建 jenkins-go 平台

为了可以在容器里编译go代码，需要构建一个 jenkins-go 镜像，这是已经构建好的 [jenkins-go镜像](https://hub.docker.com/r/zhufuyi/jenkins-go/tags)。如果想自己构建 jenkins-go 镜像，可以参考docker构建脚本[Dokerfile](https://github.com/zhufuyi/sponge/blob/main/test/server/jenkins/Dockerfile)

准备好 jenkins-go 镜像之后，还需要准备一个k8s集群(网上有很多搭建k8s集群教程)，k8s鉴权文件和命令行工具[kubectl](https://kubernetes.io/zh-cn/docs/tasks/tools/#kubectl)，确保在 jenkins-go 容器中有操作k8s的权限。

jenkins-go 启动脚本 docker-compose.yml 内容如下：

```yaml
version: "3.7"
services:
  jenkins-go:
    image: zhufuyi/jenkins-go:2.37
    restart: always
    container_name: "jenkins-go"
    ports:
      - 38080:8080
    #- 50000:50000
    volumes:
      - $PWD/jenkins-volume:/var/jenkins_home
      # docker configuration
      - /var/run/docker.sock:/var/run/docker.sock
      - /usr/bin/docker:/usr/bin/docker
      - /root/.docker/:/root/.docker/
      # k8s api configuration directory, including config file
      - /usr/local/bin/kubectl:/usr/local/bin/kubectl
      - /root/.kube/:/root/.kube/
      # go related tools
      - /opt/go/bin/golangci-lint:/usr/local/bin/golangci-lint
```

启动jenkis-go服务：

> docker-compose up -d

在浏览器访问 [http://localhost:38080](http://localhost:38080) ，第一次启动需要 admin 密钥(执行命令获取 `docker exec jenkins-go cat /var/jenkins_home/secrets/initialAdminPassword`)，然后安装推荐的插件和设置管理员账号密码，接着安装一些需要使用到的插件和一些自定义设置。

**(1) 安装插件**

```bash
# 中文插件
Locale

# 添加参数化构建插件
Extended Choice Parameter

# 添加git参数插件
Git Parameter

# 账号管理
Role-based Authorization Strategy
```

**(2) 设置中文**

点击【Manage Jenkins】->【Configure System】选项，找到【Locale】选项，输入【zh_CN】，勾选下面的选项，最后点击【应用】。

**(3) 配置全局参数**

dashboard --> 系统管理 --> 系统配置 --> 勾选环境变量

设置容器镜像的仓库地址：

```bash
# 开发环境镜像仓库
DEV_REGISTRY_HOST http://localhost:27070

# 测试环境镜像仓库
TEST_REGISTRY_HOST http://localhost:28080

# 生产环境镜像仓库
PROD_REGISTRY_HOST http://localhost:29090
```

<br>

### 6.2 创建模板

创建jenkins新任务的一种相对简单的方法是在创建新任务时导入现有模板，然后修改git存储库地址，第一次使用jenkins还没有模板，可以按照下面步骤创建一个模板：

**(1) 创建新的任务**，如图6-1所示。

![create job](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/createJob.jpg)
*图6-1 创建任务界面*

<br>

**(2) 参数化构设置**，使用参数名`GIT_parameter`，如图6-2所示。

![parametric construction](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/paramSetting.jpg)
*图6-2 设置参数化构建界面*

<br>

**(3) 设置流水线**，如图6-3所示。

![flow line](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/pipelineSetting.jpg)
*图6-3 设置流水线界面*

<br>

**(4) 构建项目**

单击左侧菜单栏上的 **Build with Parameters**，然后选择要分支或tag，如图6-4所示。

![run job](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/building.jpg)
*图6-4 参数化构建界面*

<br>

### 6.3 部署到k8s

以**章节3.1.2**的edusys服务为例，使用jenkins构建和部署到k8s。

第一次构建服务需要做一些前期准备：

(1) 把edusys代码上传到代码仓库。

(2) 准备一个docker镜像仓库，确保jenkins-go所在docker有权限上传镜像到镜像仓库。

(3) 确保在k8s集群节点有权限从镜像拉取镜像，在已登录docker镜像仓库服务器上执行命令生成密钥。

```bash
kubectl create secret generic docker-auth-secret \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson
```

(4) 在k8s创建edusys相关资源。

```bash
# 切换到目录
cd deployments/kubernetes

# 创建名称空间，名称对应spong创建服务参数project-name
kubectl apply -f ./*namespace.yml

# 创建configmap、service
kubectl apply -f ./*configmap.yml
kubectl apply -f ./*svc.yml
```

(5) 如果想使用钉钉通知查看构建部署结果，打开代码库下的 **Jenkinsfile** 文件，找到字段**tel_num**填写手机号码，找到**access_token**填写token值。

<br>

前期准备好之后，在jenkins界面创建一个新任务(名称edusys)，使用上面创建的模板(名称sponge)，然后修改git仓库，保存任务，开始参数化构建，构建结果如图6-5所示：

![run job](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/jenkins-build.jpg)
*图6-5 jenkins构建结果界面*

<br>

使用命令`kubectl get all -n edusys `查看edusys服务在k8s运行状态：

```
NAME                             READY   STATUS    RESTARTS   AGE
pod/edusys-dm-77b4bcccc5-8xt8v   1/1     Running   0          21m

NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/edusys-svc   ClusterIP   10.108.31.220   <none>        8080/TCP   27m

NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/edusys-dm   1/1     1            1           21m

NAME                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/edusys-dm-77b4bcccc5   1         1         1       21m
```

<br>

在本地测试是否可以访问

```bash
# 代理端口
kubectl port-forward --address=0.0.0.0 service/edusys-svc 8080:8080 -n edusys

# 请求
curl http://localhost:8080/api/v1/teacher/1
```

<br>

sponge生成的服务包括了Jenkinsfile、构建和上传镜像脚本、k8s部署脚本，基本不需要修改脚本就可以使用，也可以修改脚本适合自己场景。

<br><br>

如果对你有用给个star，也欢迎加入微信群交流。

![wechat-group](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/wechat-group.png)
