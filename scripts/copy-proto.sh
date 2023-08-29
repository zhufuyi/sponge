#!/bin/bash

serverDir=$1

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

sponge patch copy-proto --server-dir=$serverDir
checkResult $?
