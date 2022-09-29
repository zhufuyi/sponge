#!/bin/bash

go env -w GOPROXY=https://goproxy.cn,direct

go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.12

cd $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-health-probe@v0.4.12 \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/grpc_health_probe
