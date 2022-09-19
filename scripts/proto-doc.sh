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

protoc --proto_path=.  --proto_path=./third_party \
  --doc_out=. --doc_opt=html,proto.html \
  $allProtoFiles

checkResult $?

mv -f proto.html docs/proto.html

echo "generate proto doc file successfully, see by 'docs/proto.html'"
