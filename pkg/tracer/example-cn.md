
### 启动jaeger和elasticsearch服务

链路跟踪使用jaeger，存储使用elasticsearch，在本地使用[docker-compose](https://github.com/docker/compose/releases)启动两个服务。

**(1) elasticsearch服务**

这是 [elasticsearch服务的启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/elasticsearch)，`.env`文件是elasticsearch的启动配置，启动elasticsearch服务：

> docker-compose up -d

<br>

**(2) jaeger服务**

这是 [jaeger服务的启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/jaeger)，`.env`文件是配置jaeger信息，启动jaeger服务：

> docker-compose up -d

在浏览器访问jaeger查询主页 [http://localhost:16686](http://localhost:16686) 。

<br>

### 单服务链路跟踪示例

以`⓵基于sql创建web服务`代码为例，修改配置文件`configs/user.yml`，开启链路跟踪功能(字段enableTrace)，并且填写jaeger配置信息。

如果想跟踪redis，启用redis缓存，把yaml配置文件里的缓存类型字段**cacheType**值改为redis，并配置redis地址，同时在本地使用docker启动redis服务，这是[redis服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/redis)。

运行web服务：

```bash
# 编译和运行服务
make run
```

复制 [http://localhost:8080/swagger/index.html](http://localhost:8080/apis/swagger/index.html) 到浏览器访问swagger主页，以请求get查询为例，连续请求同一个id两次，链路跟踪如下图所示。

![one-server-trace](https://go-sponge.com/assets/images/one-server-trace.jpg)

<br>

从图中可以看到第一次请求有4个span，分别是：

- 请求接口 /api/v1/teacher/1
- 查询redis
- 查询mysql
- 设置redis缓存

说明第一次请求从redis查找，没有命中缓存，然后从mysql读取数据，最后设置缓存。

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

### 多服务链路跟踪示例

以一个极简版的电商微服务集群为例，点击查看[源码](https://github.com/zhufuyi/sponge_examples/tree/main/6_micro-cluster)，一个共四个服务**shopgw**、**product**、**inventory**、**comment**，分别修改4个服务yaml配置(在configs目录下)，开启链路跟踪功能，并且填写jaeger配置信息。

在 **product**、**inventory**、**comment** 三个服务的**internal/service**目录下找到模板文件，填充代码替代`panic("implement me")`，使得代码可以正常执行，并且手动添加一个**span**，添加随机延时。

启动 **shopgw**、**product**、**inventory**、**comment** 四个服务，在浏览器访问 [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html) ，执行get请求，链路跟踪界面如下图所示。

![multi-servers-trace](https://go-sponge.com/assets/images/multi-servers-trace.jpg)

<br>

从图中可以看到共有10个span，主要链路：

- 请求接口/api/v1/detail
- shopgw 服务调用product的grpc客户端
- product 的grpc服务端
- product 服务中手动添加的mockDAO
- shopgw 服务调用inventory的grpc客户端
- inventory 的grpc服务端
- inventory 服务中手动添加的mockDAO
- shopgw 服务调用comment的grpc客户端
- comment 的grpc服务端
- comment 服务中手动添加的mockDAO

shopgw服务串行调用了**product**、**inventory**、**comment** 三个服务获取数据，实际中可以改为并行调用会更节省时间，但是要注意控制协程数量。

<br>