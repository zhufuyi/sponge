### Start the elasticsearch service

Switch to the elasticsearch directory and start the service

> docker-compose up -d

<br>

### Start jaeger service

Switch to the jaeger directory, open the .env file, and fill in the url, login account, and password of elastic search

Start the jaeger service

> docker-compose up -d

Check whether it is normal

> docker-compose ps

Visit the http://192.168.3.37:16686 query page in your browser

<br>

#### port

| Port  | Protocol | Component | Function                                                                                     |
|-------|----------|-----------|----------------------------------------------------------------------------------------------|
| 6831  | UDP      | agent     | accept jaeger.thrift over Thrift-compact protocol (used by most SDKs)                        |
| 6832  | UDP      | agent     | accept jaeger.thrift over Thrift-binary protocol (used by Node.js SDK)                       |
| 5775  | UDP      | agent     | (deprecated) accept zipkin.thrift over compact Thrift protocol (used by legacy clients only) |
| 5778  | HTTP     | agent     | serve configs (sampling, etc.)                                                               |
| 16686 | HTTP     | query     | serve frontend                                                                               |
| 4317  | HTTP     | collector | accept OpenTelemetry Protocol (OTLP) over gRPC, if enabled                                   |
| 4318  | HTTP     | collector | accept OpenTelemetry Protocol (OTLP) over HTTP, if enabled                                   |
| 14268 | HTTP     | collector | accept jaeger.thrift directly from clients                                                   |
| 14250 | HTTP     | collector | accept model.proto                                                                           |
| 9411  | HTTP     | collector | Zipkin compatible endpoint (optional                                                         |