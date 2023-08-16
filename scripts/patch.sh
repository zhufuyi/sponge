#!/bin/bash

patchType=$1
typesPb="types-pb"
mysqlInit="mysql-redis-init"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function patchTypesPb() {
    sponge gen types-pb --out=./
    checkResult $?
}

function patchMysqlAndRedisInit() {
    sponge gen mysql-redis-init --out=./
    checkResult $?
}

if [ "X$patchType" = "X" ];then
    patchTypesPb
    patchMysqlAndRedisInit
elif [  "$patchType" = "$typesPb"  ]; then
    patchTypesPb
elif [ "$patchType" = "$mysqlInit" ]; then
    patchMysqlAndRedisInit
else
    echo "TYPE should be "", $typesPb or $mysqlInit."
    exit 1
fi
