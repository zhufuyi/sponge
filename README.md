## sponge

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/sponge)](https://goreportcard.com/report/github.com/zhufuyi/sponge)  
[![codecov](https://codecov.io/gh/zhufuyi/sponge/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/sponge)  
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/sponge.svg)](https://pkg.go.dev/github.com/zhufuyi/sponge)  
[![Go](https://github.com/zhufuyi/sponge/workflows/Go/badge.svg?branch=main)](https://github.com/zhufuyi/sponge/actions)  
[![License: MIT](https://img.shields.io/github/license/zhufuyi/sponge)](https://img.shields.io/github/license/zhufuyi/sponge)

**sponge** is a microservices framework that supports automatic generation of http service code, grpc service code and CICD execution scripts. Generate config, ecode, model, dao, handler, router, http, proto, service, grpc independent sub-modules. The low-coupling module code is combined into a complete service (similar to how a broken sponge cell automatically reorganises itself into a new sponge) and only needs to be populated with business code.

Features :

- web framework [gin](https://github.com/gin-gonic/gin)
- rpc framework [grpc](https://github.com/grpc/grpc-go)
- Configuration file parsing [viper](https://github.com/spf13/viper)
- logging [zap](go.uber.org/zap)
- Database component [gorm](gorm.io/gorm)
- Caching component [go-redis](github.com/go-redis/redis)
- Documentation [swagger](github.com/swaggo/swag)
- Validator [validator](github.com/go-playground/validator)
- Link tracking [opentelemetry](go.opentelemetry.io/otel)
- Metrics collection [prometheus](github.com/prometheus/client_golang/prometheus)
- ratelimiter](golang.org/x/time/rate)
- fuse [hystrix](github.com/afex/hystrix-go)
- Configuration Center [nacos](https://github.com/alibaba/nacos)
- Service registration and discovery [etcd](https://github.com/etcd-io/etcd), [consul](https://github.com/hashicorp/consul), [nacos](https://github.com/alibaba/) nacos)
- Package management tools [go modules](https://github.com/golang/go/wiki/Modules)
- Performance analysis [go profile](https://go.dev/blog/pprof)
- Code inspection [golangci-lint](https://github.com/golangci/golangci-lint)
- Continuous Integration CI [jenkins](https://github.com/jenkinsci/jenkins)
- Continuous Deployment CD [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

<br>  

The directory structure follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)。

```bash
.  
├── api               # grpc's proto file and corresponding code  
├── assets          # other assets used with the repository (images, logos, etc.)  
├── build            # Packaging and continuous integration  
├── cmd             # The application's directory  
├── configs         # Directory of configuration files  
├── deployments # IaaS, PaaS, system and container orchestration deployment configurations and templates  
├─ docs              # Design documentation and interface documentation  
├── internal        # Private application and library code  
│ ├── cache       # Business wrapper-based cache  
│ ├── config       # Go struct for config file mapping  
│ ├── dao          # Data access  
│ ├── ecode       # custom business error codes  
│ ├── handler    # Business function implementation for http  
│ ├── model      # Database model  
│ ├── routers    # http routing  
│ ├── server     # service entry, including http and grpc services  
│ ├── service    # Business function implementation for grpc  
├── pkg            # library code that external applications can use  
├── scripts       # Scripts for performing various build, install, analysis, etc. operations  
├── test            # Additional external test applications and test data  
└── third_party # External helpers, forked code and other third party tools  
``` 

<br>  

The development specification follows the [Uber Go Language Coding Specification](https://github.com/uber-go/guide/blob/master/style.md).

<br>  

## Quick start

### Quickly create an http project

#### Creating a new http service

**(1) Generate http service code**

> sponge http --module-name=account --server-name=account --project-name=account --repo-addr=zhufuyi --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=use

Command parameter description.

- --**module-name**: module name of the go.mod file
- --**server-name**: the name of the http service
- --**project-name**: the name of the project to which the http service belongs
- --**repo-addr**: image repository address
- --**db-dsn**: address of the msyql connection
- --**db-table**: the name of the table
- --**embed**: whether to embed gorm.Model (id and time as embedded fields), optional parameter, default is true
- --**out**: generate code to the specified directory, e.g. `$PATH/src/<http service name>`, optional, default code generated in the current directory, note: if the file exists it will be cancel generate codes.


**(2) Modify the configuration file configs/<service name>.yml**

- Modify the redis configuration
- If the field `enableTracing` is true, the jaeger address must be set
- If the field `enableRegistryDiscovery` is true, the etcd address must be set

**(3) Generate swagger documentation**

```bash
# First run, update project dependencies
make mod && make fmt

# Generate swagger documentation
make docs
```

**(4) Start the service**

Way 1: Run locally in the binary

> make run

Copy `http://localhost:8080/swagger/index.html` to your browser and test the add/remove interface.

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

### Quick creation of grpc projects

#### Creating a new grpc service

**(1) Generate grpc service code**

> sponge grpc --module-name=account --server-name=account --project-name=account --repo-addr=zhufuyi --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

Command parameter description.

- --**module-name**: module name of the go.mod file
- --**server-name**: the name of the http service
- --**project-name**: the name of the project to which the http service belongs
- --**repo-addr**: image repository address
- --**db-dsn**: address of the msyql connection
- --**db-table**: the name of the table
- --**embed**: whether to embed gorm.Model (id and time as embedded fields), optional parameter, default is true
- --**out**: generate code to the specified directory, e.g. `$PATH/src/<http service name>`, optional, default code generated in the current directory, note: if the file exists it will be cancel generate codes.

**(2) Modify the configuration file configs/<service name>.yml**

- Modify the redis configuration
- If the field `enableTracing` is true, the jaeger address must be set
- If the field `enableRegistryDiscovery` is true, the etcd address must be set

**(3) Generating grpc code**

```bash
# First run, update project dependencies
make mod && make fmt

# Generate *.pb.go
make proto
```  

**(4) Start the service**

Way 1: Run locally in the binary

> make run

Copy `http://localhost:8080/swagger/index.html` to your browser and test the add/remove interface.

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
