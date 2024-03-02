#!/bin/bash

patchType=$1
typesPb="types-pb"
initMysql="init-mysql"
initMongodb="init-mongodb"
initTidb="init-tidb"
initPostgresql="init-postgresql"
initSqlite="init-sqlite"

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

function generateInitMongodbCode() {
    sponge patch gen-db-init --db-driver=mongodb --out=./
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

function generateInitSqliteCode() {
    sponge patch gen-db-init --db-driver=sqlite --out=./
    checkResult $?
    importPkg
}

if [  "$patchType" = "$typesPb"  ]; then
    generateTypesPbCode
elif [ "$patchType" = "$initMysql" ]; then
    generateInitMysqlCode
elif [ "$patchType" = "$initMongodb" ]; then
    generateInitMongodbCode
elif [ "$patchType" = "$initTidb" ]; then
    generateInitTidbCode
elif [ "$patchType" = "$initPostgresql" ]; then
    generateInitPostgresqlCode
elif [ "$patchType" = "$initSqlite" ]; then
    generateInitSqliteCode
else
    echo "invalid patch type: '$patchType'"
    echo "supported types: $initMysql, $initMongodb, $initTidb, $initPostgresql, $initSqlite, $typesPb"
    echo "e.g. make patch TYPE=init-mysql"
    echo ""
    exit 1
fi
