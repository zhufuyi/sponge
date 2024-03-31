#!/bin/bash

webServiceName="user"
webDir="3_web_gin_pb_${webServiceName}"

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
  checkServiceStarted $webServiceName
  sleep 1
  echo "--------------------- start testing ---------------------"

  echo -e "\n\n"
  echo -e "${colorCyan}curl -X POST http://localhost:8080/api/v1/auth/register -H \"Content-Type: application/json\" -d {\"email\":\"foo@bar.com\",\"password\":\"123456\"} ${markEnd}"
  curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  checkErrCount $?

  echo -e "\n\n"
  echo -e "${colorCyan}curl -X POST http://localhost:8080/api/v1/auth/register  -H "Content-Type: application/json" -H \"X-Request-Id: qaz12wx3ed4\" -d {\"email\":\"foo@bar.com\",\"password\":\"123456\"} ${markEnd}"
  curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -H "X-Request-Id: qaz12wx3ed4" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  checkErrCount $?

  echo -e "\n--------------------- the test is over, error result: $errCount ---------------------\n"
  stopService $webServiceName
}

if [ -d "${webDir}" ]; then
  echo "service ${webDir} already exists"
else
  echo "create service ${webDir}"
  echo -e "${colorCyan}sponge web http-pb --module-name=${webServiceName} --server-name=${webServiceName} --project-name=ginpbdemo --protobuf-file=./files/user.proto --out=./${webDir} ${markEnd}"
  sponge web http-pb --module-name=${webServiceName} --server-name=${webServiceName} --project-name=ginpbdemo --protobuf-file=./files/user.proto --out=./${webDir}
  checkResult $?
fi

cd ${webDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

echo "replace the sample template code"
replaceCode ../files/web_gin_pb_content ./internal/handler/user.go
checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
