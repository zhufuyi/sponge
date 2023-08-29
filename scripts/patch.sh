#!/bin/bash

patchType=$1
typesPb="types-pb"
mysqlInit="mysql-init"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function patchTypesPb() {
    sponge patch gen-types-pb --out=./
    checkResult $?
}

function patchMysqlInit() {
    sponge patch gen-mysql-init --out=./
    checkResult $?
}

if [ "X$patchType" = "X" ];then
    patchTypesPb
    patchMysqlInit
elif [  "$patchType" = "$typesPb"  ]; then
    patchTypesPb
elif [ "$patchType" = "$mysqlInit" ]; then
    patchMysqlInit
else
    echo "TYPE should be "", $typesPb or $mysqlInit."
    exit 1
fi
