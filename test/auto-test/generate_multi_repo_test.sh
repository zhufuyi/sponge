#!/bin/bash

mysqlDsn="root:123456@(192.168.3.37:3306)/account"
mysqlTable1="user_example"
mysqlTable2="user"

postgresqlDsn="root:123456@(192.168.3.37:5432)/account"
postgresqlTable1="user_example"
postgresqlTable2="user"

sqliteDsn="../sql/sqlite/sponge.db"
sqliteTable1="user_example"
sqliteTable2="user"

mongodbDsn="root:123456@(192.168.3.37:27017)/account"
mongodbCollection1="user_example"
mongodbCollection2="user"

isOnlyGenerateCode="false"

colorCyan='\e[1;36m'
markEnd='\e[0m'

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function stopService() {
  local name=$1
  if [ "$name" == "" ]; then
    echo "name cannot be empty"
    return 1
  fi

  local processMark="./cmd/$name"
  pid=$(ps -ef | grep "${processMark}" | grep -v grep | awk '{print $2}')
  if [ "${pid}" != "" ]; then
      kill -9 ${pid}
      return 0
  fi

  return 1
}

function checkServiceStarted() {
  local name=$1
  if [ "$name" == "" ]; then
    echo "name cannot be empty"
    return 1
  fi

  local processMark="./cmd/$name"
  local timeCount=0
  # waiting for service to start
  while true; do
    sleep 1
    pid=$(ps -ef | grep "${processMark}" | grep -v grep | awk '{print $2}')
    if [ "${pid}" != "" ]; then
        break
    fi
    (( timeCount++ ))
    if (( timeCount >= 30 )); then
      echo "service startup timeout"
      return 1
    fi
  done
}

function runningHTTPService() {
  local name=$1
  if [ "$name"x = x ];then
    echo "server name cannot be empty"
    return 1
  fi

  make docs
  checkResult $?
  echo "startup service $name"
  make run &
  checkServiceStarted $name
  checkResult $?
  sleep 1
  stopService $name
  checkResult $?
}

function runningProtoService() {
  local name=$1
  if [ "$name"x = x ];then
    echo "server name cannot be empty"
    return 1
  fi

  make proto
  checkResult $?
  echo -e "startup service $name"
  make run &
  checkServiceStarted $name
  checkResult $?
  sleep 1
  stopService $name
  checkResult $?
}

# -------------------------------------------------------------------------------------------

