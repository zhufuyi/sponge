Change the `prometheus_data` directory permissions to 0777 before starting prometheus.

> docker-compose up -d

Open in your browser http://localhost:9090

<br>

### prometheus examples of common target service indicators collected

#### consoul automatic service discovery

Commonly used for cluster monitoring, using consoul's registration and discovery features, prometheus implements exporter auto-discovery and supports tagging attributes.

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
  # from configuration
  - job_name: 'node-exporter'
    #scrape_interval: 5s
    static_configs:
      - targets: ['192.168.111.128:9100']
        labels:
          env: 'dev'

  # from file
 - job_name: "node-exporter-2"
   file_sd_configs:
   - refresh_interval: 1m
     files: 
      - "/etc/prometheus/conf.d/node-exporter/*.yml"
```

<br>

Example of the contents of an imported file, using the file import method.

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

#### process

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

#### containers

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
