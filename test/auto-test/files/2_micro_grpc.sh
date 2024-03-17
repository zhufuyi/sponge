#!/bin/bash

grpcServiceName="user"
grpcDir="2_micro_grpc_${grpcServiceName}"

mysqlDSN="root:123456@(192.168.3.37:3306)/school"
mysqlTable="teacher"

colorCyan='\e[1;36m'
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
  checkServiceStarted $grpcServiceName
  sleep 1

  echo "--------------------- start testing ---------------------"
  cd internal/service

  echo -e "\n\n"
  echo -e "${colorCyan}go test -run Test_service_${mysqlTable}_methods/GetByID ${markEnd}"
  sed -i "s/Id: 0,/Id: 1,/g" ${mysqlTable}_client_test.go
  go test -run Test_service_${mysqlTable}_methods/GetByID
  checkErrCount $?

  echo -e "\n\n"
  echo -e "${colorCyan}go test -run Test_service_${mysqlTable}_methods/ListByLastID ${markEnd}"
  sed -i "s/Limit:  0,/Limit:  3,/g" ${mysqlTable}_client_test.go
  go test -run Test_service_${mysqlTable}_methods/ListByLastID
  checkErrCount $?

  echo -e "\n--------------------- the test is over, error result: $errCount ---------------------\n"
  cd -
  stopService $grpcServiceName
}

if [ -d "${grpcDir}" ]; then
  echo "service ${grpcDir} already exists"
else
  echo "create service ${grpcDir}"
  echo -e "\n${colorCyan}sponge micro rpc --module-name=${grpcServiceName} --server-name=${grpcServiceName} --project-name=grpcdemo --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${grpcDir} ${markEnd}"
  sponge micro rpc --module-name=${grpcServiceName} --server-name=${grpcServiceName} --project-name=grpcdemo --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${grpcDir}
  checkResult $?
fi

cd ${grpcDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
