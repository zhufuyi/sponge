#!/bin/bash

# grpc服务
grpcServiceName="user"
grpcDir="2_micro_grpc_${grpcServiceName}"

# grpc网关服务
rpcGwServiceName="user_gw"
rpcGwDir="5_grpc_gateway_${rpcGwServiceName}"

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

  pid=$(ps -ef | grep "./cmd/${rpcGwServiceName}" | grep -v grep | awk '{print $2}')
  if [ "${pid}" != "" ]; then
      kill -9 ${pid}
  fi
}

function testRequest() {
  echo "--------------------- 20s 后测试开始 ---------------------"
  sleep 20
  for i in {1..3}; do
    echo -e "\n\n"
    echo ${i} 'curl -X POST http://localhost:8080/api/v1/user/register -H "Content-Type: application/json" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"'
    curl -X POST http://localhost:8080/api/v1/user/register -H "Content-Type: application/json" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
    echo -e "\n\n"
    sleep 3

    echo -e "\n\n"
    echo ${i} 'curl -X POST http://localhost:8080/api/v1/user/register  -H "Content-Type: application/json" -H "X-Request-Id: qaz12wx3ed4" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"'
    curl -X POST http://localhost:8080/api/v1/user/register -H "Content-Type: application/json" -H "X-Request-Id: qaz12wx3ed4" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
    echo -e "\n\n"
    sleep 3
  done
  echo "--------------------- 测试结束！---------------------"
  stopService
}

function runGRPCService() {
  if [ -d "${grpcDir}" ]; then
    echo "微服务 ${grpcDir} 已存在"
  else
    echo "创建微服务 ${grpcDir}"
    sponge micro rpc \
      --module-name=${grpcServiceName} \
      --server-name=${grpcServiceName} \
      --project-name=grpcdemo \
      --db-dsn="root:123456@(192.168.3.37:3306)/school" \
      --db-table=teacher \
      --out=./${grpcDir}
    checkResult $?
  fi

  cd ${grpcDir}
  checkResult $?

  echo "make proto"
  make proto
  checkResult $?

  echo "make run"
  make run
  checkResult $?
}

echo "运行微服务 ${grpcDir}"
runGRPCService &

if [ -d "${rpcGwDir}" ]; then
  echo "rpc网关服务 ${rpcGwDir} 已存在"
else
  echo "创建rpc网关服务 ${rpcGwDir}"
  sponge micro rpc-gw-pb \
    --module-name=${rpcGwServiceName} \
    --server-name=${rpcGwServiceName} \
    --project-name=grpcgwdemo \
    --protobuf-file=./files/user_gw.proto \
    --out=./${rpcGwDir}
  checkResult $?

  echo "生成rpc连接服务代码"
  sponge micro rpc-conn --rpc-server-name=${grpcServiceName} --out=./${rpcGwDir}
  checkResult $?

  echo "修改配置文件的grpcClient字段"
  sed -i "s/your_grpc_service_name/user/g" ./${rpcGwDir}/configs/user_gw.yml
  checkResult $?

  echo "复制proto文件到rpc网关目录"
  cd ${rpcGwDir}
  make copy-proto SERVER=../${grpcDir}
  checkResult $?
  cd -
fi

cd ${rpcGwDir}

echo "make proto"
make proto
checkResult $?

echo "替换示例模板代码"
replaceCode ../files/rpc_gateway_content ./internal/service/user_gw.go

echo "test request"
testRequest &

echo "运行rpc网关服务 ${rpcGwDir}"
make run
checkResult $?
