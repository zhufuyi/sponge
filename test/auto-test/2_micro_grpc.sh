#!/bin/bash

grpcServiceName="user"
grpcDir="2_micro_grpc_${grpcServiceName}"

mysqlDSN="root:123456@(192.168.3.37:3306)/school"
mysqlTable="teacher"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function stopService() {
  pid=$(ps -ef | grep "./cmd/${grpcServiceName}" | grep -v grep | awk '{print $2}')
  if [ "${pid}" != "" ]; then
      kill -9 ${pid}
  fi
}

function testRequest() {
  echo "--------------------- 20s 后测试开始 ---------------------"
  sleep 20
  cd internal/service

  echo -e "\n\n"
  echo "获取详情 go test -run Test_service_${mysqlTable}_methods/GetByID id=1"
  sed -i "s/Id: 0,/Id: 1,/g" ${mysqlTable}_client_test.go
  go test -run Test_service_${mysqlTable}_methods/GetByID
  sed -i "s/Id: 1,/Id: 0,/g" ${mysqlTable}_client_test.go
  echo -e "\n\n"
  sleep 3

  echo -e "\n\n"
  echo "获取列表 go test -run Test_service_${mysqlTable}_methods/ListByLastID"
  sed -i "s/Limit: 0,/Limit: 3,/g" ${mysqlTable}_client_test.go
  go test -run Test_service_${mysqlTable}_methods/ListByLastID
  sed -i "s/Limit: 3,/Limit: 0,/g" ${mysqlTable}_client_test.go
  echo -e "\n\n"
  sleep 3

  cd -
  echo "--------------------- 测试结束！---------------------"
  stopService
}

if [ -d "${grpcDir}" ]; then
  echo "微服务 ${grpcDir} 已存在"
else
  echo "创建微服务 ${grpcDir}"
  sponge micro rpc \
    --module-name=${grpcServiceName} \
    --server-name=${grpcServiceName} \
    --project-name=grpcdemo \
    --db-dsn=${mysqlDSN} \
    --db-table=${mysqlTable} \
    --out=./${grpcDir}
  checkResult $?
fi

cd ${grpcDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

echo "test request"
testRequest &

echo "make run"
make run
checkResult $?
