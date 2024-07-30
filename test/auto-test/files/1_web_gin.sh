#!/bin/bash

webServiceName="user"
webDir="1_web_gin_${webServiceName}"

mysqlDSN="root:123456@(192.168.3.37:3306)/school"
mysqlTable="teacher"

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

  echo -e "\n--------------------- the test is over, error result: $errCount ---------------------\n"
  stopService $webServiceName
}

if [ -d "${webDir}" ]; then
  echo "service ${webDir} already exists"
else
  echo "create service ${webDir}"
  echo -e "${colorCyan}sponge web http --module-name=${webServiceName} --server-name=${webServiceName} --project-name=webdemo --extended-api=true --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${webDir} ${markEnd}"
  sponge web http --module-name=${webServiceName} --server-name=${webServiceName} --project-name=webdemo --extended-api=true --db-dsn=${mysqlDSN} --db-table=${mysqlTable} --out=./${webDir}
  checkResult $?
fi

cd ${webDir}

echo "make docs"
make docs
checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
