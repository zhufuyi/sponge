### Starting Prometheus and Grafana Services

**(1) Prometheus Service**

Here is the [script for starting the Prometheus service](https://github.com/zhufuyi/sponge/tree/main/test/server/monitor/prometheus). Start the Prometheus service:

```bash
docker-compose up -d
```

Access the Prometheus homepage in your browser at [http://localhost:9090](http://localhost:9090/).

<br>

**(2) Grafana Service**

Here is the [script for starting the Grafana service](https://github.com/zhufuyi/sponge/tree/main/test/server/monitor/grafana). Start the Grafana service:

```bash
docker-compose up -d
```

Access the main Grafana page in your browser at [http://localhost:33000](http://localhost:33000), and configure the Prometheus data source to be `http://localhost:9090`.

> [!attention] When importing JSON dashboards into Grafana, the **datasource** value in the JSON must match the name of the Prometheus data source you set in Grafana (in this case, **Prometheus**), or the graphs won't display data.

<br>

### Web Service Monitoring Example

Taking the code for the `⓵ Web Service Based on SQL` as an example, it provides a default metrics interface at [http://localhost:8080/metrics](http://localhost:8080/metrics).

**(1) Adding Monitoring Targets to Prometheus**

Open the Prometheus configuration file `prometheus.yml` and add scraping targets:

```bash
  - job_name: 'http-edusys'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:8080']
```

> [!attention] Before starting the Prometheus service, make sure to change the permissions of the `prometheus.yml` file to `0777`. Otherwise, changes made to `prometheus.yml` using `vim` won't be synchronized with the container.

Trigger the Prometheus configuration to take effect by executing `curl -X POST http://localhost:9090/-/reload`. Wait for a moment, then access [http://localhost:9090/targets](http://localhost:9090/targets) in your browser to check if the newly added scraping target is active.

<br>

**(2) Adding Monitoring Dashboards to Grafana**

Import the [HTTP monitoring dashboard](https://github.com/zhufuyi/sponge/blob/main/pkg/gin/middleware/metrics/gin_grafana.json) into Grafana. If the monitoring interface does not display data, check if the data source name in the JSON matches the Prometheus data source name in Grafana's configuration.

<br>

**(3) Load Testing the API and Observing Monitoring Data**

Use the [wrk](https://github.com/wg/wrk) tool to perform load testing on the API:

```bash
# test 1
wrk -t2 -c10 -d10s http://192.168.3.27:8080/api/v1/teacher/1

# test 2
wrk -t2 -c10 -d10s http://192.168.3.27:8080/api/v1/course/1
```

The monitoring interface will look like the following image:

![http-grafana](https://go-sponge.com/assets/images/http-grafana.jpg)

<br>

### GRPC Service Monitoring Example

Taking the code for the `⓶Create grpc service based on sql` as an example, it provides a default metrics interface at [http://localhost:8283/metrics](http://localhost:8283/metrics).

**(1) Adding Monitoring Targets to Prometheus**

Open the Prometheus configuration file `prometheus.yml` and add scraping targets:

```yaml
  - job_name: 'rpc-server-user'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:8283']
```

> [!attention] Before starting the Prometheus service, make sure to change the permissions of the `prometheus.yml` file to `0777`. Otherwise, changes made to `prometheus.yml` using `vim` won't be synchronized with the container.

Trigger the Prometheus configuration to take effect by executing `curl -X POST http://localhost:9090/-/reload`. Wait for a moment, then access [http://localhost:9090/targets](http://localhost:9090/targets) in your browser to check if the newly added scraping target is active.

<br>

**(2) Adding Monitoring Dashboards to Grafana**

Import the [grpc service monitoring dashboard](https://github.com/zhufuyi/sponge/blob/main/pkg/grpc/metrics/server_grafana.json) into Grafana. If the monitoring interface does not display data, check if the data source name in the JSON matches the Prometheus data source name in Grafana's configuration.

<br>

**(3) Load Testing GRPC Methods and Observing Monitoring Data**

Open the `internal/service/teacher_client_test.go` file using the `Goland` IDE and test various methods under **Test_teacherService_methods** or **Test_teacherService_benchmark**.

The monitoring interface will look like the following image:

![rpc-grafana](https://go-sponge.com/assets/images/rpc-grafana.jpg)

<br>

The monitoring for the grpc client is similar to the server monitoring, and you can find the [grpc client monitoring dashboard](https://github.com/zhufuyi/sponge/blob/main/pkg/grpc/metrics/client_grafana.json) here.

<br>

### Automatically Adding and Removing Monitoring Targets in Prometheus

In real-world scenarios, managing monitoring targets in Prometheus manually can be cumbersome and error-prone, especially when dealing with a large number of services. Prometheus supports dynamic configuration using service discovery with tools like `Consul` to automatically add and remove monitoring targets.

Start a local Consul service using the [Consul service start script](https://github.com/zhufuyi/sponge/tree/main/test/server/consul):

```bash
docker-compose up -d
```

Open the Prometheus configuration file `prometheus.yml` and add Consul configuration:

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

Trigger the Prometheus configuration to take effect by executing `curl -X POST http://localhost:9090/-/reload`.

Once you've configured Prometheus with Consul service discovery, you can push service address information to Consul. Push the information using a JSON file named `user_exporter.json` with content like this:

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

Wait for a moment, then open [http://localhost:9090/targets](http://localhost:9090/targets) in your browser to check if the newly added scraping target is active. Then, stop the service and wait for a moment to check if the target is automatically removed.

> [!tip] In web or grpc service, the JSON information is usually submitted to Consul automatically via program code, not through commands. After the web or grpc service starts normally, Prometheus can dynamically discover monitoring targets. When the web or grpc service stops, Prometheus automatically removes the monitoring target.

<br>
