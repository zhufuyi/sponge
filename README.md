## sponge

<p align="center">
<img align="center" width="500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/logo.png">
</p>

<div align=center>

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/sponge)](https://goreportcard.com/report/github.com/zhufuyi/sponge)
[![codecov](https://codecov.io/gh/zhufuyi/sponge/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/sponge)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/sponge.svg)](https://pkg.go.dev/github.com/zhufuyi/sponge)
[![Go](https://github.com/zhufuyi/sponge/workflows/Go/badge.svg?branch=main)](https://github.com/zhufuyi/sponge/actions)
[![License: MIT](https://img.shields.io/github/license/zhufuyi/sponge)](https://img.shields.io/github/license/zhufuyi/sponge)

</div>

**sponge** is a go microservices framework, a tool for quickly creating complete microservices codes for http or grpc. Generate `config`, `ecode`, `model`, `dao`, `handler`, `router`, `http`, `proto`, `service`, `grpc` codes from the SQL DDL, which can be combined into full services(similar to how a broken sponge cell automatically reorganises itself into a new sponge).

Features :

- Web framework [gin](https://github.com/gin-gonic/gin)
- RPC framework [grpc](https://github.com/grpc/grpc-go)
- Configuration file parsing [viper](https://github.com/spf13/viper)
- Configuration Center [nacos](https://github.com/alibaba/nacos)
- Logging [zap](go.uber.org/zap)
- Database component [gorm](gorm.io/gorm)
- Caching component [go-redis](github.com/go-redis/redis)
- Documentation [swagger](github.com/swaggo/swag)
- Authorization [authorization](github.com/golang-jwt/jwt)
- Validator [validator](github.com/go-playground/validator)
- Rate limiter [ratelimiter](golang.org/x/time/rate)
- Circuit Breaker [hystrix](github.com/afex/hystrix-go)
- Tracking [opentelemetry](go.opentelemetry.io/otel)
- Monitoring [prometheus](github.com/prometheus/client_golang/prometheus) [grafana](https://github.com/grafana/grafana)
- Service registration and discovery [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/) nacos)
- Performance analysis [go profile](https://go.dev/blog/pprof)
- Code inspection [golangci-lint](https://github.com/golangci/golangci-lint)
- Continuous Integration CI [jenkins](https://github.com/jenkinsci/jenkins)
- Continuous Deployment CD [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

<br>

The directory structure follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

```bash
.
├── api            # Grpc's proto file and corresponding code
├── assets         # Other assets used with the repository (images, logos, etc.)
├── build          # Packaging and continuous integration
├── cmd            # The application's directory
├── configs        # Directory of configuration files
├── deployments    # IaaS, PaaS, system and container orchestration deployment configurations and templates
├─ docs            # Design documentation and interface documentation
├── internal       # Private application and library code
│ ├── cache        # Business wrapper-based cache
│ ├── config       # Go struct for config file mapping
│ ├── dao          # Data access
│ ├── ecode        # Custom business error codes
│ ├── handler      # Business function implementation for http
│ ├── model        # Database model
│ ├── routers      # Http routing
│ ├── server       # Service entry, including http and grpc servers
│ ├── service      # Business function implementation for grpc
│ └── types        # Request and response types for http
├── pkg            # library code that external applications can use
├── scripts        # Scripts for performing various build, install, analysis, etc. operations
├── test           # Additional external test applications and test data
└── third_party    # External helpers, forked code and other third party tools
```

<br>

The development specification follows the [Uber Go Language Coding Specification](https://github.com/uber-go/guide/blob/master/style.md).

<br>

## Quick start

### Install

> go install github.com/zhufuyi/sponge@sponge

### Quickly create a http project

#### Creating a new http server

**(1) Generate http server code**

> sponge http --module-name=account --server-name=account --project-name=account --repo-addr=zhufuyi --db-dsn=root:123456@(127.0.0.1:3306)/test --db-table=student

**(2) Modify the configuration file `configs/<service name>.yml`**

- Modify the redis configuration
- If the field `enableTracing` is true, the jaeger address must be set
- If the field `enableRegistryDiscovery` is true, the etcd address must be set

**(3) Generate swagger documentation**

> make docs

**(4) Start up the server**

Way 1: Run locally in the binary

> make run

Copy `http://localhost:8080/swagger/index.html` to your browser and test the api interface.

Way 2: Run in docker

```bash
# Build the docker image
make docker-image

# Start the service
make deploy-docker

# Check the status of the service, if it is healthy, it started successfully
cd deployments/docker-compose
docker-compose ps
```

Way 3: Run in k8s

```bash
# Build the image
make image-build REPO_HOST=zhufuyi TAG=latest

# Push the image to the remote image repository and delete the local image after a successful upload
make image-push REPO_HOST=zhufuyi TAG=latest

# Deploy k8s
kubectl apply -f deployments/kubernetes/
make deploy-k8s

# Check the status of the service
kubectl get -f account-deployment.yml
```  

You can also use Jenkins to automatically build deployments to k8s.

<br>

#### Creating a new handler

> sponge handler --module-name=account --db-dsn=root:123456@(127.0.0.1:3306)/test --db-table=teacher --out=./account

Start up the server

> make docs && make run

Copy `http://localhost:8080/swagger/index.html` to your browser and test the api interface.

<br>

### Quick create a grpc project

#### Creating a new grpc server

**(1) Generate grpc server code**

> sponge grpc --module-name=account --server-name=account --project-name=account --repo-addr=zhufuyi --db-dsn=root:123456@(127.0.0.1:3306)/test --db-table=student

**(2) Modify the configuration file configs/<server name>.yml**

- Modify the redis configuration
- If the field `enableTracing` is true, the jaeger address must be set
- If the field `enableRegistryDiscovery` is true, the etcd address must be set

**(3) Generating grpc code**

> make proto

**(4) Start up the server**

Way 1: Run locally in the binary

> make run

Use IDE to open the file `internal/service/<table name>_client_test.go` to test the api interface of grpc, you can copy the pressure test report to your browser to view it. Or use the `go test` command to execute the test cases.

Way 2: Run in docker

```bash
# Build the docker image
make docker-image

# Start the service
make deploy-docker

# Check the status of the service, if it is healthy, it started successfully
cd deployments/docker-compose
docker-compose ps
```

Way 3: Run in k8s

```bash
# Build the image
make image-build REPO_HOST=zhufuyi TAG=latest

# Push the image to the remote image repository and delete the local image after a successful upload
make image-push REPO_HOST=zhufuyi TAG=latest

# Deploy k8s
kubectl apply -f deployments/kubernetes/
make deploy-k8s

# Check the status of the service
kubectl get -f account-deployment.yml
```  

You can also use Jenkins to automatically build deployments to k8s.  

<br>

#### Creating a new service

> sponge service --module-name=account --server-name=account --db-dsn=root:123456@(127.0.0.1:3306)/test --db-table=teacher --out=./account

Start up the server

> make proto && make run

Use IDE to open the file `internal/service/<table name>_client_test.go` to test the api interface of grpc.
