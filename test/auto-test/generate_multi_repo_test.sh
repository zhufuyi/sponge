#!/bin/bash

mysqlDsn="root:123456@(192.168.3.37:3306)/account"
mysqlTable1="user"
mysqlTable2="user_order"
mysqlTable3="user_str"

postgresqlDsn="root:123456@(192.168.3.37:5432)/account"
postgresqlTable1="user"
postgresqlTable2="user_order"
postgresqlTable3="user_str"

sqliteDsn="../sql/sqlite/sponge.db"
sqliteTable1="user"
sqliteTable2="user_order"
sqliteTable3="user_str"

mongodbDsn="root:123456@(192.168.3.37:27017)/account"
mongodbCollection1="user"
mongodbCollection2="user_order"
mongodbCollection3="people"

colorCyan='\e[1;36m'
colorGreen='\e[1;32m'
colorRed='\e[1;31m'
markEnd='\e[0m'

httpProtobufFile="files/user_http.proto"
grpcProtobufFile="files/user_rpc.proto"
mixProtobufFile="files/user_hybrid.proto"
grpcGwProtobufFile="files/user_gw.proto"
mixedProtobufFile="files/mixed.proto"

isOnlyGenerateCode="false"

isExtended="false"
if [ "$1" == "true" ]; then
  isExtended="true"
else
  isExtended="false"
fi

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function printTestResult() {
  local errCount=$1
  local serverDir=$2
  if [ ${errCount} -eq 0 ]; then
    echo -e "\n\n${colorGreen}--------------------- [${serverDir}] test result: passed ---------------------${markEnd}\n"
  else
    echo -e "\n\n${colorRed}--------------------- [${serverDir}] test result: failed ---------------------${markEnd}\n"
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
  local serverDir=$2
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
        printTestResult 0 $serverDir
        break
    fi
    (( timeCount++ ))
    if (( timeCount >= 20 )); then
      printTestResult 1 $serverDir
      return 1
    fi
  done
}

function runningHTTPService() {
  local name=$1
  local serverDir=$2
  if [ "$name"x = x ];then
    echo "server name cannot be empty"
    return 1
  fi

  make docs
  checkResult $?
  echo "startup service $name"
  make run &
  checkServiceStarted $name $serverDir
  checkResult $?
  sleep 1
  stopService $name
  checkResult $?
}

function runningProtoService() {
  local name=$1
  local serverDir=$2
  if [ "$name"x = x ];then
    echo "server name cannot be empty"
    return 1
  fi

  make proto
  checkResult $?
  echo -e "startup service $name"
  make run &
  checkServiceStarted $name $serverDir
  checkResult $?
  sleep 1
  stopService $name
  checkResult $?
}

# -------------------------------------------------------------------------------------------

