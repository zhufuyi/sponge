#!/bin/bash

# grpc server
grpcServiceName="user"
grpcDir="2_micro_grpc_${grpcServiceName}"

mysqlDSN="root:123456@(192.168.3.37:3306)/account"
mysqlTable="user"

# grpc gateway server
testServerName="user_gw"
testServerDir="5_micro_grpc_gateway_${testServerName}"
projectName="edusys"
protobufFile="files/user_gw.proto"

colorCyan='\e[1;36m'
colorGreen='\e[1;32m'
colorRed='\e[1;31m'
markEnd='\e[0m'
errCount=0

srcStr1='// define rpc clients interface here'
dstStr1='userCli userV1.UserClient'

srcStr2='func (c *userGwClient) Register(ctx context.Context, req *user_gwV1.RegisterRequest) (*user_gwV1.RegisterReply, error) {'
dstStr2='func (c *userGwClient) Register(ctx context.Context, req *user_gwV1.RegisterRequest) (*user_gwV1.RegisterReply, error) {
	rep, err := c.userCli.GetByID(ctx, &userV1.GetUserByIDRequest{Id: 1})
	if err != nil {
		return nil, err
	}
	return &user_gwV1.RegisterReply{
		Id: rep.User.Id,
	}, nil
'

srcStr3='//"user_gw/internal/rpcclient"'
dstStr3='"user_gw/internal/rpcclient"
	userV1 "user_gw/api/user/v1"'

srcStr4='//	    userGwCli: user_gwV1.NewUserGwClient(rpcclient.GetUserGwRPCConn()),'
dstStr4='	    userCli: userV1.NewUserClient(rpcclient.GetUserRPCConn()),'

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
  checkServiceStarted $grpcServiceName
  checkServiceStarted $testServerName
  sleep 1

  echo -e "start testing [${testServerName}] api:\n\n"
  echo -e "${colorCyan}curl -X POST http://localhost:8080/api/v1/user/register -H \"Content-Type: application/json\" -d {\"email\":\"foo@bar.com\",\"password\":\"123456\"} ${markEnd}"
  curl -X POST http://localhost:8080/api/v1/user/register -H "Content-Type: application/json" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  checkErrCount $?

  echo -e "\n\n"
  echo -e "${colorCyan}curl -X POST http://localhost:8080/api/v1/user/register  -H \"Content-Type: application/json\" -H \"X-Request-Id: qaz12wx3ed4\" -d {\"email\":\"foo@bar.com\",\"password\":\"123456\"} ${markEnd}"
  curl -X POST http://localhost:8080/api/v1/user/register -H "Content-Type: application/json" -H "X-Request-Id: qaz12wx3ed4" -d "{\"email\":\"foo@bar.com\",\"password\":\"123456\"}"
  checkErrCount $?

  printTestResult
  stopService $grpcServiceName
  stopService $testServerName
}

function runGRPCService() {
  echo -e "\n\n"
  if [ -d "${grpcDir}" ]; then
    echo "service ${grpcDir} already exists"
  else
    echo "create service ${grpcDir}"
    echo -e "${colorCyan}sponge micro rpc --module-name=${grpcServiceName} --server-name=${grpcServiceName} --project-name=${projectName} --db-dsn=$mysqlDSN --db-table=$mysqlTable --out=./${grpcDir} ${markEnd}"
    sponge micro rpc --module-name=${grpcServiceName} --server-name=${grpcServiceName} --project-name=${projectName} --db-dsn=$mysqlDSN --db-table=$mysqlTable --out=./${grpcDir}
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

function replaceContent() {
    local file="$1"
    local src="$2"
    local dst="$3"

    if [ ! -f "$file" ]; then
        echo "File not found!"
        return 1
    fi

    # Use sed for multiline substitution to ensure special characters are parsed correctly
    sed -i.bak -e "/$(echo "$src" | sed 's/[]\/$*.^[]/\\&/g')/{
        r /dev/stdin
        d
    }" "$file" <<< "$dst"
	checkResult $?
}

function replaceContents() {
    local file="$1"
    shift

    if (( $# % 2 != 0 )); then
        echo "Error: Parameters must be in pairs of 'src' and 'dst'"
        return 1
    fi

    while (( "$#" )); do
        local src="$1"
        local dst="$2"
        replaceContent "$file" "$src" "$dst"
        shift 2
    done
}

echo "running service ${grpcDir}"
runGRPCService &

echo -e "\n\n"

if [ -d "${testServerDir}" ]; then
  echo "service ${testServerDir} already exists"
else
  echo "create service ${testServerDir}"
  echo -e "${colorCyan}sponge micro rpc-gw-pb --module-name=${testServerName} --server-name=${testServerName} --project-name=${projectName} --protobuf-file=${protobufFile} --out=./${testServerDir} ${markEnd}"
  sponge micro rpc-gw-pb --module-name=${testServerName} --server-name=${testServerName} --project-name=${projectName} --protobuf-file=${protobufFile} --out=./${testServerDir}
  checkResult $?

  echo -e "${colorCyan}sponge micro rpc-conn --rpc-server-name=${grpcServiceName} --out=./${testServerDir} ${markEnd}"
  sponge micro rpc-conn --rpc-server-name=${grpcServiceName} --out=./${testServerDir}
  checkResult $?

  echo "modify grpcClient field of configuration file"
  sed -i "s/your_grpc_service_name/user/g" ./${testServerDir}/configs/user_gw.yml
  checkResult $?

  echo "copy the proto file to the grpc gateway service directory"
  cd ${testServerDir}
  make copy-proto SERVER=../${grpcDir}
  checkResult $?
  cd -
fi

cd ${testServerDir}

echo "make proto"
make proto
checkResult $?

echo "replace template code"
replaceContents ./internal/service/user_gw.go "$srcStr1" "$dstStr1" "$srcStr2" "$dstStr2" "$srcStr3" "$dstStr3" "$srcStr4" "$dstStr4"

testRequest &

echo "running service ${testServerDir}"
make run
checkResult $?
