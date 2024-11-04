#!/bin/bash

testServerName="user"
testServerDir="4_micro_grpc_pb_${testServerName}"
projectName="edusys"
protobufFile="files/user_rpc.proto"
protoServiceNameCamelFCL="userRpc"

colorCyan='\e[1;36m'
colorGreen='\e[1;32m'
colorRed='\e[1;31m'
markEnd='\e[0m'
errCount=0

srcStr1='func (s *userRpc) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterReply, error) {'
dstStr1='func (s *userRpc) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterReply, error) {
    return &userV1.RegisterReply{
        Id: 111,
    }, nil
'
srcStr2='Email:    "",'
dstStr2='Email:    "foo@bar.com",'
srcStr3='Password: "",'
dstStr3='Password: "123456",'

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

  cd internal/service
  echo -e "start testing [${testServerName}] api:\n\n"
  echo -e "${colorCyan}go test -run Test_service_${protoServiceNameCamelFCL}_methods/Register ${markEnd}"
  go test -run Test_service_${protoServiceNameCamelFCL}_methods/Register
  checkErrCount $?

  cd -
  printTestResult
  stopService $testServerName
}

function replaceContent() {
    local file="$1"
    local src="$2"
    local dst="$3"

    if [ ! -f "$file" ]; then
        echo "file $file not found!"
        return 1
    fi

    # Use sed for multiline substitution to ensure special characters are parsed correctly
    sed -i.bak -e "/$(echo "$src" | sed 's/[]\/$*.^[]/\\&/g')/{
        r /dev/stdin
        d
    }" "$file" <<< "$dst"
	checkResult $?
}

echo -e "\n\n"

if [ -d "${testServerDir}" ]; then
  echo "service ${testServerDir} already exists"
else
  echo "create service ${testServerDir}"
  echo -e "${colorCyan}sponge micro rpc-pb --module-name=${testServerName} --server-name=${testServerName} --project-name=${projectName} --protobuf-file=${protobufFile} --out=./${testServerDir} ${markEnd}"
  sponge micro rpc-pb --module-name=${testServerName} --server-name=${testServerName} --project-name=${projectName} --protobuf-file=${protobufFile} --out=./${testServerDir}
  checkResult $?
fi

cd ${testServerDir}
checkResult $?

echo "make proto"
make proto
checkResult $?

echo "replace template code"
replaceContent ./internal/service/user_rpc.go "$srcStr1" "$dstStr1"
checkResult $?
replaceContent ./internal/service/user_rpc_client_test.go "$srcStr2" "$dstStr2"
checkResult $?
replaceContent ./internal/service/user_rpc_client_test.go "$srcStr3" "$dstStr3"
checkResult $?

testRequest &

echo "make run"
make run
checkResult $?
