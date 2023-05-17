[sponge](https://github.com/zhufuyi/sponge) 是一个强大的生成web和微服务代码工具，也是一个基于gin和grpc封装的微服务框架。sponge拥有丰富的生成代码命令，生成不同的功能代码可以组合成完整的服务(类似人为打散的海绵细胞可以自动重组成一个新的海绵)。微服务代码功能包括日志、服务注册与发现、注册中心、限流、熔断、链路跟踪、指标监控、pprof性能分析、统计、缓存、CICD等功能。代码解耦模块化设计，很容易构建出从开发到部署的完整工程代码，让使用go语言开发更便捷、轻松、高效。

<br>

### 生成代码的命令框架

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

### 支持生成的代码类型

- 基于mysql表生成完整的web服务代码。
- 基于mysql表生成handler代码，包括CRUD代码。
- 基于mysql表生成dao代码，包括CRUD代码。
- 基于mysql表生成model代码。
- 基于mysql表生成完整的rpc服务代码。
- 基于mysql表生成service代码，包括CRUD代码。
- 基于mysql表生成protobuf代码。
- 基于protobuf生成web服务代码。
- 基于protobuf生成rpc服务代码。
- 基于protobuf生成rpc网关服务代码。
- 基于yaml生成对应的go结构体代码。
- 根据参数生成rpc连接代码。
- 根据参数生成cache代码。

生成的这些代码可以组合成实际项目的web或微服务代码，开发者只需专注写业务逻辑代码，在线UI界面demo： [https://go-sponge.com/ui](https://go-sponge.com/ui)

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

在浏览器访问 `http://localhost:24631`，在页面上操作生成代码。

<br>

生成web项目代码示例：

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

💡 如果不想使用ui界面，可以使用sponge命令行生成代码，命令行帮助信息里面有丰富的示例，有一些生成代码命令比使用UI界面更便捷。

<br>

### 文档

[sponge 使用文档](https://go-sponge.com/zh-cn/)

<br>

### 使用示例

- [生成完整的web服务项目代码](https://www.bilibili.com/read/cv23018269)
- [生成通用的web服务项目代码](https://www.bilibili.com/read/cv23040234)
- [生成完整的微服务(gRPC)项目代码](https://www.bilibili.com/read/cv23064432)
- [生成通用的微服务(gRPC)项目代码](https://www.bilibili.com/read/cv23099236)
- [生成rpc网关服务项目代码](https://www.bilibili.com/read/cv23189890)
- [生成微服务集群项目代码](https://www.bilibili.com/read/cv23255594)

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

如果对您有帮助给个star⭐，欢迎加入[go sponge微信群交流](https://pan.baidu.com/s/1NZgPb2v_8tAnBuwyeFyE_g?pwd=spon)。

![wechat-group](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/wechat-group.png)

<br><br>
