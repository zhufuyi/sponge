#!/bin/bash

serverName="serverNameExample_mixExample"

binaryFile="cmd/${serverName}/${serverName}"

osType=$(uname -s)
if [ "${osType%%_*}"x = "MINGW64"x ];then
    binaryFile="${binaryFile}.exe"
fi

if [ -f "${binaryFile}" ] ;then
     rm "${binaryFile}"
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
