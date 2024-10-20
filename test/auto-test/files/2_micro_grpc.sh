#!/bin/bash

testServerName="user"
testServerDir="2_micro_grpc_${testServerName}"

mysqlDSN="root:123456@(192.168.3.37:3306)/school"
mysqlTable="teacher"

colorCyan='\e[1;36m'
colorGreen='\e[1;32m'
colorRed='\e[1;31m'
markEnd='\e[0m'
errCount=0

function checkResult() {
  result=$1
  if [ ${result} -ne 0 ]; then
      exit ${result}
  fi
}

function checkErrCount() {
  result=$1
  if [ ${result} -ne 0 ]; then
      ((errCount++))
  fi
}

function printTestResult() {
  if [ ${errCount} -eq 0 ]; then
    echo -e "\n\n${colorGreen}--------------------- [${testServerDir}] test result: passed ---------------------${markEnd}\n"
  else
    echo -e "\n\n${colorRed}--------------------- [${testServerDir}] test result: failed ${errCount} ---------------------${markEnd}\n"
  fi
}

function stopService() {
  local name=$1
  if [ "$name" == "" ]; then
    echo "name cannot be empty"
    exit 1
  fi

  local processMark="./cmd/$name"
  pid=$(ps -ef | grep "${processMark}" | grep -v grep | awk '{print $2}')
  if [ "${pid}" != "" ]; then
      kill -9 ${pid}
  fi
}

function checkServiceStarted() {
  local name=$1
  if [ "$name" == "" ]; then
    echo "name cannot be empty"
    exit 1
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
      exit 1
    fi
  done
}

function testRequest() {
  checkServiceStarted $testServerName
  sleep 1

  cd internal/service
  echo -e "start testing [${testServerName}] api:\n\n"
  echo -e "${colorCyan}go test -run Test_service_${mysqlTable}_methods/GetByID ${markEnd}"
  sed -i "s/Id: 0,/Id: 1,/g" ${mysqlTable}_client_test.go
  go test -run Test_service_${mysqlTable}_methods/GetByID
  checkErrCount $?

  echo -e "\n\n"
  echo -e "${colorCyan}go test -run Test_service_${mysqlTable}_methods/ListByLastID ${markEnd}"
  sed -i "s/Limit:  0,/Limit:  3,/g" ${mysqlTable}_client_test.go
  go test -run Test_service_${mysqlTable}_methods/ListByLastID
  checkErrCount $?

  cd -
  printTestResult
  stopService $testServerName
}

echo -e "\n\n"

if [ -d "${testServerDir}" ]; then
  echo "service ${testServerDir} already exists"
else
  echo "create service ${testServerDir}"
  echo -e "\n${colorCyan}sponge micro rpc --module-name=${testServerName} --server-name=${testServerName} --project-name=grpcdemo --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${testServerDir} ${markEnd}"
  sponge micro rpc --module-name=${testServerName} --server-name=${testServerName} --project-name=grpcdemo --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${testServerDir}
  checkResult $?
fi


cd ${testServerDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

#cp -r ../pkg .
#checkResult $?
#sed -i "s/github.com\/zhufuyi\/sponge\/pkg\/grpc\/benchmark/${testServerName}\/pkg\/grpc\/benchmark/g" internal/service/${mysqlTable}_client_test.go
#checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
