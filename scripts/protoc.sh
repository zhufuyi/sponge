#!/bin/bash

# 插件版本
# protoc                               v3.20.1
# protoc-gen-go                   v1.28.0
# protoc-gen-validate           v0.6.7

# proto文件所在的目录
protoBasePath="api"
allProtoFiles=""

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function listFiles(){
    cd $1
    items=$(ls)

    for item in $items; do
        if [ -d "$item" ]; then
            listFiles $item
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

# 获取所有proto文件路径
listFiles $protoBasePath

# 生成文件 *_pb.go, *_grpc_pb.go, *_pb.validate.go文件在同一目录
protoc --proto_path=. --proto_path=./third_party \
 --go_out=. --go_opt=paths=source_relative \
 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
 --validate_out=lang=go:. --validate_opt=paths=source_relative \
 $allProtoFiles

checkResult $?

echo "exec protoc command successfully."
