#!/bin/bash

patchType=$1
typesPb="types-pb"
initMysql="init-mysql"
initTidb="init-tidb"
initPostgresql="init-postgresql"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function importPkg() {
    go mod tidy
}

function generateTypesPbCode() {
    sponge patch gen-types-pb --out=./
    checkResult $?
}

function generateInitMysqlCode() {
    sponge patch gen-db-init --db-driver=mysql --out=./
    checkResult $?
    importPkg
}

function generateInitTidbCode() {
    sponge patch gen-db-init --db-driver=tidb --out=./
    checkResult $?
    importPkg
}

function generateInitPostgresqlCode() {
    sponge patch gen-db-init --db-driver=postgresql --out=./
    checkResult $?
    importPkg
}

if [  "$patchType" = "$typesPb"  ]; then
    generateTypesPbCode
elif [ "$patchType" = "$initMysql" ]; then
    generateInitMysqlCode
elif [ "$patchType" = "$initTidb" ]; then
    generateInitTidbCode
elif [ "$patchType" = "$initPostgresql" ]; then
    generateInitPostgresqlCode
else
    echo "invalid patch type: '$patchType'"
    echo "supported types: $initMysql, $initTidb, $initPostgresql, $typesPb"
    echo "e.g. make patch TYPE=init-mysql"
    exit 1
fi
