[sponge](https://github.com/zhufuyi/sponge) 是一个微服务框架，一个快速创建web和微服务代码工具。sponge拥有丰富的生成代码命令，一共生成12种不同功能代码，这些功能代码可以组合成完整的服务(类似人为打散的海绵细胞可以自动重组成一个新的海绵)。微服务代码功能包括日志、服务注册与发现、注册中心、限流、熔断、链路跟踪、指标监控、pprof性能分析、统计、缓存、CICD等功能。代码解耦模块化设计，包括了从开发到部署完整工程，常用代码和脚本是自动生成，只需在按照代码模板去编写业务逻辑代码，使得开发效率提高不少。

<br>

### sponge 生成代码框架

生成代码基于**Yaml**、**SQL DDL**和**Protocol buffers**三种方式，每种方式拥有生成不同功能代码，生成代码的框架图如下所示：

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>

<br>

### 微服务框架

sponge创建的微服务代码框架如图下图所示，这是典型的微服务分层结构，具有高性能，高扩展性，包含常用的服务治理功能。

<p align="center">
<img width="1200px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>

<br>

### 安装

**(1) 安装 sponge**

```bash
go install github.com/zhufuyi/sponge/cmd/sponge@latest
```

**(2) 安装依赖插件和工具**

```bash
sponge init
```

如果有插件安装出错(`protoc`除外), 执行命令重试 `sponge tools --install`

<br>

**(3) 安装 protoc**

protoc 下载地址 `https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3`， 然后把 **protoc** 二进制文件添加到系统path下。

<br>

💡 注意：

> 如果使用windows环境, 还需要安装额外依赖工具, 安装详情看 [windows dependency tools](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-doc-cn.md#21-window-environment-installation-dependencies).

<br>

### 快速开始

启动命令行的UI服务：

```bash
sponge run
```

在浏览器访问 `http://localhost:24631`。

💡 注意：

> 不要在sponge二进制文件所在的目录下执行 "sponge run"命令，否则生成代码时会报错：
>
>> exec: "sponge": cannot run executable found relative to current directory
>>

<br>

根据sql生成web项目代码的一个示例如下图所示：

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/web-http.gif">
</p>

下载web项目代码后，执行命令启动服务：

```bash
# 更新swagger文档
make docs

# 编译和运行服务
make run
```

<br>

生成其他类型代码可以自己尝试。

<br>

### 文档

[sponge 使用文档](sponge-doc-cn.md)

<br><br>
