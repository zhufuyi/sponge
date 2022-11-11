#!/bin/bash

# proto文件所在的目录
protoBasePath="api"
allProtoFiles=""

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# 把生成*.pb.go代码中导入无用的package添加到这里
function deleteUnusedPkg() {
  file=$1
  sed -i "s#_ \"github.com/envoyproxy/protoc-gen-validate/validate\"##g" ${file}
  sed -i "s#_ \"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options\"##g" ${file}
  sed -i "s#_ \"github.com/srikrsna/protoc-gen-gotag/tagger\"##g" ${file}
  sed -i "s#_ \"google.golang.org/genproto/googleapis/api/annotations\"##g" ${file}
}

function listProtoFiles(){
    cd $1
    items=$(ls)

    for item in $items; do
        if [ -d "$item" ]; then
            listProtoFiles $item
        else
            if [ "${item#*.}"x = "proto"x ];then
              file=$(pwd)/${item}
              protoFile="${protoBasePath}${file#*${protoBasePath}}"
              allProtoFiles="${allProtoFiles} ${protoFile}"
            fi
        fi
    done
    cd ..
}

function listPbGoFiles(){
    cd $1
    items=$(ls)

    for item in $items; do
        if [ -d "$item" ]; then
            listPbGoFiles $item
        else
            if [ "${item#*.}"x = "pb.go"x ];then
              deleteUnusedPkg $item
            fi
        fi
    done
    cd ..
}

# 获取所有proto文件路径
listProtoFiles $protoBasePath

if [ "$protoBasePath"x = x ];then
  echo "Error: not found protobuf file in path $protoBasePath"
  exit 1
fi

# 生成文件 *_pb.go, *_grpc_pb.go，
protoc --proto_path=. --proto_path=./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  $allProtoFiles

checkResult $?

# 生成文件*_pb.validate.go
protoc --proto_path=. --proto_path=./third_party \
  --validate_out=lang=go:. --validate_opt=paths=source_relative \
  $allProtoFiles

checkResult $?

# 对*_pb.go字段嵌入tag
protoc --proto_path=. --proto_path=./third_party \
  --gotag_out=:. --gotag_opt=paths=source_relative \
  $allProtoFiles

checkResult $?
# todo generate router code for gin here
# delete the templates code start

# 生成swagger文档，所有文件合并到docs/apis.swagger.json
protoc --proto_path=. --proto_path=./third_party \
  --openapiv2_out=. --openapiv2_opt=logtostderr=true --openapiv2_opt=allow_merge=true --openapiv2_opt=merge_file_name=docs/apis.json \
  $allProtoFiles

checkResult $?

# 共生成4个文件，分别是注册路由文件_*router.pb.go(保存在protobuf文件同一目录)、注入路由文件*_service.pb.go(默认保存路径在internal/routers)、
# 逻辑代码模板文件*_logic.go(默认保存路径在internal/service), 返回错误码模板文件*_http.go(默认保存路径在internal/ecode)
protoc --proto_path=. --proto_path=./third_party \
  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=service \
  --go-gin_opt=moduleName=github.com/zhufuyi/sponge --go-gin_opt=serverName=serverNameExample \
  $allProtoFiles

checkResult $?
# delete the templates code end
listPbGoFiles $protoBasePath

go mod tidy
checkResult $?
echo "exec protoc command successfully."
