### 启动Prometheus和Grafana服务

**(1) prometheus服务**

这是 [prometheus服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/monitor/prometheus)，启动prometheus服务：

```bash
docker-compose up -d
```

在浏览器访问prometheus主页 [http://localhost:9090](http://localhost:9090/) 。

<br>

**(2) grafana服务**

这是 [grafana服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/monitor/grafana)，启动grafana服务：

```bash
docker-compose up -d
```

在浏览器访问 grafana 主页面 [http://localhost:33000](http://localhost:33000) ，设置prometheus的数据源 `http://localhost:9090` 。

> [!attention] 在grafana导入监控面板的json的**datasource**值，必须与在grafana设置的prometheus的数据源名称(这里是**Prometheus**)要一致，否则图标上无法显示数据。

<br>

### web服务监控示例

以`⓵基于sql创建web服务`代码为例，默认提供指标接口 [http://localhost:8080/metrics](http://localhost:8080/metrics) 。

**(1) 在prometheus添加监控目标**

打开prometheus配置文件`prometheus.yml`，添加采集目标：

```bash
  - job_name: 'http-edusys'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:8080']
```

> [!attention] 在启动Prometheus服务前，必须将文件`prometheus.yml`权限改为`0777`，否则使用vim修改`prometheus.yml`文件无法同步到容器中。

执行请求使prometheus配置生效 `curl -X POST http://localhost:9090/-/reload` ，稍等一会，然后在浏览器访问 [http://localhost:9090/targets](http://localhost:9090/targets)，检查新添加的采集目标是否生效。

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

监控界面如下图所示：

![http-grafana](https://go-sponge.com/assets/images/http-grafana.jpg)

<br>

### grpc服务监控示例

以`⓶基于sql创建grpc服务`代码为例，默认提供指标接口 [http://localhost:8283/metrics](http://localhost:8283/metrics) 。

**(1) 在prometheus添加监控目标**

打开prometheus配置文件`prometheus.yml`，添加采集目标：

```yaml
  - job_name: 'rpc-server-user'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:8283']
```

> [!attention] 在启动Prometheus服务前，必须将文件`prometheus.yml`权限改为`0777`，否则使用vim修改`prometheus.yml`文件无法同步到容器中。

执行请求使prometheus配置生效 `curl -X POST http://localhost:9090/-/reload` ，稍等一会，然后在浏览器访问 [http://localhost:9090/targets](http://localhost:9090/targets)， 检查新添加的采集目标是否生效。

<br>

**(2) 在grafana添加监控面板**

把 [grpc server 监控面板](https://github.com/zhufuyi/sponge/blob/main/pkg/grpc/metrics/server_grafana.json) 导入到grafana，如果监控界面没有数据显示，检查json里的数据源名称与grafana配置prometheus数据源名称是否一致。

<br>

**(3) 压测grpc api，观察监控数据**

使用`Goland` IDE打开`internal/service/teacher_client_test.go`文件，对**Test_teacherService_methods** 或 **Test_teacherService_benchmark** 下各个方法进行测试。

监控界面如下图所示。
![rpc-grafana](https://go-sponge.com/assets/images/rpc-grafana.jpg)

<br>

上面是grpc服务端的监控，grpc的客户端的监控也类似，[grpc client 监控面板](https://github.com/zhufuyi/sponge/blob/main/pkg/grpc/metrics/client_grafana.json) 。

<br>

### 在prometheus自动添加和移除监控目标

实际使用中服务数量比较多，手动添加监控目标到prometheus比较繁琐，也容易出错。prometheus支持使用`consul`的服务注册与发现进行动态配置，自动添加和移除监控目标。

在本地启动 consul 服务，这是 [consul 服务启动脚本](https://github.com/zhufuyi/sponge/tree/main/test/server/consul)，启动consul服务：

```bash
docker-compose up -d
```

打开 prometheus 配置 prometheus.yml，添加consul配置：

```yaml
  - job_name: 'consul-micro-exporter'
    consul_sd_configs:
      - server: 'localhost:8500'
        services: []  
    relabel_configs:
      - source_labels: [__meta_consul_tags]
        regex: .*user.*
        action: keep
      - regex: __meta_consul_service_metadata_(.+)
        action: labelmap
```

执行请求使prometheus配置生效 `curl -X POST http://localhost:9090/-/reload` 。

在prometheus配置好consul服务发现之后，接着把服务的地址信息推送到consul，推送信息 user_exporter.json 文件内容如下：

```json
{
  "ID": "user-exporter",
  "Name": "user",
  "Tags": [
    "user-exporter"
  ],
  "Address": "localhost",
  "Port": 8283,
  "Meta": {
    "env": "dev",
    "project": "user"
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

> curl -XPUT --data @user_exporter.json http://localhost:8500/v1/agent/service/register

稍等一会，然后在浏览器打开 [http://localhost:9090/targets](http://localhost:9090/targets)  检查新添加的采集目标是否生效。然后关闭服务，稍等一会，检查是否自动移除采集目标。

> [!tip] 在web或grpc务中，通常是使用程序代码自动把json信息提交给consul，不是通过命令，web或grpc务正常启动服务后，Prometheus就可以动态获取到监控目标，web或grpc务停止后，Prometheus自动移除监控目标。

<br>
