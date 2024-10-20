#!/bin/bash

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# check if replaceCode exists
which replaceCode
checkResult $?

function main() {
  bash files/1_web_http.sh
  bash files/2_micro_grpc.sh
  bash files/3_web_http_pb.sh
  bash files/4_micro_grpc_pb.sh
  bash files/5_micro_grpc_gateway_pb.sh
  bash files/6_micro_grpc_http_pb.sh
}

main
