SHELL := /bin/bash

PROJECT_NAME := "github.com/zhufuyi/sponge"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v /api/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)


.PHONY: install
# installation of dependent tools
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/envoyproxy/protoc-gen-validate@v0.6.7
	go install github.com/srikrsna/protoc-gen-gotag@v0.6.2
	go install github.com/mohuishou/protoc-gen-go-gin@v0.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.0
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
	go install github.com/swaggo/swag/cmd/swag@v1.8.6
	go install github.com/ofabry/go-callvis@v0.6.1
	go install golang.org/x/pkgsite/cmd/pkgsite@latest


.PHONY: mod
# add missing and remove unused modules
mod:
	go mod tidy


.PHONY: fmt
# go format *.go files
fmt:
	gofmt -s -w .


.PHONY: ci-lint
# check the code specification against the rules in the .golangci.yml file
ci-lint: fmt
	golangci-lint run ./...


.PHONY: dep
# download dependencies to the directory vendor
dep:
	go mod download


.PHONY: test
# go test *_test.go files, the parameter -count=1 means that caching is disabled
test:
	go test -count=1 -short ${PKG_LIST}


.PHONY: cover
# generate test coverage
cover:
	go test -short -coverprofile=cover.out -covermode=atomic ${PKG_LIST}
	go tool cover -html=cover.out


.PHONY: docs
# generate swagger docs, the host address can be changed via parameters, e.g. make docs HOST=192.168.3.37
docs: mod fmt
	@bash scripts/swag-docs.sh $(HOST)


.PHONY: graph
# generate interactive visual function dependency graphs
graph:
	@echo "generating graph ......"
	@cp -f cmd/serverNameExample/main.go .
	go-callvis -skipbrowser -format=svg -nostd -file=serverNameExample github.com/zhufuyi/sponge
	@rm -f main.go serverNameExample.gv


.PHONY: proto
# generate *.pb.go codes from *.proto files
proto: mod fmt
	@bash scripts/protoc.sh


.PHONY: proto-doc
# generate doc from *.proto files
proto-doc:
	@bash scripts/proto-doc.sh


.PHONY: build
# build serverNameExample for linux amd64 binary
build:
	@echo "building 'serverNameExample', binary file will output to 'cmd/serverNameExample'"
	@cd cmd/serverNameExample && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct go build -gcflags "all=-N -l"

# delete the templates code start
.PHONY: build-sponge
# build sponge for linux amd64 binary
build-sponge:
	@echo "building 'sponge', binary file will output to 'cmd/sponge'"
	@cd cmd/sponge && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct go build
# delete the templates code end

.PHONY: run
# run server
run:
	@bash scripts/run.sh


.PHONY: docker-image
# build docker image, use binary files to build
docker-image: build
	@bash scripts/grpc_health_probe.sh
	@mv -f cmd/serverNameExample/serverNameExample build/
	@mkdir -p build/configs && cp -f configs/serverNameExample.yml build/configs/
	docker build -t project-name-example.server-name-example:latest build/
	@rm -rf build/serverNameExample build/configs build/grpc_health_probe


.PHONY: image-build
# build docker image with parameters, use binary files to build, e.g. make image-build REPO_HOST=addr TAG=latest
image-build:
	@bash scripts/image-build.sh $(REPO_HOST) $(TAG)


.PHONY: image-build2
# build docker image with parameters, phase II build, e.g. make image-build2 REPO_HOST=addr TAG=latest
image-build2:
	@bash scripts/image-build2.sh $(REPO_HOST) $(TAG)


.PHONY: image-push
# push docker image to remote repositories, e.g. make image-push REPO_HOST=addr TAG=latest
image-push:
	@bash scripts/image-push.sh $(REPO_HOST) $(TAG)


.PHONY: deploy-k8s
# deploy service to k8s
deploy-k8s:
	@bash scripts/deploy-k8s.sh


.PHONY: deploy-docker
# deploy service to docker
deploy-docker:
	@bash scripts/deploy-docker.sh


.PHONY: clean
# clean binary file, cover.out, redundant dependency packages
clean:
	@rm -vrf cmd/serverNameExample/serverNameExample
	@rm -vrf cover.out
	go mod tidy
	@echo "clean finished"


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
