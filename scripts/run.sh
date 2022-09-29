#!/bin/bash

serverName="serverNameExample"

binaryFile="cmd/${serverName}/${serverName}"

if [ -f "${serverName}" ] ;then
     rm "${serverName}"
fi

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

sleep 0.2

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${serverName}
go build -o ${binaryFile} cmd/${serverName}/main.go
checkResult $?

# 运行服务
./${binaryFile}
