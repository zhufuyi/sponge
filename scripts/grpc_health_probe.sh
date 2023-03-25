#!/bin/bash

DOCKERFILE_PATH=$(pwd)/build/grpc_health_probe

GOPROXY=https://goproxy.cn,direct go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.12

cd $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-health-probe@v0.4.12 \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "all=-s -w" -o "${DOCKERFILE_PATH}"
