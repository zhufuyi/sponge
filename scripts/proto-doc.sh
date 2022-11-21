#!/bin/bash

# the directory where the proto files are located
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

# get all proto file paths
listFiles $protoBasePath

protoc --proto_path=.  --proto_path=./third_party \
  --doc_out=. --doc_opt=html,apis.html \
  $allProtoFiles

checkResult $?

mv -f apis.html docs/apis.html

echo "generate proto doc file successfully, see by 'docs/apis.html'"
