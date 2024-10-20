#!/bin/bash

testServerName="user"
testServerDir="3_web_http_pb_${testServerName}"

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

  echo -e "start testing [${testServerName}] api:\n\n"
  echo -e "${colorCyan}curl -X POST http://localhost:8080/api/v1/auth/register -H \"Content-Type: application/json\" -d {\"email\":\"foo@bar.com\",\"password\":\"123456\"} ${markEnd}"
  curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  checkErrCount $?

  echo -e "\n\n"
  echo -e "${colorCyan}curl -X POST http://localhost:8080/api/v1/auth/register  -H "Content-Type: application/json" -H \"X-Request-Id: qaz12wx3ed4\" -d {\"email\":\"foo@bar.com\",\"password\":\"123456\"} ${markEnd}"
  curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -H "X-Request-Id: qaz12wx3ed4" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  checkErrCount $?

  printTestResult
  stopService $testServerName
}

echo -e "\n\n"

if [ -d "${testServerDir}" ]; then
  echo "service ${testServerDir} already exists"
else
  echo "create service ${testServerDir}"
  echo -e "${colorCyan}sponge web http-pb --module-name=${testServerName} --server-name=${testServerName} --project-name=ginpbdemo --protobuf-file=./files/user.proto --out=./${testServerDir} ${markEnd}"
  sponge web http-pb --module-name=${testServerName} --server-name=${testServerName} --project-name=ginpbdemo --protobuf-file=./files/user.proto --out=./${testServerDir}
  checkResult $?
fi

cd ${testServerDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

echo "replace the sample template code"
replaceCode ../files/web_http_pb_content ./internal/handler/user.go
checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
