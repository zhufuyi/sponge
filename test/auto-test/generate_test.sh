#!/bin/bash

mysqlDsn="root:123456@(192.168.3.37:3306)/account"
mongodbDsn="root:123456@(192.168.3.37:27017)/account"
postgresqlDsn="root:123456@(192.168.3.37:5432)/account"
sqliteDsn="../sql/sqlite/sponge.db"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function buildCode() {
  serverName=$1
  if [ "$serverName"x = x ];then
    echo "please input serverName"
    exit 1
  fi

  go build cmd/${serverName}/main.go
  checkResult $?
  echo "build successfully"
  rm -f main main.exe
}

function generate_http_mysql() {
  local outDir="./http-mysql"
  echo "generate http-mysql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge web http --module-name=user --server-name=user --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=user_example --out=$outDir
  checkResult $?

  sponge web handler --db-driver=mysql  --db-dsn=$mysqlDsn --db-table=user --out=$outDir
  checkResult $?

  cd $outDir
  make docs
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

function generate_http_postgresql() {
  local outDir="./http-postgresql"
  echo "generate http-postgresql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge web http --module-name=user --server-name=user --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=user_example --embed=false --out=$outDir
  checkResult $?

  cd $outDir
  make docs
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

function generate_http_sqlite() {
  local outDir="./http-sqlite"
  echo "generate http-sqlite service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge web http --module-name=user --server-name=user --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=user_example --embed=false --out=$outDir
  checkResult $?

  cd $outDir
  make docs
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

function generate_http_mongodb() {
  local outDir="./http-mongodb"
  echo "generate http-mongodb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge web http --module-name=user --server-name=user --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --embed=false --out=$outDir
  checkResult $?

  cd $outDir
  make docs
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

# ---------------------------------------------------------------

function generate_grpc_mysql() {
  local outDir="./grpc-mysql"
  echo "generate grpc-mysql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge micro rpc --module-name=user --server-name=user --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=user_example --out=$outDir
  checkResult $?

  cd $outDir
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

function generate_grpc_postgresql() {
  local outDir="./grpc-postgresql"
  echo "generate grpc-postgresql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge micro rpc --module-name=user --server-name=user --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=user_example --embed=false --out=$outDir
  checkResult $?

  cd $outDir
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

function generate_grpc_sqlite() {
  local outDir="./grpc-sqlite"

  echo "generate grpc-sqlite service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge micro rpc --module-name=user --server-name=user --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=user_example --embed=false --out=$outDir
  checkResult $?

  cd $outDir
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

function generate_grpc_mongodb() {
  local outDir="./grpc-mongodb"
  echo "generate grpc-mongodb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge micro rpc --module-name=user --server-name=user --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --embed=false --out=$outDir
  checkResult $?

  cd $outDir
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

# ---------------------------------------------------------------

function generate_http_pb_mysql() {
  local outDir="./http-pb-mysql"
  echo "generate http-pb service code, use database type mysql"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge web http-pb --module-name=user --server-name=user --project-name=edusys --protobuf-file=./files/user.proto --out=$outDir
  checkResult $?

  sponge web handler-pb --db-driver=mysql  --db-dsn=$mysqlDsn --db-table=user_example --out=$outDir
  checkResult $?

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mysql
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

function generate_http_pb_mongodb() {
  local outDir="./http-pb-mongodb"
  echo "generate http-pb service code, use database type mongodb"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge web http-pb --module-name=user --server-name=user --project-name=edusys --protobuf-file=./files/user.proto --out=$outDir
  checkResult $?

  sponge web handler-pb --db-driver=mongodb  --db-dsn=$mongodbDsn --db-table=user_example --out=$outDir
  checkResult $?

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mongodb
  checkResult $?
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}

# ---------------------------------------------------------------

function generate_grpc_pb_mysql() {
  local outDir="./grpc-pb-mysql"
  echo "generate grpc-pb service code, use database type mysql"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge micro rpc-pb --module-name=user --server-name=user --project-name=edusys --protobuf-file=./files/user2.proto --out=$outDir
  checkResult $?

  sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=user_example --out=$outDir
  checkResult $?

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mysql
  checkResult $?
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}


function generate_grpc_pb_mongodb() {
  local outDir="./grpc-pb-mongodb"
  echo "generate mongodb-pb service code, use database type mongodb"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge micro rpc-pb --module-name=user --server-name=user --project-name=edusys --protobuf-file=./files/user2.proto --out=$outDir
  checkResult $?

  sponge micro service --db-driver=mongodb  --db-dsn=$mongodbDsn --db-table=user_example --out=$outDir
  checkResult $?

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mongodb
  make proto
  checkResult $?
  buildCode user
  #make run
  cd - && echo -e "\n\n\n\n"
}


# ---------------------------------------------------------------

function generate_grpc_gw_pb() {
  local outDir="./grpc-gw-pb"
  echo "generate grpc-gw-pb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
    return
  fi

  sponge micro rpc-gw-pb --module-name=edusys --server-name=edusys --project-name=edusys --protobuf-file=./files/user_gw.proto --out=$outDir
  checkResult $?

  sponge micro rpc-conn --rpc-server-name=user --out=$outDir
  checkResult $?

  cd $outDir
  make copy-proto SERVER=../grpc-mysql
  make proto
  checkResult $?
  # modify yaml configuration field "grpcClient"
  buildCode edusys
  #make run
  cd - && echo -e "\n\n\n\n"
}

function main() {
  generate_http_mysql
  generate_http_postgresql
  generate_http_sqlite
  generate_http_mongodb

  generate_grpc_mysql
  generate_grpc_postgresql
  generate_grpc_sqlite
  generate_grpc_mongodb

  generate_http_pb_mysql
  generate_http_pb_mongodb

  generate_grpc_pb_mysql
  generate_grpc_pb_mongodb

  generate_grpc_gw_pb
}

main
