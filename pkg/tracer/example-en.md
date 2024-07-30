### Starting Jaeger and Elasticsearch Services

Distributed tracing uses Jaeger for tracing and Elasticsearch for storage. You can start both services locally using [docker-compose](https://github.com/docker/compose/releases).

**(1) Elasticsearch Service**

Here is the [script for starting the Elasticsearch service](https://github.com/zhufuyi/sponge/tree/main/test/server/elasticsearch). The `.env` file contains Elasticsearch configuration. To start the Elasticsearch service, run:

> docker-compose up -d

<br>

**(2) Jaeger Service**

Here is the [script for starting the Jaeger service](https://github.com/zhufuyi/sponge/tree/main/test/server/jaeger). The `.env` file contains Jaeger configuration. To start the Jaeger service, run:

> docker-compose up -d

Access the Jaeger query homepage in your browser at [http://localhost:16686](http://localhost:16686).

<br>

### Single-Service Distributed Tracing Example

Taking the code for the `⓵ Web Service Based on SQL` as an example, modify the configuration file `configs/user.yml` to enable distributed tracing (set the `enableTrace` field to true) and provide Jaeger configuration details.

If you want to trace Redis and use Redis caching, change the cache type field **cacheType** to "redis" in the YAML configuration file and configure the Redis address. Additionally, start a Redis service locally using Docker with this [script](https://github.com/zhufuyi/sponge/tree/main/test/server/redis).

Run the web service:

```bash
# Compile and run the service
make run
```

Copy [http://localhost:8080/swagger/index.html](http://localhost:8080/apis/swagger/index.html) into your browser to access the Swagger homepage. As an example, for a GET request, make two consecutive requests with the same ID. The distributed tracing results are shown in the following image:

![one-server-trace](https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/one-server-trace.jpg)

From the image, you can see that the first request consists of 4 spans:

- Request to the interface /api/v1/teacher/1
- Redis query
- MySQL query
- Redis cache set

This indicates that the first request checked Redis, did not find a cache, retrieved data from MySQL, and finally set the cache.

The second request only has 2 spans:

- Request to the interface /api/v1/teacher/1
- Redis query

This means that the second request directly hit the cache, skipping the MySQL query and cache setting processes.

These spans are automatically generated, but often you may need to manually add custom spans. Here's an example of adding a span:

```go
import "github.com/zhufuyi/sponge/pkg/tracer"

tags := map[string]interface{}{"foo": "bar"}
_, span := tracer.NewSpan(ctx, "spanName", tags)  
defer span.End()
```

<br>

### Multi-Service Distributed Tracing Example

Taking a simplified e-commerce microservices cluster as an example, you can see the [source code here](https://github.com/zhufuyi/sponge_examples/tree/main/6_micro-cluster). This cluster consists of four services: **shopgw**, **product**, **inventory**, and **comment**. Modify the YAML configuration for each of these services (located in the `configs` directory) to enable distributed tracing and provide Jaeger configuration details.

In the **product**, **inventory**, and **comment** services, locate the template files in the **internal/service** directory and replace `panic("implement me")` with code that allows the service to run correctly. Additionally, manually add a **span** and introduce random delays in the code.

Start the **shopgw**, **product**, **inventory**, and **comment** services. Access [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html) in your browser and execute a GET request. The distributed tracing interface will look like the image below:

![multi-servers-trace](https://raw.githubusercontent.com/zhufuyi/sponge_examples/main/assets/multi-servers-trace.jpg)

From the image, you can see a total of 10 spans in the primary trace:

- request to the /api/v1/detail interface
- shopgw service invoking the grpc client-side of product
- grpc server-side in the product service
- manually added mockDAO in the product service
- shopgw service invoking the grpc client of inventory
- grpc server in the inventory service
- manually added mockDAO in the inventory service
- shopgw service invoking the grpc client of comment
- grpc server in the comment service
- manually added mockDAO in the comment service

The shopgw service sequentially calls the **product**, **inventory**, and **comment** services to fetch data. In practice, you can optimize this by making parallel calls to save time, but be mindful of controlling the number of concurrent goroutines.

<br>
