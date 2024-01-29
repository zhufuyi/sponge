#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "all=-s -w" -o replaceCode
mv -f replaceCode* $(go env GOBIN)
