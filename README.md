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

**sponge** is a microservices framework for quickly creating http or grpc code. Generate codes `config`, `ecode`, `model`, `dao`, `handler`, `router`, `http`, `proto`, `service`, `grpc` from SQL DDL, these codes can be combined into complete services (similar to how a broken sponge cell can automatically reorganize into a new sponge).

Features :

- Web framework [gin](https://github.com/gin-gonic/gin)
- RPC framework [grpc](https://github.com/grpc/grpc-go)
- Configuration file parsing [viper](https://github.com/spf13/viper)
- Configuration Center [nacos](https://github.com/alibaba/nacos)
- Logging [zap](https://go.uber.org/zap)
- Database component [gorm](https://gorm.io/gorm)
- Caching component [go-redis](https://github.com/go-redis/redis) [ristretto](github.com/dgraph-io/ristretto)
- Documentation [swagger](https://github.com/swaggo/swag)
- Authorization [jwt](https://github.com/golang-jwt/jwt)
- Validator [validator](https://github.com/go-playground/validator)
- Rate limiter [aegis](https://github.com/go-kratos/aegis)
- Circuit Breaker [aegis](https://github.com/go-kratos/aegis)
- Tracking [opentelemetry](https://go.opentelemetry.io/otel)
- Monitoring [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana)
- Service registration and discovery [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/)
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

```bash
go install github.com/zhufuyi/sponge/cmd/sponge@latest

sponge update
```

<br>

### Quickly create a http project

#### Creating a new http server

**(1) Generate http server code**

```bash
sponge http \
  --module-name=account \
  --server-name=account \
  --project-name=account \
  --repo-addr=zhufuyi \
  --db-dsn=root:123456@(127.0.0.1:3306)/test \
  --db-table=student
```

**(2) If using the default configuration, skip this step. Modify the configuration file configs/<server name>.yml**

- If the field `cacheType` is `redis`, the `redis` address must be set.
- If the field `enableTracing` is true, the `jaeger` address must be set.
- If the field `enableRegistryDiscovery` is true, the configuration of the corresponding type of `registryDiscoveryType` must be set.

**(3) Generate swagger documentation**

> make docs

**(4) Start up the server**

Way 1: Run locally in the binary

> make run

Copy `http://localhost:8080/swagger/index.html` to your browser and test the api interface.

Way 2: Run in docker. Prerequisite: `docker` and `docker-compose` are already installed.

```bash
# Build the docker image
make docker-image

# Start the service
make deploy-docker

# Check the status of the service, if it is healthy, it started successfully
cd deployments/docker-compose
docker-compose ps
```

Way 3: Run in k8s. Prerequisite: `docker` and `kubectl` are already installed.

```bash
# Build the image
make image-build REPO_HOST=zhufuyi TAG=latest

# Push the image to the remote image repository and delete the local image after a successful upload
make image-push REPO_HOST=zhufuyi TAG=latest

# Deploy to k8s
kubectl apply -f deployments/kubernetes/*namespace.yml
kubectl apply -f deployments/kubernetes/
make deploy-k8s

# Check the status of the service
kubectl get -f account-deployment.yml
```

You can also use Jenkins to automatically build deployments to k8s.

<br>

#### Creating a new handler

add a new **handler** to an existing http server.

```bash
sponge handler \
  --module-name=account \
  --db-dsn=root:123456@(127.0.0.1:3306)/test \
  --db-table=teacher \
  --out=./account
```

Start up the server

> make docs && make run

Copy `http://localhost:8080/swagger/index.html` to your browser and test the api interface.

<br>

### Quick create a grpc project

#### Creating a new grpc server

**(1) Generate grpc server code**

```bash
sponge grpc \
  --module-name=account \
  --server-name=account \
  --project-name=account \
  --repo-addr=zhufuyi \
  --db-dsn=root:123456@(127.0.0.1:3306)/test \
  --db-table=student
```

**(2) If using the default configuration, skip this step. Modify the configuration file configs/<server name>.yml**

- If the field `cacheType` is `redis`, the `redis` address must be set.
- If the field `enableTracing` is true, the `jaeger` address must be set.
- If the field `enableRegistryDiscovery` is true, the configuration of the corresponding type of `registryDiscoveryType` must be set.

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

# Deploy to k8s
kubectl apply -f deployments/kubernetes/*namespace.yml
kubectl apply -f deployments/kubernetes/
make deploy-k8s

# Check the status of the service
kubectl get -f account-deployment.yml
```

You can also use Jenkins to automatically build deployments to k8s.

<br>

#### Creating a new service

add a new **service** to an existing grpc server.

```bash
sponge service \
  --module-name=account \
  --server-name=account \
  --db-dsn=root:123456@(127.0.0.1:3306)/test \
  --db-table=teacher \
  --out=./account
```

Start up the server

> make proto && make run

Use IDE to open the file `internal/service/<table name>_client_test.go` to test the api interface of grpc.

<br>

## Give a Star! ⭐

If you like it, and it's of use for you, please give a star, thanks!

<br>

## License

See the [LICENSE](LICENSE) file for licensing information.
