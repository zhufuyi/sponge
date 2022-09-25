SHELL := /bin/bash

PROJECT_NAME := "github.com/zhufuyi/sponge"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)


.PHONY: init
# init env
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
# make ci-lint
ci-lint:
	golangci-lint run ./...


.PHONY: build
# make build, Build the binary file
build:
	@cd cmd/sponge && go build
	@echo "build finished, binary file in path 'cmd/sponge'"


.PHONY: run
# make run, run app
run:
	@bash scripts/run.sh


.PHONY: dep
# make dep Get the dependencies
dep:
	@go mod download


.PHONY: fmt
# make fmt
fmt:
	@gofmt -s -w .


.PHONY: test
# make test
test:
	go test -short ${PKG_LIST}


.PHONY: cover
# make cover
cover:
	go test -short -coverprofile coverage.out -covermode=atomic ${PKG_LIST}
	go tool cover -html=coverage.out


.PHONY: docker
# generate docker image
docker:
	docker build -t sponge:latest -f build/Dockerfile


.PHONY: clean
# make clean
clean:
	@-rm -vrf sponge
	@-rm -vrf cover.out
	@-rm -vrf coverage.txt
	@go mod tidy
	@echo "clean finished"


.PHONY: docs
# gen swagger doc
docs:
	@swag init -g cmd/sponge/main.go
	@echo "see docs by: http://localhost:8080/swagger/index.html"


.PHONY: graph
# make graph 生成交互式的可视化Go程序调用图，生成完毕后会在浏览器自动打开
graph:
	@echo "generating graph ......"
	@go-callvis github.com/zhufuyi/sponge


.PHONY: mockgen
# make mockgen gen mock file
mockgen:
	cd ./internal &&  for file in `egrep -rnl "type.*?interface" ./repository | grep -v "_test" `; do \
		echo $$file ; \
		cd .. && mockgen -destination="./internal/mock/$$file" -source="./internal/$$file" && cd ./internal ; \
	done


.PHONY: proto
# generate proto struct only
proto:
	@bash scripts/protoc.sh


.PHONY: proto-doc
# generate proto doc
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
