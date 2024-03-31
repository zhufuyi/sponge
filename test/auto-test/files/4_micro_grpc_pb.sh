#!/bin/bash

grpcServiceName="user"
grpcDir="4_micro_grpc_pb_${grpcServiceName}"

colorCyan='\033[1;36m'
markEnd='\033[0m'
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
  echo -e "${colorCyan}go test -run Test_service_user_methods/Register ${markEnd}"
  go test -run Test_service_user_methods/Register
  checkErrCount $?

  cd -
  echo -e "\n--------------------- the test is over, error result: $errCount ---------------------\n"
  stopService $grpcServiceName
}

if [ -d "${grpcDir}" ]; then
  echo "service ${grpcDir} already exists"
else
  echo "create service ${grpcDir}"
  echo -e "${colorCyan}sponge micro rpc-pb --module-name=${grpcServiceName} --server-name=${grpcServiceName} --project-name=grpcpbdemo --protobuf-file=./files/user2.proto --out=./${grpcDir} ${markEnd}"
  sponge micro rpc-pb --module-name=${grpcServiceName} --server-name=${grpcServiceName} --project-name=grpcpbdemo --protobuf-file=./files/user2.proto --out=./${grpcDir}
  checkResult $?
fi

cd ${grpcDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

echo "replace the sample template code"
replaceCode ../files/micro_grpc_pb_content ./internal/service/user2.go
checkResult $?
replaceCode ../files/micro_grpc_pb_content ./internal/service/user2_client_test.go
checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
