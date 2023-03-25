#!/bin/bash

serverName="serverNameExample_mixExample"

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

go build -o ${binaryFile} cmd/${serverName}/main.go
checkResult $?

# running server
./${binaryFile}