function generate_http_mysql() {
  local serverName="user"
  local outDir="./http-mysql"
  echo "start generating http-mysql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --out=$outDir ${markEnd}"
    sponge web handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_http_postgresql() {
  local serverName="user"
  local outDir="./http-postgresql"
  echo "start generating http-postgresql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --out=$outDir ${markEnd}"
    sponge web handler --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_http_sqlite() {
  local serverName="user"
  local outDir="./http-sqlite"
  echo "start generating http-sqlite service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --out=$outDir ${markEnd}"
    sponge web handler --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --out=$outDir
    checkResult $?

    sed -E -i 's/\\\\sql\\\\/\\\\\.\.\\\\\sql\\\\/g' ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_http_mongodb() {
  local serverName="user"
  local outDir="./http-mongodb"
  echo "start generating http-mongodb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --out=$outDir ${markEnd}"
    sponge web handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

# ---------------------------------------------------------------

function generate_grpc_mysql() {
  local serverName="user"
  local outDir="./grpc-mysql"
  echo "start generating grpc-mysql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_grpc_postgresql() {
  local serverName="user"
  local outDir="./grpc-postgresql"
  echo "start generating grpc-postgresql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --out=$outDir ${markEnd}"
    sponge micro service --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_grpc_sqlite() {
  local serverName="user"
  local outDir="./grpc-sqlite"
  echo "start generating grpc-sqlite service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --out=$outDir ${markEnd}"
    sponge micro service --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --out=$outDir
    checkResult $?

    sed -E -i 's/\\\\sql\\\\/\\\\\.\.\\\\\sql\\\\/g' ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_grpc_mongodb() {
  local serverName="user"
  local outDir="./grpc-mongodb"
  echo "start generating grpc-mongodb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=user --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

# ---------------------------------------------------------------

function generate_http_pb_mysql() {
  local serverName="user"
  local outDir="./http-pb-mysql"
  echo "start generating http-pb-mysql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user.proto --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user.proto --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler-pb --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir ${markEnd}"
    sponge web handler-pb --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir
    checkResult $?

    mysqlDsnTmp=$(echo "$mysqlDsn" | sed -E 's/\(/\\\(/g' | sed -E 's/\)/\\\)/g' | sed -E 's/\//\\\//g')
    sed -E -i "s/root:123456@\(192.168.3.37:3306\)\/account/${mysqlDsnTmp}/g" ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mysql
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_http_pb_mongodb() {
  local serverName="user"
  local outDir="./http-pb-mongodb"
  echo "start generating http-pb-mongodb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user.proto --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user.proto --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler-pb --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --out=$outDir ${markEnd}"
    sponge web handler-pb --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --out=$outDir
    checkResult $?

    sed -E -i 's/\"mysql\"/\"mongodb\"/g' ${outDir}/configs/${serverName}.yml
    sed -E -i 's/mysql:/mongodb:/g' ${outDir}/configs/${serverName}.yml
    mongodbDsnTmp=$(echo "$mongodbDsn" | sed -E 's/\(/\\\(/g' | sed -E 's/\)/\\\)/g' | sed -E 's/\//\\\//g')
    sed -E -i "s/root:123456@\(192.168.3.37:3306\)\/account/${mongodbDsnTmp}/g" ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mongodb
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

# ---------------------------------------------------------------

function generate_grpc_pb_mysql() {
  local serverName="user"
  local outDir="./grpc-pb-mysql"
  echo "start generating grpc-pb-mysql service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user2.proto --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user2.proto --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --out=$outDir
    checkResult $?

    mysqlDsnTmp=$(echo "$mysqlDsn" | sed -E 's/\(/\\\(/g' | sed -E 's/\)/\\\)/g' | sed -E 's/\//\\\//g')
    sed -E -i "s/root:123456@\(192.168.3.37:3306\)\/account/${mysqlDsnTmp}/g" ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mysql
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}


function generate_grpc_pb_mongodb() {
  local serverName="user"
  local outDir="./grpc-pb-mongodb"
  echo "start generating grpc-pb-mongodb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user2.proto --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/user2.proto --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --out=$outDir
    checkResult $?

    sed -E -i 's/\"mysql\"/\"mongodb\"/g' ${outDir}/configs/${serverName}.yml
    sed -E -i 's/mysql:/mongodb:/g' ${outDir}/configs/${serverName}.yml
    mongodbDsnTmp=$(echo "$mongodbDsn" | sed -E 's/\(/\\\(/g' | sed -E 's/\)/\\\)/g' | sed -E 's/\//\\\//g')
    sed -E -i "s/root:123456@\(192.168.3.37:3306\)\/account/${mongodbDsnTmp}/g" ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  make patch TYPE=types-pb
  make patch TYPE=init-mongodb
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

# ---------------------------------------------------------------

function generate_http_pb_mixed() {
  local serverName="user"
  local outDir="./http-pb-mixed"
  echo "start generating http-pb-mixed service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/mixed.proto --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/mixed.proto --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_grpc_pb_mixed() {
  local serverName="user"
  local outDir="./grpc-pb-mixed"
  echo "start generating grpc-pb-mixed service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/mixed.proto --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=user --project-name=edusys --protobuf-file=./files/mixed.proto --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

# ---------------------------------------------------------------

function generate_grpc_gw_pb_mixed() {
  local serverName="user_gw"
  local outDir="./grpc-gw-pb-mixed"
  echo "start generating grpc-gw-pb-mixed service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=./files/mixed.proto --out=$outDir ${markEnd}"
    sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=./files/mixed.proto --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
}

function generate_grpc_gw_pb() {
  local serverName="user_gw"
  local outDir="./grpc-gw-pb"
  echo "start generating grpc-gw-pb service code"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=./files/user_gw.proto --out=$outDir ${markEnd}"
    sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=./files/user_gw.proto --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro rpc-conn --rpc-server-name=user --out=$outDir ${markEnd}"
    sponge micro rpc-conn --rpc-server-name=user --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  make copy-proto SERVER=../grpc-mysql
  checkResult $?
  runningProtoService $serverName
  checkResult $?
  sleep 1
  cd -

  echo -e "\n\n--------------------- $outDir test passed ---------------------\n\n"
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

  generate_http_pb_mixed
  generate_grpc_pb_mixed

  generate_grpc_gw_pb_mixed
  generate_grpc_gw_pb
}

main
