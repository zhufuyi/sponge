#!/bin/bash

webServiceName="user"
webDir="1_web_gin_${webServiceName}"

mysqlDSN="root:123456@(192.168.3.37:3306)/school"
mysqlTable="teacher"

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
  echo "$获取详情 [GET] curl http://localhost:8080/api/v1/${mysqlTable}/1"
  curl http://localhost:8080/api/v1/${mysqlTable}/1
  echo -e "\n\n"
  sleep 3

  echo -e "\n\n"
  echo "获取列表1 [GET] curl http://localhost:8080/api/v1/${mysqlTable}/list"
  curl http://localhost:8080/api/v1/${mysqlTable}/list
  echo -e "\n\n"
  sleep 3

  echo -e "\n\n"
  echo '获取列表2 [POST] curl -X POST http://localhost:8080/api/v1/${mysqlTable}/list -H "X-Request-Id: qaz12wx3ed4" -H "Content-Type: application/json" -d "{\"columns\":[{\"exp\":\">\",\"name\":\"id\",\"value\":1}],\"page\":0,\"size\":10}"'
  curl -X POST http://localhost:8080/api/v1/${mysqlTable}/list -H "X-Request-Id: qaz12wx3ed4" -H "Content-Type: application/json" -d "{\"columns\":[{\"exp\":\">\",\"name\":\"id\",\"value\":1}],\"page\":0,\"size\":10}"
  echo -e "\n\n"
  sleep 3

  echo ""
  echo "--------------------- 测试结束！---------------------"
  stopService
}

if [ -d "${webDir}" ]; then
  echo "web服务 ${webDir} 已存在"
else
  echo "创建web服务 ${webDir}"
  sponge web http \
    --module-name=${webServiceName} \
    --server-name=${webServiceName} \
    --project-name=webdemo \
    --db-dsn=${mysqlDSN} \
    --db-table=${mysqlTable} \
    --out=./${webDir}
  checkResult $?
fi

cd ${webDir}

echo "make docs"
make docs
checkResult $?

echo "test request"
testRequest &

echo "make run"
make run
checkResult $?
