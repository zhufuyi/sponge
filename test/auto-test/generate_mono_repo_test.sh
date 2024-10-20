#!/bin/bash

projectName="edusys"

mysqlDsn="root:123456@(192.168.3.37:3306)/account"
mysqlTable1="user_example"
mysqlTable2="user"
mysqlTable3="user_account"

postgresqlDsn="root:123456@(192.168.3.37:5432)/account"
postgresqlTable1="user_example"
postgresqlTable2="user"

sqliteDsn="../sql/sqlite/sponge.db"
sqliteTable1="user_example"
sqliteTable2="user"

mongodbDsn="root:123456@(192.168.3.37:27017)/account"
mongodbCollection1="user_example"
mongodbCollection2="user"
mongodbCollection3="userInfo"

colorCyan='\e[1;36m'
colorGreen='\e[1;32m'
colorRed='\e[1;31m'
markEnd='\e[0m'

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

  make patch TYPE=types-pb
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
  local serverName="mono_01_http_mysql"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --embed=true --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --embed=true --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --embed=true --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_02_http_postgresql"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_03_http_sqlite"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_04_http_mongodb"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web http --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_05_grpc_mysql"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_06_grpc_postgresql"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=postgresql --db-dsn=$postgresqlDsn --db-table=$postgresqlTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_07_grpc_sqlite"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=sqlite --db-dsn=$sqliteDsn --db-table=$sqliteTable2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_08_grpc_mongodb"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro rpc --server-name=$serverName --module-name=$projectName --project-name=$projectName --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection2 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_09_http_pb_mysql"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler-pb --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler-pb --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_10_http_pb_mongodb"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge web handler-pb --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge web handler-pb --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=$mongodbCollection1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_11_grpc_pb_mysql"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user2.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user2.proto --suited-mono-repo=true --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_12_grpc_pb_mongodb"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user2.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user2.proto --suited-mono-repo=true --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_13_grpc_http_pb_mysql"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro grpc-http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro grpc-http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service-handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service-handler --db-driver=mysql --db-dsn=$mysqlDsn --db-table=$mysqlTable1 --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_14_grpc_http_pb_mongodb"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro grpc-http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro grpc-http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user.proto --suited-mono-repo=true --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro service-handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --suited-mono-repo=true --extended-api=$isExtended --out=$outDir ${markEnd}"
    sponge micro service-handler --db-driver=mongodb --db-dsn=$mongodbDsn --db-table=user_example --suited-mono-repo=true --extended-api=$isExtended --out=$outDir
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
  local serverName="mono_15_http_pb_mixed"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge web http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/mixed.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge web http-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/mixed.proto --suited-mono-repo=true --out=$outDir
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
  local serverName="mono_16_grpc_pb_mixed"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/mixed.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro rpc-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/mixed.proto --suited-mono-repo=true --out=$outDir
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
  local serverName="mono_17_grpc_gw_pb_mixed"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-gw-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/mixed.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro rpc-gw-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/mixed.proto --suited-mono-repo=true --out=$outDir
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
  local serverName="mono_18_grpc_gw_pb"
  local outDir="$serverName"
  echo -e "\n\n"
  echo -e "create service code to directory $outDir"
  if [ -d "${outDir}" ]; then
    echo -e "$outDir already exists\n\n"
  else
    echo -e "\n${colorCyan}sponge micro rpc-gw-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user_gw.proto --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro rpc-gw-pb --server-name=$serverName --module-name=$projectName --project-name=$projectName --protobuf-file=./files/user_gw.proto --suited-mono-repo=true --out=$outDir
    checkResult $?

    echo -e "\n${colorCyan}sponge micro rpc-conn --rpc-server-name=user --suited-mono-repo=true --out=$outDir ${markEnd}"
    sponge micro rpc-conn --rpc-server-name=user --suited-mono-repo=true --out=$outDir
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
