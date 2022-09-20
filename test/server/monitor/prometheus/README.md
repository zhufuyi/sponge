启动前把 `prometheus_data` 目录权限改为0777，启动prometheus：

> docker-compose up -d


在浏览器打开 http://localhost:9090

<br>

### prometheus 采集常用目标服务指标示例

#### consoul 服务自动发现

常用于集群监控，利用consoul的注册与发现功能，prometheus实现exporter自动发现，并且支持标签属性。

```yaml
  - job_name: 'consul-node-exporter'
    consul_sd_configs:
      - server: '192.168.3.36:8500'
        services: []  
    relabel_configs:
      - source_labels: [__meta_consul_tags]
        regex: .*node-exporter.*
        action: keep
      - regex: __meta_consul_service_metadata_(.+)
        action: labelmap
```

<br>

#### linux

```yaml
  - job_name: 'node-exporter'
    #scrape_interval: 5s
    static_configs:
      - targets: ['192.168.111.128:9100']
        labels:
          env: 'dev'

  # 通过导入文件方式
 - job_name: "node-exporter-2"
   file_sd_configs:
   - refresh_interval: 1m
     files: 
      - "/etc/prometheus/conf.d/node-exporter/*.yml"
```

使用文件导入方式，导入的文件内容示例：

```yaml
- targets: ['node-exporter:9100']
  labels:
    env: 'dev'
```

<br>

#### windows

```yaml
 - job_name: 'windows-exporter'
#   scrape_interval: 5s
   static_configs:
     - targets: ['192.168.6.169:9182']
       labels:
          env: 'dev'
```

<br>

#### 进程

```yaml
  - job_name: 'process-exporter'
#    scrape_interval: 5s
    static_configs:
      - targets: ['192.168.111.128:9256']
        labels:
          env: 'dev'
```

<br>

### grpc

```yaml
 - job_name: 'hello_grpc_server'
   scrape_interval: 2s
   static_configs:
     - targets: ['192.168.3.27:9092']
       labels:
         env: 'dev'

 - job_name: 'hello_grpc_client'
   scrape_interval: 2s
   static_configs:
     - targets: ['192.168.3.27:9094']
       labels:
         env: 'dev'

  - job_name: 'gin_exporter'
    scrape_interval: 2s
    static_configs:
      - targets: ['192.168.3.27:6060']
        labels:
          env: 'dev'
```

<br>


#### prometheus

```yaml
  - job_name: 'prometheus'
    #scrape_interval: 5s
    static_configs:
      - targets: ['192.168.3.36:9090']
```

</br>

#### 容器

```yaml
  - job_name: 'cadvisor'
    #scrape_interval: 15s
    static_configs:
      - targets: ['192.168.111.128:9192']
```

<br>

#### url

```yaml
  - job_name: 'blackbox-exporter'
    metrics_path: /probe
    params:
      module: [http_2xx]  # Look for a HTTP 200 response.
    static_configs:
      - targets:
        - https://www.baidu.com
        - http://myexample.com:8080
        - https://prometheus.io
        - https://zhuyasen.com
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 192.168.111.128:9115  # The blackbox exporter's real hostname:port.
```
</br>

#### thanos

```yaml
  - job_name: 'thanos'
#    scrape_interval: 5s
    static_configs:
      - targets:
          - 'thanos-sidecar-1:10902'
#          - 'thanos-sidecar-2:10902'
          - 'thanos-querier:10902'
          - 'thanos-store-gateway:10902'
          - 'thanos-compactor:10902'
          - 'thanos-ruler:10902'
```

<br>

#### minio

```yaml
  - job_name: 'minio'
#    scrape_interval: 5s
    static_configs:
      - targets: ['minio:9000']
    metrics_path: /minio/prometheus/metrics
```

</br>

#### mysql

```yaml
  - job_name: 'mysql-exporter'
#    scrape_interval: 15s
    static_configs:
      - targets: ['192.168.111.128:9104']
        labels:
          env: 'dev'
```

</br>

#### redis

```yaml
  - job_name: 'redis-exporter'
#    scrape_interval: 5s
    static_configs:
      - targets: ['192.168.111.128:9121']
        labels:
          env: 'dev'
```

</br>

#### rabbitmq

```yaml
 - job_name: 'rabbitmq-exporter'
#   scrape_interval: 5s
   static_configs:
     - targets: ['192.168.7.76:9419']
       labels:
         env: 'dev'
```

</br>

#### kafka

```yaml
 - job_name: 'kafka-exporter'
#   scrape_interval: 5s
   static_configs:
     - targets: ['10.201.0.112:9308']
       labels:
         env: 'dev'
```

</br>

#### loki

```yaml
  - job_name: 'loki-exporter'
    #scrape_interval: 5s
    static_configs:
      - targets: ['192.168.111.128:3100']
        labels:
          env: 'dev'

  - job_name: 'loki-canary-exporter'
    #scrape_interval: 5s
    static_configs:
      - targets: ['192.168.111.128:3500']
        labels:
          env: 'dev'

  - job_name: 'promtail-exporter'
    #scrape_interval: 5s
    static_configs:
      - targets: ['192.168.111.128:3101']
        labels:
          env: 'dev'
```

</br>
