SHELL := /bin/bash

PROJECT_NAME := "github.com/zhufuyi/sponge"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
PROJECT_FILES := $(shell ls)


.PHONY: init
# installation of dependent tools
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.0
	go install github.com/envoyproxy/protoc-gen-validate@v0.6.7
	go install github.com/mohuishou/protoc-gen-go-gin@v0.1.0
	go install github.com/srikrsna/protoc-gen-gotag@v0.6.2
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.1
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
	go install github.com/swaggo/swag/cmd/swag@v1.8.6
	go install github.com/ofabry/go-callvis@v0.6.1


.PHONY: ci-lint
# check the code specification against the rules in the .golangci.yml file
ci-lint:
	golangci-lint run ./...


.PHONY: build
# go build the linux amd64 binary file
build:
	@cd cmd/serverNameExample && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l"
	@echo "build finished, binary file in path 'cmd/serverNameExample'"


.PHONY: run
# run server
run:
	@bash scripts/run.sh


.PHONY: dep
# download dependencies to the directory vendor
dep:
	@go mod download


.PHONY: fmt
# go format *.go files
fmt:
	@gofmt -s -w .


.PHONY: test
# go test *_test.go files
test:
	go test -short ${PKG_LIST}


.PHONY: cover
# generate test coverage
cover:
	go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	go tool cover -html=cover.out


.PHONY: docker
# build docker image
docker:
	@tar zcf serverNameExample.tar.gz ${PROJECT_FILES}
	@mv -f serverNameExample.tar.gz build/
	docker build -t project-name-example/server-name-example:latest build/
	@rm -rf build/serverNameExample.tar.gz


.PHONY: docker-image
# copy the binary file to build the docker image, skip the compile to binary in docker
docker-image: build
	@bash scripts/grpc_health_probe.sh
	@mv -f cmd/serverNameExample/serverNameExample build/ && cp -f /tmp/grpc_health_probe build/
	@mkdir -p build/configs && cp -f configs/serverNameExample.yml build/configs/
	docker build -f build/Dockerfile_cp -t project-name-example/server-name-example:latest build/
	@rm -rf build/serverNameExample build/configs/ build/grpc_health_probe


.PHONY: clean
# clean binary file, cover.out, redundant dependency packages
clean:
	@rm -vrf cmd/serverNameExample/serverNameExample
	@rm -vrf cover.out
	@go mod tidy
	@echo "clean finished"


.PHONY: docs
# generate swagger doc
docs:
	@swag init -g cmd/serverNameExample/main.go
	@echo "see docs by: http://localhost:8080/swagger/index.html"


.PHONY: graph
# generate interactive visual function dependency graphs
graph:
	@echo "generating graph ......"
	@cd cmd/serverNameExample
	@go-callvis -nostd github.com/zhufuyi/sponge


.PHONY: proto
# generate *.pb.go codes from *.proto files
proto:
	@bash scripts/protoc.sh


.PHONY: proto-doc
# generate doc from *.proto files
proto-doc:
	@bash scripts/proto-doc.sh


# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m  %-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all
