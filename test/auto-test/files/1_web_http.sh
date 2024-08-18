#!/bin/bash

testServerName="user"
testServerDir="1_web_http_${testServerName}"

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

  echo -e "start testing [${testServerName}] api:\n\n"
  echo -e "${colorCyan}curl -X GET http://localhost:8080/api/v1/${mysqlTable}/1 ${markEnd}"
  curl -X GET http://localhost:8080/api/v1/${mysqlTable}/1
  checkErrCount $?

  echo -e "\n\n"
  echo -e "${colorCyan}curl -X GET http://localhost:8080/api/v1/${mysqlTable}/list ${markEnd}"
  curl -X GET http://localhost:8080/api/v1/${mysqlTable}/list
  checkErrCount $?

  echo -e "\n\n"
  echo -e "${colorCyan}curl -X POST http://localhost:8080/api/v1/${mysqlTable}/list -H \"X-Request-Id: qaz12wx3ed4\" -H \"Content-Type: application/json\" -d {\"columns\":[{\"exp\":\">\",\"name\":\"id\",\"value\":1}],\"page\":0,\"limit\":10} ${markEnd}"
  curl -X POST http://localhost:8080/api/v1/${mysqlTable}/list -H "X-Request-Id: qaz12wx3ed4" -H "Content-Type: application/json" -d "{\"columns\":[{\"exp\":\">\",\"name\":\"id\",\"value\":1}],\"page\":0,\"limit\":10}"
  checkErrCount $?

  printTestResult
  stopService $testServerName
}

echo -e "\n\n"

if [ -d "${testServerDir}" ]; then
  echo "service ${testServerDir} already exists"
else
  echo "create service ${testServerDir}"
  echo -e "${colorCyan}sponge web http --module-name=${testServerName} --server-name=${testServerName} --project-name=webdemo --extended-api=true --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${testServerDir} ${markEnd}"
  sponge web http --module-name=${testServerName} --server-name=${testServerName} --project-name=webdemo --extended-api=true --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${testServerDir}
  checkResult $?
fi

cd ${testServerDir}

echo "make docs"
make docs
checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