function generate_http_mysql() {
  local serverName="user"
  local outDir="multi-01-http-mysql"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_http_postgresql() {
  local serverName="user"
  local outDir="multi-02-http-postgresql"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2,$postgresqlTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2,$postgresqlTable3 --extended-api=$isExtended --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_http_sqlite() {
  local serverName="user"
  local outDir="multi-03-http-sqlite"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2,$sqliteTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2,$sqliteTable3 --extended-api=$isExtended --out=$outDir
    checkResult $?

    sed -E -i 's/\\\\sql\\\\/\\\\\.\.\\\\\sql\\\\/g' ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_http_mongodb() {
  local serverName="user"
  local outDir="multi-04-http-mongodb"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningHTTPService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

# ---------------------------------------------------------------

function generate_grpc_mysql() {
  local serverName="user"
  local outDir="multi-05-grpc-mysql"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_grpc_postgresql() {
  local serverName="user"
  local outDir="multi-06-grpc-postgresql"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2,$postgresqlTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2,$postgresqlTable3 --extended-api=$isExtended --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_grpc_sqlite() {
  local serverName="user"
  local outDir="multi-07-grpc-sqlite"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2,$sqliteTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2,$sqliteTable3 --extended-api=$isExtended --out=$outDir
    checkResult $?

    sed -E -i 's/\\\\sql\\\\/\\\\\.\.\\\\\sql\\\\/g' ${outDir}/configs/${serverName}.yml
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_grpc_mongodb() {
  local serverName="user"
  local outDir="multi-08-grpc-mongodb"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$serverName --project-name=edusys --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

# ---------------------------------------------------------------

function generate_http_pb_mysql() {
  local serverName="user"
  local outDir="multi-09-http-pb-mysql"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${httpProtobufFile} --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${httpProtobufFile} --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler-pb--db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended  --out=$outDir ${markEnd}"
    sponge web handler-pb --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler-pb --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler-pb --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir
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
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_http_pb_mongodb() {
  local serverName="user"
  local outDir="multi-10-http-pb-mongodb"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${httpProtobufFile} --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${httpProtobufFile} --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler-pb --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1,$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler-pb --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1,$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir
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
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

# ---------------------------------------------------------------

function generate_grpc_pb_mysql() {
  local serverName="user"
  local outDir="multi-11-grpc-pb-mysql"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${grpcProtobufFile} --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${grpcProtobufFile} --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir
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
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_grpc_pb_mongodb() {
  local serverName="user"
  local outDir="multi-12-grpc-pb-mongodb"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${grpcProtobufFile} --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${grpcProtobufFile} --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1,$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1,$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir
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
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

# ---------------------------------------------------------------

function generate_grpc_http_pb_mysql() {
  local serverName="user"
  local outDir="multi-13-grpc-http-pb-mysql"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro grpc-http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixProtobufFile} --out=$outDir ${markEnd}"
    sponge micro grpc-http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixProtobufFile} --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service-handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service-handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service-handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service-handler  --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2,$mysqlTable3 --extended-api=$isExtended --out=$outDir
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
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_grpc_http_pb_mongodb() {
  local serverName="user"
  local outDir="multi-14-grpc-http-pb-mongodb"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro grpc-http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixProtobufFile} --out=$outDir ${markEnd}"
    sponge micro grpc-http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixProtobufFile} --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service-handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1,$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service-handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1,$mongodbCollection2,$mongodbCollection3 --extended-api=$isExtended --out=$outDir
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
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}
# ---------------------------------------------------------------

function generate_http_pb_mixed() {
  local serverName="user"
  local outDir="multi-15-http-pb-mixed"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixedProtobufFile} --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixedProtobufFile} --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_grpc_pb_mixed() {
  local serverName="user"
  local outDir="multi-16-grpc-pb-mixed"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixedProtobufFile} --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=$serverName --project-name=edusys --protobuf-file=${mixedProtobufFile} --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

# ---------------------------------------------------------------

function generate_grpc_gw_pb_mixed() {
  local serverName="user_gw"
  local outDir="multi-17-grpc-gw-pb-mixed"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=${mixedProtobufFile} --out=$outDir ${markEnd}"
    sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=${mixedProtobufFile} --out=$outDir
    checkResult $?
  fi

  if [ "$isOnlyGenerateCode" == "true" ]; then
    echo -e "\n\n"
    return
  fi

  cd $outDir
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
}

function generate_grpc_gw_pb() {
  local serverName="user_gw"
  local outDir="multi-18-grpc-gw-pb"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=${grpcGwProtobufFile} --out=$outDir ${markEnd}"
    sponge micro rpc-gw-pb --server-name=$serverName --module-name=edusys --project-name=edusys --protobuf-file=${grpcGwProtobufFile} --out=$outDir
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
  make copy-proto SERVER=../multi-05-grpc-mysql
  checkResult $?
  runningProtoService $serverName $outDir
  checkResult $?
  sleep 1
  cd -
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

  generate_grpc_http_pb_mysql
  generate_grpc_http_pb_mongodb

  generate_http_pb_mixed
  generate_grpc_pb_mixed

  generate_grpc_gw_pb_mixed
  generate_grpc_gw_pb
}

main
