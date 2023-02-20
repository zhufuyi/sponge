[sponge](https://github.com/zhufuyi/sponge) 是一个微服务框架，也是一个快速创建web和微服务代码工具。sponge拥有丰富的生成代码命令，一共生成12种不同功能代码，这些功能代码可以组合成完整的服务(类似人为打散的海绵细胞可以自动重组成一个新的海绵)。微服务代码功能包括日志、服务注册与发现、注册中心、限流、熔断、链路跟踪、指标监控、pprof性能分析、统计、缓存、CICD等功能。代码解耦模块化设计，创建的web和rpc服务包括了从开发到部署完整工程代码和脚本，让使用go语言开发更加轻松和高效。

<br>

### sponge 生成代码框架

生成代码基于**Yaml**、**SQL**和**Protobuf**三种方式，每种方式拥有生成不同功能代码，生成代码的框架图如下所示：

<p align="center">
<img width="1500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/sponge-framework.png">
</p>

<br>

### 微服务框架

sponge创建的微服务代码框架如下图所示，这是典型的微服务分层结构，具有高性能，高扩展性，包含了常用的服务治理功能。

<p align="center">
<img width="1000px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/microservices-framework.png">
</p>

<br>

### 在线生成代码demo

在线生成代码demo： [https://go-sponge.com/ui](https://go-sponge.com/ui/micro-rpc-gw-pb)

💡 警告：

> 有部分生成代码需要填写mysql账号和密码，不要在这里尝试，以免被暴露风险。
> 服务器资源有限，后面有可能无法使用，建议在本地搭建同样的生成代码平台来使用。

<br>

### 快速安装

**(1) 安装 sponge**

```bash
go install github.com/zhufuyi/sponge/cmd/sponge@latest
```

<br>

**(2) 安装 protoc**

protoc 下载地址 `https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3`， 然后把 **protoc** 二进制文件添加到系统path下。

<br>

**(3) 安装依赖插件和工具**

```bash
sponge init
```
如果有插件安装出错 执行命令重试 `sponge tools --install`

<br>

💡 注意：

> 如果使用windows环境, 还需要安装额外依赖工具, 安装详情看 [windows dependency tools](https://go-sponge.com/zh-cn/sponge-install?id=window%e7%8e%af%e5%a2%83%e5%ae%89%e8%a3%85%e4%be%9d%e8%b5%96%e5%b7%a5%e5%85%b7).

<br>

### 快速开始

启动命令行的UI服务：

```bash
sponge run
```

在浏览器访问 `http://localhost:24631`。

💡 注意：

> 不要在sponge二进制文件所在的目录下执行 `sponge run` 命令。

<br>

根据sql生成web项目代码示例：

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

这是生成12种功能代码的其中的一个示例，生成其他11种功能代码可以在ui界面上尝试。

<br>

### 文档

[sponge 使用文档](sponge-doc-cn.md)

<br>

### 视频

- [sponge的形成过程](https://www.bilibili.com/video/BV1s14y1F7Fz/)
- [sponge的框架介绍](https://www.bilibili.com/video/BV13u4y1F7EU/)
- [十分钟搭建一个拥有多个微服务的go语言工程项目](https://www.bilibili.com/video/BV1pP4y1y7hA/)

<br>

觉得对你有用的话请给个star⭐，欢迎加入[go sponge微信群交流](https://pan.baidu.com/s/1NZgPb2v_8tAnBuwyeFyE_g?pwd=spon)。

![wechat-group](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/wechat-group.png)

<br><br>
