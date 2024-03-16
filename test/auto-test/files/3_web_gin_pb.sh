#!/bin/bash

webServiceName="user"
webDir="3_web_gin_pb_${webServiceName}"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function stopService() {
  pid=$(ps -ef | grep "./cmd/${webServiceName}" | grep -v grep | awk '{print $2}')
  if [ "${pid}" != "" ]; then
      kill -9 ${pid}
  fi
}

function testRequest() {
  echo "--------------------- 20s 后测试开始 ---------------------"
  sleep 20

  echo -e "\n\n"
  echo 'curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"'
  curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  echo -e "\n\n"
  sleep 3

  echo -e "\n\n"
  echo 'curl -X POST http://localhost:8080/api/v1/auth/register  -H "Content-Type: application/json" -H "X-Request-Id: qaz12wx3ed4" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"'
  curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -H "X-Request-Id: qaz12wx3ed4" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  echo -e "\n\n"
  sleep 3

  echo "--------------------- 测试结束！---------------------"
  stopService
}

if [ -d "${webDir}" ]; then
  echo "web服务 ${webDir} 已存在"
else
  echo "创建web服务 ${webDir}"
  sponge web http-pb \
    --module-name=${webServiceName} \
    --server-name=${webServiceName} \
    --project-name=ginpbdemo \
    --protobuf-file=./files/user.proto \
    --out=./${webDir}
  checkResult $?
fi

cd ${webDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

echo "替换示例模板代码"
replaceCode ../files/web_gin_pb_content ./internal/handler/user.go
checkResult $?

echo "test request"
testRequest &

echo "make run"
make run
checkResult $?
