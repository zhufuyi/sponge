#!/bin/bash

SERVER_NAME="serverNameExample"
DEPLOY_FILE="deployments/kubernetes/${SERVER_NAME}-deployment.yml"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# 判断文件是否存在
if [ ! -f "${DEPLOY_FILE}" ];then
  echo "部署文件文件${DEPLOY_FILE}不存在"
  checkResult 1
fi

# 检查是否授权操作k8s
echo "kubectl version"
kubectl version
checkResult $?

echo "kubectl delete -f ${DEPLOY_FILE} --ignore-not-found"
kubectl delete -f ${DEPLOY_FILE} --ignore-not-found
checkResult $?

sleep 1

echo "kubectl apply -f ${DEPLOY_FILE}"
kubectl apply -f ${DEPLOY_FILE}
