#!/bin/bash

grpcServiceName="user"
grpcDir="4_micro_grpc_pb_${grpcServiceName}"

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
  for i in {1..3}; do
    echo -e "\n\n"
    echo "${i} go test -run Test_service_user_methods/Register"
    go test -run Test_service_user_methods/Register
    echo -e "\n\n"
  done
  cd -
  echo "--------------------- 测试结束！---------------------"
  stopService
}

if [ -d "${grpcDir}" ]; then
  echo "微服务 ${grpcDir} 已存在"
else
  echo "创建微服务 ${grpcDir}"
  sponge micro rpc-pb \
    --module-name=${grpcServiceName} \
    --server-name=${grpcServiceName} \
    --project-name=grpcpbdemo \
    --protobuf-file=./files/user2.proto \
    --out=./${grpcDir}
  checkResult $?
fi

cd ${grpcDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

echo "替换示例模板代码"
replaceCode ../files/micro_grpc_pb_content ./internal/service/user2.go
checkResult $?
replaceCode ../files/micro_grpc_pb_content2 ./internal/service/user2_client_test.go
checkResult $?

echo "test request"
testRequest &

echo "make run"
make run
checkResult $?
